package emailit

import (
	"bytes"
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

// Config holds Emailit API configuration
type Config struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
	Retries int
}

// Emailit is the email sending client
type Emailit struct {
	config     Config
	httpClient *http.Client
}

// SendEmailRequest represents an email sending request
type SendEmailRequest struct {
	To           string
	Subject      string
	TemplateID   string
	TemplateData map[string]interface{}
}

// Response represents the Emailit API response
type Response struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
}

// NewEmailit creates a new Emailit client
func NewEmailit(config Config) *Emailit {
	return &Emailit{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// SendEmail sends an email using Emailit API
func (e *Emailit) SendEmail(ctx context.Context, req SendEmailRequest) error {
	ctx, span := otel.Tracer("emailit").Start(ctx, "emailit.SendEmail")
	defer span.End()

	span.SetAttributes(
		attribute.String("email.to", req.To),
		attribute.String("email.template_id", req.TemplateID),
		attribute.String("emailit.api_url", e.config.BaseURL),
	)

	var lastErr error
	for attempt := 0; attempt <= e.config.Retries; attempt++ {
		if attempt > 0 {
			span.AddEvent(fmt.Sprintf("retry attempt %d", attempt))
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond) // Exponential backoff
		}

		err := e.sendEmailOnce(ctx, req)
		if err == nil {
			span.SetStatus(codes.Ok, "")
			return nil
		}

		lastErr = err
		span.RecordError(err)
	}

	span.SetStatus(codes.Error, lastErr.Error())
	return fmt.Errorf("failed to send email after %d retries: %w", e.config.Retries, lastErr)
}

// sendEmailOnce performs a single email sending request
func (e *Emailit) sendEmailOnce(ctx context.Context, req SendEmailRequest) error {
	// Build request payload
	payload := map[string]interface{}{
		"to":            req.To,
		"subject":       req.Subject,
		"template_id":   req.TemplateID,
		"template_data": req.TemplateData,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build request URL
	url := fmt.Sprintf("%s/emails", e.config.BaseURL)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	httpReq.Header.Set("Authorization", "Bearer "+e.config.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := e.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		var apiResp Response
		if json.Unmarshal(body, &apiResp) == nil && apiResp.Error != "" {
			return fmt.Errorf("API error: %s", apiResp.Error)
		}
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResp Response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API-level errors
	if apiResp.Error != "" {
		return fmt.Errorf("API error: %s", apiResp.Error)
	}

	// Verify email was sent
	if apiResp.Status != "sent" {
		return fmt.Errorf("email not sent, status: %s", apiResp.Status)
	}

	return nil
}

