package worker

import (
	"context"
	"httpfetcher/internal/logger"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Result struct {
	JobID    string
	Duration time.Duration
	Success  bool
	Error    error
	Attempts int
}

type Pool struct {
	Workers   int
	Retries   int
	RateLimit time.Duration
	Timeout   time.Duration
	Seed      int
	Results   []Result
	mu        sync.Mutex
}

// JobFunc - intentionally generic to allow reuse beyond HTTP fetching.
// Could unify workers and fetcher logic, but this design enables the worker Pool to handle any job type (file processing, API calls, etc.)
type JobFunc func(ctx context.Context, job string) error

func NewPool(workers, retries int, rateLimit time.Duration, timeout time.Duration, seed int) *Pool {
	return &Pool{
		Workers:   workers,
		Retries:   retries,
		RateLimit: rateLimit,
		Timeout:   timeout,
		Seed:      seed,
		Results:   []Result{},
	}
}

// PerformJob - orchestrates the worker pool
func (w *Pool) PerformJob(ctx context.Context, jobs <-chan string, jobFunc JobFunc, wg *sync.WaitGroup) {
	log := logger.FromContext(ctx)
	workerChan := make(chan struct{}, w.Workers)

	for job := range jobs {
		wg.Add(1)
		log.Debug("Starting Go routine to process job", zap.String("job", job))
		go w.processJob(ctx, job, jobFunc, wg, workerChan)
	}
}

// processJob - handles a single job with retries
func (w *Pool) processJob(ctx context.Context, job string, jobFunc JobFunc, wg *sync.WaitGroup, workerChan chan struct{}) {
	log := logger.FromContext(ctx)
	defer wg.Done()
	defer w.releaseWorker(workerChan)

	workerChan <- struct{}{} // Acquire worker

	startTime := time.Now()

	// Apply rate limiting
	w.applyRateLimit()

	// Execute job with retries
	result := w.executeWithRetries(ctx, job, jobFunc, startTime)

	// Record result
	w.recordResult(result)
	if result.Success {
		log.Debug("Successfully processed job", zap.String("job", job))
	} else {
		log.Error("Failed to process job", zap.String("job", job), zap.Error(result.Error))
	}
}

// executeWithRetries - handles the retry logic
func (w *Pool) executeWithRetries(ctx context.Context, job string, jobFunc JobFunc, startTime time.Time) Result {
	log := logger.FromContext(ctx)
	var err error
	for attempt := 1; attempt <= w.Retries; attempt++ {
		// Add timeout to the context
		ctx, cancel := context.WithTimeout(ctx, w.Timeout)
		defer cancel()

		err = jobFunc(ctx, job)
		if err == nil {
			return Result{
				JobID:    job,
				Duration: time.Since(startTime).Round(time.Millisecond),
				Success:  true,
				Attempts: attempt,
			}
		}

		if attempt < w.Retries {
			delay := w.calculateBackoff(attempt)
			log.Debug("Retrying job", zap.String("job", job), zap.Duration("delay", delay), zap.Int("attempt", attempt+1), zap.Int("retries", w.Retries))
			time.Sleep(delay)
		}
	}

	// All retries failed
	return Result{
		JobID:    job,
		Duration: time.Since(startTime).Round(time.Millisecond),
		Success:  false,
		Error:    err,
		Attempts: w.Retries,
	}
}

// Helper methods
func (w *Pool) releaseWorker(workerChan chan struct{}) {
	<-workerChan
}

func (w *Pool) applyRateLimit() {
	if w.RateLimit > 0 {
		time.Sleep(time.Second / w.RateLimit)
	}
}

func (w *Pool) calculateBackoff(attempt int) time.Duration {
	baseDelay := time.Duration(1<<attempt) * time.Second
	jitter := time.Duration(w.Seed%100) * time.Millisecond
	return baseDelay + jitter
}

func (w *Pool) recordResult(result Result) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Results = append(w.Results, result)
}
