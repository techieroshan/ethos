# Email Services Integration Guide

## Overview

The Ethos BFF integrates three email services:
- **Checker**: Email validation (prevents temporary/disposable emails)
- **Emailit**: Production email sending (templated emails)
- **Mailpit**: Local email testing (development)

## Configuration

### Checker (Email Validation)

```bash
CHECKER_API_KEY=your-api-key
CHECKER_BASE_URL=https://api.checker.com
CHECKER_TIMEOUT=5s
CHECKER_RETRIES=2
```

### Emailit (Production Email Sending)

```bash
EMAILIT_API_KEY=your-api-key
EMAILIT_BASE_URL=https://api.emailit.com
EMAILIT_TIMEOUT=10s
EMAILIT_RETRIES=2
```

### Mailpit (Local Development)

```bash
MAILPIT_ENABLED=true
MAILPIT_SMTP_HOST=localhost
MAILPIT_SMTP_PORT=1025
MAILPIT_FROM_EMAIL=noreply@ethos.test
```

## Usage

### Email Validation

Email validation is automatically integrated into user registration:

```go
// In AuthService.Register()
if s.emailChecker != nil {
    valid, err := s.emailChecker.ValidateEmail(ctx, req.Email)
    // Rejects temporary/disposable emails
}
```

### Email Sending

Email sending is integrated into registration flow:

```go
// Verification email sent automatically on registration
if s.emailSender != nil {
    emailReq := email.SendEmailRequest{
        To:           req.Email,
        Subject:      "Verify Your Email",
        TemplateID:   templates.TemplateEmailVerification,
        TemplateData: map[string]interface{}{"name": req.Name},
    }
    s.emailSender.SendEmail(ctx, emailReq)
}
```

## Email Templates

Templates are defined in `pkg/email/templates/templates.go`:

- `TemplateEmailVerification` - Email verification
- `TemplatePasswordReset` - Password reset
- `TemplateAccountDeletion` - Account deletion confirmation
- `TemplateSecurityAlert` - Security alerts

## Local Development with Mailpit

1. Start Mailpit:
```bash
docker run -d -p 1025:1025 -p 8025:8025 axllent/mailpit
```

2. Configure `.env`:
```bash
MAILPIT_ENABLED=true
MAILPIT_SMTP_HOST=localhost
MAILPIT_SMTP_PORT=1025
```

3. View emails at `http://localhost:8025`

## Production Setup

1. Get Emailit API key
2. Configure templates in Emailit dashboard
3. Set `EMAILIT_API_KEY` in production environment
4. Ensure `MAILPIT_ENABLED=false`

## Error Handling

- Email validation failures return `VALIDATION_FAILED` error
- Email sending failures are logged but don't fail registration
- Retry logic handles transient failures
- OpenTelemetry traces all email operations

## Testing

Email services have comprehensive test coverage:
- Unit tests for each client
- Integration tests for registration flow
- Mock servers for API testing

