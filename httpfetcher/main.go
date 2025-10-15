package main

import (
	"context"
	"errors"
	"flag"
	"httpfetcher/internal/fetcher"
	"httpfetcher/internal/logger"
	"httpfetcher/pkg/worker"
	"math/rand/v2"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	workers int
	timeout int
	rps     int
	retries int
	urls    string
	seed    int
	wg      sync.WaitGroup
	noDebug bool
)

func main() {
	flag.IntVar(&workers, "workers", 3, "number of workers")
	flag.IntVar(&timeout, "timeout", 30, "timeout in seconds")
	flag.IntVar(&rps, "rps", 5, "number of requests per second")
	flag.IntVar(&retries, "retries", 3, "number of retries")
	flag.IntVar(&seed, "seed", int(time.Now().Unix()), "Random seed for exponatial backoff (default: current time)")
	flag.StringVar(&urls, "urls", "", "urls to fetch")
	flag.BoolVar(&noDebug, "disable-debug", true, "disable debug mode")
	flag.Parse()

	// Initialize Zap logger in development mode by default
	// Future implmentation can extend the production logger to include more fields, etc.
	log := logger.Init(noDebug)
	defer log.Sync()

	log.Info("Starting HTTP Fetcher")
	if err := validateFlags(); err != nil {
		log.Error("Invalid flags", zap.Error(err))
		return
	}

	// conver urls to array
	urlsArray := strings.Split(urls, ",")
	if len(urlsArray) == 0 {
		log.Error("urls are required")
		return
	}

	// Add https:// prefix to URLs that don't have a protocol
	for i, url := range urlsArray {
		url = strings.TrimSpace(url)
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			urlsArray[i] = "https://" + url
		} else {
			urlsArray[i] = url
		}
	}
	log.Debug("URLs to process", zap.Strings("urls", urlsArray))

	// Line up
	jobs := make(chan string, len(urlsArray))
	for _, url := range urlsArray {
		jobs <- url
	}
	close(jobs)

	// Add signal context (for graceful shutdowns)
	ctxSg, cancelSg := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelSg()

	// Add logger to context
	ctxSg = logger.WithLogger(ctxSg, log)

	// Add Cancel option to context
	ctx, cancel := context.WithCancel(ctxSg)
	defer cancel()

	log.Debug("Worker pool configuration", zap.Int("workers", workers), zap.Int("retries", retries), zap.Int("rps", rps), zap.Int("seed", seed))
	workerPool := &worker.Pool{
		Workers:   workers,
		Retries:   retries,
		RateLimit: time.Duration(rps),
		Timeout:   time.Duration(timeout) * time.Second,
		Seed:      rand.IntN(seed),
	}

	log.Info("Starting worker pool execution")
	workerPool.PerformJob(ctx, jobs, fetcher.FetchUrl, &wg)
	wg.Wait()

	// Count successful vs failed jobs
	successful := 0
	for _, result := range workerPool.Results {
		if result.Success {
			successful++
		}
	}

	log.Info("All tasks completed",
		zap.Int("total_jobs", len(workerPool.Results)),
		zap.Int("successful", successful),
		zap.Int("failed", len(workerPool.Results)-successful),
	)
}

func validateFlags() error {
	if workers <= 0 {
		return errors.New("workers must be greater than 0")
	}
	if retries <= 0 {
		return errors.New("retries must be greater than 0")
	}
	if rps <= 0 {
		return errors.New("rps must be greater than 0")
	}
	if seed <= 0 {
		return errors.New("seed must be greater than 0")
	}
	return nil
}
