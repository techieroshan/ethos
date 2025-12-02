package templates

// EmailTemplate represents an email template
type EmailTemplate struct {
	ID      string
	Subject string
	Body    string
}

// Predefined email templates
const (
	TemplateEmailVerification = "email_verification"
	TemplatePasswordReset     = "password_reset"
	TemplateAccountDeletion   = "account_deletion"
	TemplateSecurityAlert     = "security_alert"
)

// GetTemplate returns template configuration for Emailit
func GetTemplate(templateID string) map[string]interface{} {
	templates := map[string]map[string]interface{}{
		TemplateEmailVerification: {
			"subject": "Verify Your Email Address",
			"template_id": "email_verification",
		},
		TemplatePasswordReset: {
			"subject": "Reset Your Password",
			"template_id": "password_reset",
		},
		TemplateAccountDeletion: {
			"subject": "Account Deletion Confirmation",
			"template_id": "account_deletion",
		},
		TemplateSecurityAlert: {
			"subject": "Security Alert",
			"template_id": "security_alert",
		},
	}

	if template, ok := templates[templateID]; ok {
		return template
	}

	return map[string]interface{}{
		"subject":     "Notification",
		"template_id": templateID,
	}
}

