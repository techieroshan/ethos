package email

import (
	"context"
)

// EmailSender defines the interface for sending emails
type EmailSender interface {
	SendEmail(ctx context.Context, req SendEmailRequest) error
}

// SendEmailRequest represents an email sending request
type SendEmailRequest struct {
	To           string
	Subject      string
	TemplateID   string
	TemplateData map[string]interface{}
}

// TemplateData holds template variables
type TemplateData map[string]interface{}

// NoOpEmailSender is a no-op email sender for testing
type NoOpEmailSender struct{}

// SendEmail implements EmailSender interface (no-op)
func (n *NoOpEmailSender) SendEmail(ctx context.Context, req SendEmailRequest) error {
	// No-op: do nothing, useful for testing
	return nil
}

// NewNoOpEmailSender creates a no-op email sender
func NewNoOpEmailSender() EmailSender {
	return &NoOpEmailSender{}
}

