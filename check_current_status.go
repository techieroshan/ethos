package main

import (
	"fmt"
	"strings"

	"ethos/pkg/email/templates"
)

func main() {
	fmt.Println("🔍 CHECKING CURRENT EMAIL COMPATIBILITY STATUS")
	fmt.Println("==============================================")

	templatesToCheck := []struct {
		id       string
		name     string
		category string
	}{
		{templates.TemplateEmailVerification, "Email Verification", "Core Auth"},
		{templates.TemplatePasswordReset, "Password Reset", "Core Auth"},
		{templates.TemplateWelcomeStandardUser, "Welcome Standard User", "Welcome"},
		{templates.TemplateWelcomeOrganizationAdmin, "Welcome Organization Admin", "Welcome"},
		{templates.TemplateWelcomePlatformAdmin, "Welcome Platform Admin", "Welcome"},
		{templates.TemplateWelcomeModerator, "Welcome Community Moderator", "Welcome"},
		{templates.TemplateWelcomeSupportEngineer, "Welcome Support Engineer", "Welcome"},
		{templates.TemplateWelcomeComplianceAuditor, "Welcome Compliance Auditor", "Welcome"},
		{templates.TemplateFeedbackSubmitted, "Feedback Submitted", "Feedback"},
		{templates.TemplateFeedbackModerated, "Feedback Moderated", "Feedback"},
		{templates.TemplateAppealSubmitted, "Appeal Submitted", "Appeals"},
		{templates.TemplateAppealStatusUpdate, "Appeal Status Update", "Appeals"},
		{templates.TemplateAppealResolved, "Appeal Resolved", "Appeals"},
		{templates.TemplateOrgInvitation, "Organization Invitation", "Organizations"},
		{templates.TemplateEscalationReceived, "Escalation Received", "Alerts"},
		{templates.TemplateModerationAlert, "Moderation Alert", "Alerts"},
		{templates.TemplateSecurityAlert, "Security Alert", "Alerts"},
		{templates.TemplateAuditReportAvailable, "Audit Report Available", "Reports"},
	}

	fmt.Printf("%-30s %-15s %-10s\n", "TEMPLATE NAME", "CATEGORY", "SCORE")
	fmt.Println(strings.Repeat("=", 60))

	for _, tmpl := range templatesToCheck {
		html, err := templates.RenderTemplate(tmpl.id, map[string]interface{}{
			"Name":         "Test User",
			"Email":        "test@example.com",
			"VerifyURL":    "https://example.com/verify",
			"ResetURL":     "https://example.com/reset",
			"DashboardURL": "https://example.com/dashboard",
			"FeedbackTitle": "Test Feedback",
			"FeedbackID":   "FB-TEST-001",
		})

		if err != nil {
			fmt.Printf("%-30s %-15s ❌ ERROR\n", tmpl.name, tmpl.category)
			continue
		}

		score := calculateRealCompatibilityScore(html)
		percentage := (score * 100) / 100

		status := "❌ FAIL"
		if percentage >= 95 {
			status = "✅ PASS"
		} else if percentage >= 85 {
			status = "⚠️  OK"
		}

		fmt.Printf("%-30s %-15s %3d%% %s\n", tmpl.name, tmpl.category, percentage, status)
	}

	fmt.Println("\n🔧 ANALYSIS:")
	fmt.Println("   Templates need complete refactoring to achieve 95%+ compatibility")
	fmt.Println("   Current scores show they still have unsupported elements")
}

func calculateRealCompatibilityScore(html string) int {
	score := 0

	// HTML Structure (20 points)
	if strings.Contains(html, "<!DOCTYPE html>") { score += 2 }
	if strings.Contains(html, "<html") && strings.Contains(html, "</html>") { score += 2 }
	if strings.Contains(html, "<head>") && strings.Contains(html, "</head>") { score += 2 }
	if strings.Contains(html, "<body>") && strings.Contains(html, "</body>") { score += 2 }
	if strings.Contains(html, `charset="UTF-8"`) { score += 2 }
	if strings.Contains(html, `name="viewport"`) { score += 2 }
	if strings.Contains(html, "<title>") && strings.Contains(html, "</title>") { score += 2 }
	if strings.Contains(html, `lang="`) { score += 2 }
	if strings.Contains(html, "1320 Pepperhill Ln") { score += 2 } // CANSPAM

	// Table-based layout (30 points) - Critical for Outlook
	if strings.Contains(html, "<table") && strings.Contains(html, "<tr") && strings.Contains(html, "<td") {
		score += 15 // Has table structure
	}
	if strings.Contains(html, `width="100%"`) { score += 5 }
	if strings.Contains(html, `border="0"`) { score += 5 }
	if strings.Contains(html, `cellspacing="0"`) { score += 5 }

	// CSS Compatibility (25 points)
	if !strings.Contains(html, "display:flex") && !strings.Contains(html, "display: flex") { score += 5 }
	if !strings.Contains(html, "display:grid") && !strings.Contains(html, "display: grid") { score += 5 }
	if !strings.Contains(html, "linear-gradient") { score += 5 }
	if !strings.Contains(html, "transform:") { score += 5 }
	if strings.Contains(html, "Arial, sans-serif") || strings.Contains(html, "font-family: Arial") { score += 5 }

	// Responsive & Mobile (15 points)
	if strings.Contains(html, "max-width:") || strings.Contains(html, "max-width") { score += 5 }
	if strings.Contains(html, "margin: 0 auto") { score += 5 }
	if strings.Contains(html, "font-size:") { score += 5 }

	// Email Client Optimizations (10 points)
	if strings.Contains(html, `style="`) { score += 3 } // Inline styles
	if !strings.Contains(html, "position:") { score += 2 } // No absolute positioning
	if !strings.Contains(html, "float:") { score += 2 } // No floats
	if strings.Contains(html, "text-decoration: none") { score += 3 } // Proper link styling

	return score
}
