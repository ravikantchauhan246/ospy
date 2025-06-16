package monitor

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// CheckResult represents the result of a website check

type CheckResult struct{
	URL          string
	Status       int
	ResponseTime time.Duration
	Error        error
	Timestamp    time.Time
	IsUp         bool
}

// Checker handles HTTP requests to websites

type Checker struct{
	client *http.Client
}

// NewChecker creates a new HTTP checker with specified timeout
func NewChecker(timeout time.Duration) *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Check performs an HTTP request to the given URL
func (c *Checker) Check(ctx context.Context, url string) CheckResult {
	start := time.Now()
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return CheckResult{
			URL:       url,
			Error:     fmt.Errorf("failed to create request: %w", err),
			Timestamp: time.Now(),
			IsUp:      false,
		}
	}

	resp, err := c.client.Do(req)
	responseTime := time.Since(start)

	result := CheckResult{
		URL:          url,
		ResponseTime: responseTime,
		Timestamp:    time.Now(),
	}

	if err != nil {
		result.Error = fmt.Errorf("request failed: %w", err)
		result.IsUp = false
		return result
	}
	defer resp.Body.Close()

	result.Status = resp.StatusCode
	result.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 400

	return result
}

