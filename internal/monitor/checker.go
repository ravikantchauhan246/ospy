package monitor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// CheckResult represents the result of a website check
type CheckResult struct {
	WebsiteName  string
	URL          string
	Status       int
	ResponseTime time.Duration
	Error        error
	Timestamp    time.Time
	IsUp         bool
	Message      string
}

// Checker handles HTTP requests to websites
type Checker struct {
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
	website := Website{
		Name: url,
		URL:  url,
		Method: "GET",
	}
	return c.CheckWebsite(ctx, website)
}

// CheckWebsite performs an HTTP request to the given website
func (c *Checker) CheckWebsite(ctx context.Context, website Website) CheckResult {
	start := time.Now()
	
	method := website.Method
	if method == "" {
		method = "GET"
	}
	
	req, err := http.NewRequestWithContext(ctx, method, website.URL, nil)
	if err != nil {
		return CheckResult{
			WebsiteName: website.Name,
			URL:         website.URL,
			Error:       fmt.Errorf("failed to create request: %w", err),
			Timestamp:   time.Now(),
			IsUp:        false,
			Message:     "Failed to create HTTP request",
		}
	}

	// Add custom headers
	for key, value := range website.Headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	responseTime := time.Since(start)

	result := CheckResult{
		WebsiteName:  website.Name,
		URL:          website.URL,
		ResponseTime: responseTime,
		Timestamp:    time.Now(),
	}

	if err != nil {
		result.Error = fmt.Errorf("request failed: %w", err)
		result.IsUp = false
		result.Message = "HTTP request failed"
		return result
	}
	defer resp.Body.Close()

	result.Status = resp.StatusCode

	// Check if status is expected
	expectedStatus := website.ExpectedStatus
	if expectedStatus == 0 {
		expectedStatus = 200
	}

	if resp.StatusCode == expectedStatus {
		result.IsUp = true
		result.Message = fmt.Sprintf("Status %d (as expected)", resp.StatusCode)
	} else {
		result.IsUp = false
		result.Message = fmt.Sprintf("Status %d (expected %d)", resp.StatusCode, expectedStatus)
	}

	// Check content if specified
	if website.CheckContent != "" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			result.Error = fmt.Errorf("failed to read response body: %w", err)
			result.IsUp = false
			result.Message = "Failed to read response body"
			return result
		}

		if !strings.Contains(string(body), website.CheckContent) {
			result.IsUp = false
			result.Message = fmt.Sprintf("Content check failed: '%s' not found", website.CheckContent)
		}
	}

	return result
}