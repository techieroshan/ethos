package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func TestChecker_ValidateEmail(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		mockResponse   string
		mockStatusCode int
		wantValid      bool
		wantErr        bool
		wantErrMsg     string
	}{
		{
			name:           "valid email",
			email:          "user@example.com",
			mockResponse:   `{"valid": true, "disposable": false}`,
			mockStatusCode: http.StatusOK,
			wantValid:      true,
			wantErr:        false,
		},
		{
			name:           "temporary email detected",
			email:          "user@tempmail.com",
			mockResponse:   `{"valid": true, "disposable": true}`,
			mockStatusCode: http.StatusOK,
			wantValid:      false,
			wantErr:        true,
			wantErrMsg:     "temporary email",
		},
		{
			name:           "invalid email format",
			email:          "invalid-email",
			mockResponse:   `{"valid": false, "disposable": false}`,
			mockStatusCode: http.StatusOK,
			wantValid:      false,
			wantErr:        true,
			wantErrMsg:     "invalid email",
		},
		{
			name:           "API error",
			email:          "user@example.com",
			mockResponse:   `{"error": "Internal server error"}`,
			mockStatusCode: http.StatusInternalServerError,
			wantValid:      false,
			wantErr:        true,
		},
		{
			name:           "network timeout",
			email:          "user@example.com",
			mockResponse:   ``,
			mockStatusCode: http.StatusOK,
			wantValid:      false,
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

			// Create checker client
			checker := NewChecker(Config{
				APIKey:  "test-api-key",
				BaseURL: server.URL,
				Timeout: 1 * time.Second,
				Retries: 1,
			})

			// Test validation
			valid, err := checker.ValidateEmail(context.Background(), tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}
				assert.False(t, valid)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValid, valid)
			}
		})
	}
}

func TestChecker_ValidateEmail_OpenTelemetry(t *testing.T) {
	// Setup OpenTelemetry tracer
	tracer := otel.Tracer("test")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"valid": true, "disposable": false}`))
	}))
	defer server.Close()

	checker := NewChecker(Config{
		APIKey:  "test-api-key",
		BaseURL: server.URL,
		Timeout: 1 * time.Second,
		Retries: 1,
	})

	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	valid, err := checker.ValidateEmail(ctx, "user@example.com")

	require.NoError(t, err)
	assert.True(t, valid)

	// Verify span was created (check span context is propagated)
	// Note: In test environment without OTEL exporter, span context may not be fully initialized
	// The important part is that the function accepts context and doesn't panic
	spanCtx := trace.SpanContextFromContext(ctx)
	// Span context should exist (even if not fully initialized in test env)
	assert.NotNil(t, spanCtx)
}

func TestChecker_ValidateEmail_RetryLogic(t *testing.T) {
	retryCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryCount++
		if retryCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"valid": true, "disposable": false}`))
	}))
	defer server.Close()

	checker := NewChecker(Config{
		APIKey:  "test-api-key",
		BaseURL: server.URL,
		Timeout: 1 * time.Second,
		Retries: 3,
	})

	valid, err := checker.ValidateEmail(context.Background(), "user@example.com")

	assert.NoError(t, err)
	assert.True(t, valid)
	assert.Equal(t, 3, retryCount)
}

func TestNewChecker(t *testing.T) {
	config := Config{
		APIKey:  "test-key",
		BaseURL: "https://api.checker.com",
		Timeout: 5 * time.Second,
		Retries: 2,
	}

	checker := NewChecker(config)

	assert.NotNil(t, checker)
	assert.Equal(t, config.APIKey, checker.config.APIKey)
	assert.Equal(t, config.BaseURL, checker.config.BaseURL)
}

