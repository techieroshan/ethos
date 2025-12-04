package templates

import (
	"strings"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		name        string
		templateID  string
		data        map[string]interface{}
		expectError bool
		contains    []string
	}{
		{
			name:       "email verification template",
			templateID: TemplateEmailVerification,
			data: map[string]interface{}{
				"Name":      "John Doe",
				"VerifyURL": "https://example.com/verify/123",
			},
			expectError: false,
			contains:    []string{"Verify Your Email", "John Doe", "Ethos Platform", "1320 Pepperhill Ln"},
		},
		{
			name:       "password reset template",
			templateID: TemplatePasswordReset,
			data: map[string]interface{}{
				"Name":     "Jane Smith",
				"ResetURL": "https://example.com/reset/abc123",
			},
			expectError: false,
			contains:    []string{"Reset Your Password", "Jane Smith", "Ethos Platform", "1320 Pepperhill Ln"},
		},
		{
			name:       "welcome standard user template",
			templateID: TemplateWelcomeStandardUser,
			data: map[string]interface{}{
				"Name":        "Test User",
				"Role":        "Standard User",
				"DashboardURL": "https://example.com/dashboard",
			},
			expectError: false,
			contains:    []string{"Welcome to Ethos", "Test User", "Ethos Platform", "1320 Pepperhill Ln"},
		},
		{
			name:       "account deletion template",
			templateID: TemplateAccountDeletion,
			data: map[string]interface{}{
				"Name": "Delete Me",
			},
			expectError: false,
			contains:    []string{"Account Deletion", "Delete Me", "Ethos Platform", "1320 Pepperhill Ln"},
		},
		{
			name:       "invalid template ID",
			templateID: "invalid_template",
			data:       map[string]interface{}{},
			expectError: true,
			contains:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RenderTemplate(tt.templateID, tt.data)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Check that result contains expected strings
			for _, contain := range tt.contains {
				if !strings.Contains(result, contain) {
					t.Errorf("expected result to contain '%s', but it didn't", contain)
				}
			}

			// Check CANSPAM compliance - should contain physical address
			if !strings.Contains(result, "1320 Pepperhill Ln") {
				t.Error("template should contain physical mailing address for CANSPAM compliance")
			}

			// Check for unsubscribe link
			if !strings.Contains(result, "unsubscribe") {
				t.Error("template should contain unsubscribe link for CANSPAM compliance")
			}

			// Should be HTML content
			if !strings.Contains(result, "<!DOCTYPE html>") {
				t.Error("template should render as HTML")
			}
		})
	}
}

func TestGetTemplate(t *testing.T) {
	tests := []struct {
		name       string
		templateID string
		expected   map[string]interface{}
	}{
		{
			name:       "email verification template config",
			templateID: TemplateEmailVerification,
			expected: map[string]interface{}{
				"subject":     "Verify Your Email Address - Ethos Platform",
				"template_id": TemplateEmailVerification,
			},
		},
		{
			name:       "password reset template config",
			templateID: TemplatePasswordReset,
			expected: map[string]interface{}{
				"subject":     "Reset Your Password - Ethos Platform",
				"template_id": TemplatePasswordReset,
			},
		},
		{
			name:       "unknown template",
			templateID: "unknown",
			expected: map[string]interface{}{
				"subject":     "Notification from Ethos Platform",
				"template_id": "unknown",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTemplate(tt.templateID)

			if result["subject"] != tt.expected["subject"] {
				t.Errorf("expected subject %s, got %s", tt.expected["subject"], result["subject"])
			}

			if result["template_id"] != tt.expected["template_id"] {
				t.Errorf("expected template_id %s, got %s", tt.expected["template_id"], result["template_id"])
			}
		})
	}
}
