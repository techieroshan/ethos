package emailit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

func TestEmailit_SendEmail(t *testing.T) {
	tests := []struct {
		name           string
		to             string
		subject        string
		templateID     string
		templateData   map[string]interface{}
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		wantErrMsg     string
	}{
		{
			name:           "successful email send",
			to:             "user@example.com",
			subject:        "Test Email",
			templateID:     "verification",
			templateData:   map[string]interface{}{"name": "John"},
			mockResponse:   `{"message_id": "msg-123", "status": "sent"}`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "API error",
			to:             "user@example.com",
			subject:        "Test Email",
			templateID:     "verification",
			templateData:   map[string]interface{}{"name": "John"},
			mockResponse:   `{"error": "Invalid API key"}`,
			mockStatusCode: http.StatusUnauthorized,
			wantErr:        true,
			wantErrMsg:     "API error",
		},
		{
			name:           "network timeout",
			to:             "user@example.com",
			subject:        "Test Email",
			templateID:     "verification",
			templateData:   map[string]interface{}{"name": "John"},
			mockResponse:   ``,
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
		{
			name:           "invalid email address",
			to:             "invalid-email",
			subject:        "Test Email",
			templateID:     "verification",
			templateData:   map[string]interface{}{"name": "John"},
			mockResponse:   `{"error": "Invalid recipient email"}`,
			mockStatusCode: http.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.name == "network timeout" {
					time.Sleep(2 * time.Second) // Simulate timeout
					return
				}
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create emailit client
			client := NewEmailit(Config{
				APIKey:  "test-api-key",
				BaseURL: server.URL,
				Timeout: 1 * time.Second,
				Retries: 1,
			})

			// Test email sending
			err := client.SendEmail(context.Background(), SendEmailRequest{
				To:          tt.to,
				Subject:     tt.subject,
				TemplateID:  tt.templateID,
				TemplateData: tt.templateData,
			})

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

func TestEmailit_SendEmail_OpenTelemetry(t *testing.T) {
	// Setup OpenTelemetry tracer
	tracer := otel.Tracer("test")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message_id": "msg-123", "status": "sent"}`))
	}))
	defer server.Close()

	client := NewEmailit(Config{
		APIKey:  "test-api-key",
		BaseURL: server.URL,
		Timeout: 1 * time.Second,
		Retries: 1,
	})

	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	err := client.SendEmail(ctx, SendEmailRequest{
		To:          "user@example.com",
		Subject:     "Test",
		TemplateID:  "verification",
		TemplateData: map[string]interface{}{"name": "John"},
	})

	require.NoError(t, err)
}

func TestEmailit_SendEmail_RetryLogic(t *testing.T) {
	retryCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryCount++
		if retryCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message_id": "msg-123", "status": "sent"}`))
	}))
	defer server.Close()

	client := NewEmailit(Config{
		APIKey:  "test-api-key",
		BaseURL: server.URL,
		Timeout: 1 * time.Second,
		Retries: 3,
	})

	err := client.SendEmail(context.Background(), SendEmailRequest{
		To:          "user@example.com",
		Subject:     "Test",
		TemplateID:  "verification",
		TemplateData: map[string]interface{}{"name": "John"},
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, retryCount)
}

func TestNewEmailit(t *testing.T) {
	config := Config{
		APIKey:  "test-key",
		BaseURL: "https://api.emailit.com",
		Timeout: 5 * time.Second,
		Retries: 2,
	}

	client := NewEmailit(config)

	assert.NotNil(t, client)
	assert.Equal(t, config.APIKey, client.config.APIKey)
	assert.Equal(t, config.BaseURL, client.config.BaseURL)
}

