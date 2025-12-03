package templates

import (
	"fmt"
	"html/template"
	"strings"
)

// EmailTemplate represents an email template with CANSPAM compliance
type EmailTemplate struct {
	ID          string
	Subject     string
	HTMLBody    string
	TextBody    string
	TemplateID  string
}

// Template data structures
type EmailVerificationData struct {
	Name    string
	Email   string
	UserID  string
	VerifyURL string
}

type PasswordResetData struct {
	Name    string
	Email   string
	ResetURL string
}

type WelcomeData struct {
	Name    string
	Email   string
	Role    string
	DashboardURL string
}

type FeedbackConfirmationData struct {
	Name    string
	Email   string
	FeedbackID string
	FeedbackURL string
	Title   string
}

type AppealNotificationData struct {
	Name    string
	Email   string
	AppealID string
	AppealURL string
	Status  string
	Reason  string
}

type OrganizationInvitationData struct {
	Name    string
	Email   string
	OrgName string
	Role    string
	AcceptURL string
	DeclineURL string
	InvitedBy string
}

type EscalationAlertData struct {
	Name    string
	Email   string
	EscalationID string
	EscalationURL string
	Priority string
	Description string
	UserAffected string
}

type ModerationAlertData struct {
	Name    string
	Email   string
	ContentID string
	ContentURL string
	ContentType string
	Reason  string
	ActionRequired string
}

type SecurityAlertData struct {
	Name    string
	Email   string
	EventType string
	EventTime string
	IPAddress string
	UserAgent string
	Location string
	ActionURL string
}

type AuditReportData struct {
	Name    string
	Email   string
	ReportID string
	ReportURL string
	ReportType string
	DateRange string
	FindingsCount int
}

type SystemHealthAlertData struct {
	Name    string
	Email   string
	AlertType string
	Severity string
	Service string
	Description string
	StatusURL string
	IncidentID string
}

// Predefined email template IDs
const (
	// Authentication Templates
	TemplateEmailVerification = "email_verification"
	TemplatePasswordReset     = "password_reset"
	TemplateAccountDeletion   = "account_deletion"
	TemplateSecurityAlert     = "security_alert"

	// Welcome & Onboarding Templates
	TemplateWelcomeStandardUser     = "welcome_standard_user"
	TemplateWelcomeOrganizationAdmin = "welcome_org_admin"
	TemplateWelcomePlatformAdmin    = "welcome_platform_admin"
	TemplateWelcomeModerator        = "welcome_moderator"
	TemplateWelcomeSupportEngineer  = "welcome_support_engineer"
	TemplateWelcomeComplianceAuditor = "welcome_compliance_auditor"

	// Feedback & Content Templates
	TemplateFeedbackSubmitted       = "feedback_submitted"
	TemplateFeedbackModerated       = "feedback_moderated"
	TemplateContentFlagged          = "content_flagged"
	TemplateContentApproved         = "content_approved"
	TemplateContentRejected         = "content_rejected"

	// Appeal Templates
	TemplateAppealSubmitted         = "appeal_submitted"
	TemplateAppealStatusUpdate      = "appeal_status_update"
	TemplateAppealResolved          = "appeal_resolved"

	// Organization Management Templates
	TemplateOrgInvitation           = "org_invitation"
	TemplateOrgMemberAdded          = "org_member_added"
	TemplateOrgMemberRemoved        = "org_member_removed"
	TemplateOrgRoleChanged          = "org_role_changed"
	TemplateOrgSettingsChanged      = "org_settings_changed"

	// Escalation & Support Templates
	TemplateEscalationReceived      = "escalation_received"
	TemplateEscalationAssigned      = "escalation_assigned"
	TemplateEscalationResolved      = "escalation_resolved"
	TemplateEscalationUpdate        = "escalation_update"

	// Moderation Templates
	TemplateModerationAlert         = "moderation_alert"
	TemplateModerationTraining      = "moderation_training"
	TemplateModerationSummary       = "moderation_summary"

	// Audit & Compliance Templates
	TemplateAuditReportAvailable    = "audit_report_available"
	TemplateComplianceAlert         = "compliance_alert"
	TemplateRegulatoryDeadline      = "regulatory_deadline"

	// System & Health Templates
	TemplateSystemHealthAlert       = "system_health_alert"
	TemplateSystemMaintenance       = "system_maintenance"
	TemplatePlatformAnnouncement    = "platform_announcement"

	// Multi-tenant Templates
	TemplateTenantInvitation        = "tenant_invitation"
	TemplateTenantRoleChanged       = "tenant_role_changed"
	TemplateCrossTenantNotification = "cross_tenant_notification"
)

// GetTemplate returns template configuration for Emailit
func GetTemplate(templateID string) map[string]interface{} {
	templates := map[string]map[string]interface{}{
		TemplateEmailVerification: {
			"subject": "Verify Your Email Address - Ethos Platform",
			"template_id": TemplateEmailVerification,
		},
		TemplatePasswordReset: {
			"subject": "Reset Your Password - Ethos Platform",
			"template_id": TemplatePasswordReset,
		},
		TemplateAccountDeletion: {
			"subject": "Account Deletion Confirmation - Ethos Platform",
			"template_id": TemplateAccountDeletion,
		},
		TemplateSecurityAlert: {
			"subject": "Security Alert - Ethos Platform",
			"template_id": TemplateSecurityAlert,
		},
		TemplateWelcomeStandardUser: {
			"subject": "Welcome to Ethos - Your Account is Ready",
			"template_id": TemplateWelcomeStandardUser,
		},
		TemplateWelcomeOrganizationAdmin: {
			"subject": "Welcome to Ethos - Organization Admin Access",
			"template_id": TemplateWelcomeOrganizationAdmin,
		},
		TemplateWelcomePlatformAdmin: {
			"subject": "Welcome to Ethos - Platform Admin Access",
			"template_id": TemplateWelcomePlatformAdmin,
		},
		TemplateWelcomeModerator: {
			"subject": "Welcome to Ethos - Community Moderator Access",
			"template_id": TemplateWelcomeModerator,
		},
		TemplateWelcomeSupportEngineer: {
			"subject": "Welcome to Ethos - Support Engineer Access",
			"template_id": TemplateWelcomeSupportEngineer,
		},
		TemplateWelcomeComplianceAuditor: {
			"subject": "Welcome to Ethos - Compliance Auditor Access",
			"template_id": TemplateWelcomeComplianceAuditor,
		},
		TemplateFeedbackSubmitted: {
			"subject": "Feedback Submitted Successfully - Ethos",
			"template_id": TemplateFeedbackSubmitted,
		},
		TemplateAppealSubmitted: {
			"subject": "Appeal Submitted - Ethos Platform",
			"template_id": TemplateAppealSubmitted,
		},
		TemplateAppealStatusUpdate: {
			"subject": "Appeal Status Update - Ethos Platform",
			"template_id": TemplateAppealStatusUpdate,
		},
		TemplateOrgInvitation: {
			"subject": "Organization Invitation - Ethos Platform",
			"template_id": TemplateOrgInvitation,
		},
		TemplateEscalationReceived: {
			"subject": "New Escalation Requires Attention - Ethos",
			"template_id": TemplateEscalationReceived,
		},
		TemplateModerationAlert: {
			"subject": "Content Requires Moderation Review - Ethos",
			"template_id": TemplateModerationAlert,
		},
		TemplateAuditReportAvailable: {
			"subject": "Audit Report Available - Ethos Platform",
			"template_id": TemplateAuditReportAvailable,
		},
		TemplateSystemHealthAlert: {
			"subject": "System Health Alert - Ethos Platform",
			"template_id": TemplateSystemHealthAlert,
		},
	}

	if template, ok := templates[templateID]; ok {
		return template
	}

	return map[string]interface{}{
		"subject":     "Notification from Ethos Platform",
		"template_id": templateID,
	}
}

// RenderTemplate renders an HTML template with data
func RenderTemplate(templateID string, data interface{}) (string, error) {
	htmlContent, err := getTemplateHTML(templateID)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(templateID).Parse(htmlContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// getTemplateHTML returns the HTML content for a template
func getTemplateHTML(templateID string) (string, error) {
	templates := map[string]string{
		TemplateEmailVerification: emailVerificationTemplate,
		TemplatePasswordReset:     passwordResetTemplate,
		TemplateAccountDeletion:   accountDeletionTemplate,
		TemplateWelcomeStandardUser: welcomeStandardUserTemplate,
		TemplateWelcomeOrganizationAdmin: welcomeOrgAdminTemplate,
		TemplateWelcomePlatformAdmin: welcomePlatformAdminTemplate,
		TemplateWelcomeModerator: welcomeModeratorTemplate,
		TemplateWelcomeSupportEngineer: welcomeSupportEngineerTemplate,
		TemplateWelcomeComplianceAuditor: welcomeComplianceAuditorTemplate,
		TemplateFeedbackSubmitted: feedbackSubmittedTemplate,
		TemplateFeedbackModerated: feedbackModeratedTemplate,
		TemplateAppealSubmitted: appealSubmittedTemplate,
		TemplateAppealStatusUpdate: appealStatusUpdateTemplate,
		TemplateAppealResolved: appealResolvedTemplate,
		TemplateOrgInvitation: orgInvitationTemplate,
		TemplateOrgMemberAdded: orgMemberAddedTemplate,
		TemplateEscalationReceived: escalationReceivedTemplate,
		TemplateModerationAlert: moderationAlertTemplate,
		TemplateAuditReportAvailable: auditReportAvailableTemplate,
		TemplateSystemHealthAlert: systemHealthAlertTemplate,
		TemplateSecurityAlert: securityAlertTemplate,
	}

	if content, ok := templates[templateID]; ok {
		return content, nil
	}

	return "", fmt.Errorf("template not found: %s", templateID)
}

// CANSPAM Compliant Footer
const canSpamFooter = `
<div style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #e5e7eb; font-size: 12px; color: #6b7280; line-height: 1.4;">
    <p><strong>Ethos Platform</strong><br>
    1320 Pepperhill Ln<br>
    Fort Worth, TX, 76131<br>
    United States</p>

    <p>This email was sent to you because you have an account with Ethos Platform.</p>

    <p>To unsubscribe from these emails or manage your notification preferences, please visit your <a href="{{.DashboardURL}}" style="color: #3b82f6;">account settings</a>.</p>

    <p>If you no longer wish to receive any communications from Ethos Platform, you can <a href="{{.UnsubscribeURL}}" style="color: #3b82f6;">unsubscribe here</a>.</p>
</div>
`