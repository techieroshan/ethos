package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Config holds Checker API configuration
type Config struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
	Retries int
}

// Checker is the email validation client
type Checker struct {
	config     Config
	httpClient *http.Client
}

// Response represents the Checker API response
type Response struct {
	Valid      bool   `json:"valid"`
	Disposable bool   `json:"disposable"`
	Error      string `json:"error,omitempty"`
}

// NewChecker creates a new Checker client
func NewChecker(config Config) *Checker {
	return &Checker{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// ValidateEmail validates an email address using the Checker API
func (c *Checker) ValidateEmail(ctx context.Context, email string) (bool, error) {
	ctx, span := otel.Tracer("checker").Start(ctx, "checker.ValidateEmail")
	defer span.End()

	span.SetAttributes(
		attribute.String("email", email),
		attribute.String("checker.api_url", c.config.BaseURL),
	)

	var lastErr error
	for attempt := 0; attempt <= c.config.Retries; attempt++ {
		if attempt > 0 {
			span.AddEvent(fmt.Sprintf("retry attempt %d", attempt))
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond) // Exponential backoff
		}

		valid, err := c.validateEmailOnce(ctx, email)
		if err == nil {
			span.SetAttributes(attribute.Bool("email.valid", valid))
			span.SetStatus(codes.Ok, "")
			return valid, nil
		}

		lastErr = err
		span.RecordError(err)
	}

	span.SetStatus(codes.Error, lastErr.Error())
	return false, fmt.Errorf("failed to validate email after %d retries: %w", c.config.Retries, lastErr)
}

// validateEmailOnce performs a single validation request
func (c *Checker) validateEmailOnce(ctx context.Context, email string) (bool, error) {
	// Build request URL
	url := fmt.Sprintf("%s/validate?email=%s", c.config.BaseURL, email)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key header
	req.Header.Set("X-API-Key", c.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResp Response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return false, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API-level errors
	if apiResp.Error != "" {
		return false, fmt.Errorf("API error: %s", apiResp.Error)
	}

	// Validate email
	if !apiResp.Valid {
		return false, fmt.Errorf("invalid email format")
	}

	// Check if disposable/temporary
	if apiResp.Disposable {
		return false, fmt.Errorf("temporary email addresses are not allowed")
	}

	return true, nil
}

