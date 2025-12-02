package checker

import "context"

// EmailChecker defines the interface for email validation
type EmailChecker interface {
	ValidateEmail(ctx context.Context, email string) (bool, error)
}

