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

