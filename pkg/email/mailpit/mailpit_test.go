package mailpit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"ethos/pkg/email"
)

func TestMailpit_SendEmail(t *testing.T) {
	tests := []struct {
		name        string
		req         email.SendEmailRequest
		wantErr     bool
		wantErrMsg  string
		skipIfNoSMTP bool
	}{
		{
			name: "successful email send via SMTP",
			req: email.SendEmailRequest{
				To:           "test@example.com",
				Subject:      "Test Email",
				TemplateID:   "verification",
				TemplateData: map[string]interface{}{"name": "John"},
			},
			wantErr: false,
			skipIfNoSMTP: true,
		},
		{
			name: "invalid email address",
			req: email.SendEmailRequest{
				To:           "invalid-email",
				Subject:      "Test",
				TemplateID:   "verification",
				TemplateData: map[string]interface{}{},
			},
			wantErr:    true,
			wantErrMsg: "invalid email",
		},
		{
			name: "missing required fields",
			req: email.SendEmailRequest{
				To:           "",
				Subject:      "Test",
				TemplateID:   "verification",
				TemplateData: map[string]interface{}{},
			},
			wantErr:    true,
			wantErrMsg: "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipIfNoSMTP {
				t.Skip("Skipping SMTP test - requires Mailpit running")
			}

			client := NewMailpit(Config{
				SMTPHost:  "localhost",
				SMTPPort:  1025,
				FromEmail: "noreply@ethos.test",
			})

			err := client.SendEmail(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMailpit_GetEmails(t *testing.T) {
	client := NewMailpit(Config{
		SMTPHost:  "localhost",
		SMTPPort:  8025, // Mailpit web UI port
		FromEmail: "noreply@ethos.test",
	})

	// This test requires Mailpit to be running
	t.Skip("Skipping - requires Mailpit running")

	emails, err := client.GetEmails(context.Background(), 10)

	require.NoError(t, err)
	assert.NotNil(t, emails)
}

func TestNewMailpit(t *testing.T) {
	config := Config{
		SMTPHost:  "localhost",
		SMTPPort:  1025,
		FromEmail: "noreply@ethos.test",
	}

	client := NewMailpit(config)

	assert.NotNil(t, client)
	assert.Equal(t, config.SMTPHost, client.config.SMTPHost)
	assert.Equal(t, config.SMTPPort, client.config.SMTPPort)
	assert.Equal(t, config.FromEmail, client.config.FromEmail)
}

func TestMailpit_ValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				SMTPHost:  "localhost",
				SMTPPort:  1025,
				FromEmail: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing SMTP host",
			config: Config{
				SMTPPort:  1025,
				FromEmail: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "missing from email",
			config: Config{
				SMTPHost: "localhost",
				SMTPPort: 1025,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

