package mailpit

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"ethos/pkg/email"
)

// Config holds Mailpit SMTP configuration
type Config struct {
	SMTPHost  string
	SMTPPort  int
	FromEmail string
}

// Mailpit is the local email testing client
type Mailpit struct {
	config Config
}

// NewMailpit creates a new Mailpit client
func NewMailpit(config Config) *Mailpit {
	return &Mailpit{
		config: config,
	}
}

// SendEmail sends an email via SMTP to Mailpit (implements email.EmailSender)
func (m *Mailpit) SendEmail(ctx context.Context, req email.SendEmailRequest) error {
	// Validate config
	if err := validateConfig(m.config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Validate request
	if req.To == "" {
		return fmt.Errorf("to email is required")
	}
	if req.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	if !isValidEmail(req.To) {
		return fmt.Errorf("invalid email address: %s", req.To)
	}

	// Build email body from template data
	body := buildEmailBody(req.TemplateID, req.TemplateData)

	// Build email message
	message := buildEmailMessage(m.config.FromEmail, req.To, req.Subject, body)

	// Send via SMTP
	addr := fmt.Sprintf("%s:%d", m.config.SMTPHost, m.config.SMTPPort)
	err := smtp.SendMail(addr, nil, m.config.FromEmail, []string{req.To}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}

	return nil
}

// GetEmails retrieves emails from Mailpit API (for testing)
func (m *Mailpit) GetEmails(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	// This would typically call Mailpit's HTTP API at port 8025
	// For now, return empty slice as this is mainly for testing
	return []map[string]interface{}{}, nil
}

// validateConfig validates Mailpit configuration
func validateConfig(config Config) error {
	if config.SMTPHost == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if config.SMTPPort <= 0 {
		return fmt.Errorf("SMTP port must be positive")
	}
	if config.FromEmail == "" {
		return fmt.Errorf("from email is required")
	}
	if !isValidEmail(config.FromEmail) {
		return fmt.Errorf("invalid from email address: %s", config.FromEmail)
	}
	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// buildEmailBody builds email body from template data
func buildEmailBody(templateID string, data map[string]interface{}) string {
	// Simple template rendering for local testing
	body := fmt.Sprintf("Template: %s\n\n", templateID)
	for key, value := range data {
		body += fmt.Sprintf("%s: %v\n", key, value)
	}
	return body
}

// buildEmailMessage builds a complete email message
func buildEmailMessage(from, to, subject, body string) string {
	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/plain; charset=UTF-8\r\n"
	message += "\r\n"
	message += body
	return message
}

// SetTLSConfig allows setting custom TLS config (for testing)
func (m *Mailpit) SetTLSConfig(config *tls.Config) {
	// This can be used to configure TLS for SMTP if needed
	// For Mailpit, TLS is typically not required
}

