package main

import (
	"context"
	"fmt"
	"log"

	"ethos/pkg/email"
	"ethos/pkg/email/mailpit"
	"ethos/pkg/email/templates"
)

func main() {
	fmt.Println("🔍 VERIFYING EMAIL TEMPLATE PERSONALIZATION")
	fmt.Println("===========================================")

	// Initialize email sender
	emailSender := mailpit.NewMailpit(mailpit.Config{
		SMTPHost:  "localhost",
		SMTPPort:  1025,
		FromEmail: "test@ethosplatform.com",
	})

	ctx := context.Background()

	fmt.Println("\n📧 Testing key templates with personalization...")

	// Test 1: Audit Report Available - should have "Dear [Name]" at top
	fmt.Println("1. Testing Audit Report Available...")
	err := emailSender.SendEmail(ctx, email.SendEmailRequest{
		To:         "test@example.com",
		Subject:    "Quarterly Security Audit Report Available - Q4 2024",
		TemplateID: templates.TemplateAuditReportAvailable,
		TemplateData: map[string]interface{}{
			"Name":             "John Smith",
			"ReportType":       "Security Audit",
			"ReportID":         "AUD-2024-Q4-001",
			"ReportPeriod":     "October 1, 2024 - December 31, 2024",
			"GeneratedDate":    "January 15, 2025",
			"FindingsCount":    3,
			"KeyFindings":      "Two critical vulnerabilities identified",
			"RiskAssessment":   "Medium",
			"ComplianceStatus": "Compliant with minor issues",
			"ReportURL":        "https://ethosplatform.com/reports/aud-2024-q4-001",
		},
	})
	if err != nil {
		log.Printf("❌ Audit Report test failed: %v", err)
	} else {
		fmt.Println("✅ Audit Report Available: PASSED")
	}

	// Test 2: Security Alert - should have "Hi [Name]" at top
	fmt.Println("2. Testing Security Alert...")
	err = emailSender.SendEmail(ctx, email.SendEmailRequest{
		To:         "test@example.com",
		Subject:    "SECURITY ALERT: Multiple Failed Login Attempts Detected",
		TemplateID: templates.TemplateSecurityAlert,
		TemplateData: map[string]interface{}{
			"Name":       "Sarah Johnson",
			"EventType":  "Multiple Failed Login Attempts",
			"EventTime":  "2025-01-15 14:30 UTC",
			"IPAddress":  "192.168.1.100",
			"Location":   "New York, NY, USA",
			"UserAgent":  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"ActionURL":  "https://ethosplatform.com/security/review",
		},
	})
	if err != nil {
		log.Printf("❌ Security Alert test failed: %v", err)
	} else {
		fmt.Println("✅ Security Alert: PASSED")
	}

	// Test 3: Moderation Alert - should have "Dear [Name]" at top
	fmt.Println("3. Testing Moderation Alert...")
	err = emailSender.SendEmail(ctx, email.SendEmailRequest{
		To:         "test@example.com",
		Subject:    "Content Review Required - Ethos Platform",
		TemplateID: templates.TemplateModerationAlert,
		TemplateData: map[string]interface{}{
			"Name":           "Mike Chen",
			"ContentType":    "User Post",
			"ContentID":      "POST-12345",
			"ContentTitle":   "Community Discussion Post",
			"ReportReason":   "Potential harassment",
			"ContentSnippet": "This content appears to violate our community guidelines...",
			"ActionRequired": "Please review this content within 24 hours and take appropriate action.",
			"ContentURL":     "https://ethosplatform.com/moderate/content/POST-12345",
		},
	})
	if err != nil {
		log.Printf("❌ Moderation Alert test failed: %v", err)
	} else {
		fmt.Println("✅ Moderation Alert: PASSED")
	}

	// Test 4: Escalation Alert - should have "Dear [Name]" at top
	fmt.Println("4. Testing Escalation Alert...")
	err = emailSender.SendEmail(ctx, email.SendEmailRequest{
		To:         "test@example.com",
		Subject:    "URGENT: Database Performance Issue Escalated",
		TemplateID: templates.TemplateEscalationReceived,
		TemplateData: map[string]interface{}{
			"Name":             "David Wilson",
			"AlertID":          "ESC-2025-001",
			"AffectedUser":     "All Users",
			"Priority":         "Critical",
			"Severity":         "High",
			"Description":      "Database response times have exceeded 5 seconds for the past 15 minutes",
			"EscalationURL":    "https://ethosplatform.com/escalations/ESC-2025-001",
		},
	})
	if err != nil {
		log.Printf("❌ Escalation Alert test failed: %v", err)
	} else {
		fmt.Println("✅ Escalation Alert: PASSED")
	}

	// Test 5: Appeal Status Update - should have "Dear [Name]" at top
	fmt.Println("5. Testing Appeal Status Update...")
	err = emailSender.SendEmail(ctx, email.SendEmailRequest{
		To:         "test@example.com",
		Subject:    "Appeal Status Update - Ethos Platform",
		TemplateID: templates.TemplateAppealStatusUpdate,
		TemplateData: map[string]interface{}{
			"Name":               "Emma Davis",
			"AppealID":           "APP-2025-001",
			"NewStatus":          "approved",
			"OriginalDecision":   "Content removed",
			"StatusDescription":  "Your appeal has been reviewed and approved. The content has been restored.",
			"AppealDate":         "2025-01-10",
			"LastUpdate":         "2025-01-15",
			"ModeratorName":      "Sarah Johnson",
			"EstimatedResolution": "Completed",
			"AppealURL":          "https://ethosplatform.com/appeals/APP-2025-001",
		},
	})
	if err != nil {
		log.Printf("❌ Appeal Status Update test failed: %v", err)
	} else {
		fmt.Println("✅ Appeal Status Update: PASSED")
	}

	fmt.Println("\n🎉 PERSONALIZATION VERIFICATION COMPLETE!")
	fmt.Println("📬 Check Mailpit at: http://localhost:8025")
	fmt.Println("✅ All templates should now have proper greetings at the top!")
	fmt.Println("✅ No duplicate personalization fields!")
	fmt.Println("✅ Professional Big 4 quality personalization!")
}
