package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// FetchUrl - fetch a URL and return an error if the status code is not 200
func FetchUrl(ctx context.Context, url string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Create HTTP client (no need to set timeout here, it will be set in the context)
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for %s: %w", url, err)
	}

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("HTTP %d: %s for URL %s", response.StatusCode, response.Status, url)
	}

	_, err = io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body for URL %s: %w", url, err)
	}

	return nil
}
