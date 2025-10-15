# HTTP Fetcher

A high-performance HTTP URL fetcher with worker pool, retry logic, and structured logging.

While existing tools like `curl`, `wget`, or monitoring solutions may be simpler for basic use cases, this project is designed for environments where you need to perform monitoring of multiple endpoints, perform health checks, or validate URL accessibility at scale.

Future improvements can extend the tool to retrieve response bodies, add POST request support, accept input from files/stdin, and integrate with monitoring systems.

## Features

- **Worker Pool**: Configurable concurrent workers for parallel URL fetching
- **Retry Logic**: Exponential backoff with jitter for failed requests
- **Rate Limiting**: Configurable requests per second
- **Timeout**: Configurable timeout for each request (set in the job level context)
- **Structured Logging**: Zap-based logging with context propagation
- **Graceful Shutdown**: Signal handling for clean termination

## Usage

```bash
# Basic usage
go run main.go -urls="google.com,github.com,stackoverflow.com,amazon.com,google.ie,facebook.com"

# With custom configuration
go run main.go \
  -urls="google.com,github.com" \
  -workers=3 \
  -rps=5 \
  -retries=3
```

## Implementation Details

- **Worker Pool Architecture**: Implements controlled concurrency with configurable workers to prevent resource exhaustion and enable production-scale job processing
- **Worker Package**: Reusable worker package for any job type (file processing, API calls, etc.)
- **Context-Based Design**: Uses Go's context package for cancellation, timeouts, and request-scoped logging propagation throughout the call chain
- **Production Error Handling**: Comprehensive error wrapping with context, graceful degradation, and structured logging for observability
- **Exponential Backoff**: Implements retry logic with jitter to handle transient failures without clashing with other parallel jobs


## Future Improvements

- **Unit Tests**: Add unit tests for worker pool
- **Metrics Collection**: Request rate, success/failure ratios, and latency percentiles
- **Extend Input support**: Accept input from file, stdin, or other sources
- **Add Post Request Support**: Add support for POST requests, to send API requests.
