package profile

import (
	"context"
	"time"

	authModel "ethos/internal/auth/model"
)

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	Name      string `json:"name"`
	PublicBio string `json:"public_bio"`
}

// UpdatePreferencesRequest represents a preferences update request
type UpdatePreferencesRequest struct {
	NotifyOnLogin *bool   `json:"notify_on_login"`
	Locale        *string `json:"locale"`
}

// OptOutRequest represents an opt-out request
type OptOutRequest struct {
	From   string `json:"from" binding:"required,oneof=public_search analytics_use"`
	Reason string `json:"reason,omitempty"`
}

// AnonymizeResponse represents the response from anonymization
type AnonymizeResponse struct {
	Status             string    `json:"status"`
	ExpectedCompletion time.Time `json:"expected_completion"`
}

// DeleteRequest represents a deletion request
type DeleteRequest struct {
	Confirm bool   `json:"confirm" binding:"required"`
	Reason  string `json:"reason,omitempty"`
}

// DeleteResponse represents the response from deletion request
type DeleteResponse struct {
	Status             string    `json:"status"`
	ExpectedCompletion time.Time `json:"expected_completion"`
}

// PreferencesUpdate represents a request to update user preferences
type PreferencesUpdate struct {
	Notifications *NotificationPreferences `json:"notifications,omitempty"`
	Privacy       *PrivacyPreferences      `json:"privacy,omitempty"`
	Display       *DisplayPreferences      `json:"display,omitempty"`
}

// NotificationPreferences represents user notification preferences
type NotificationPreferences struct {
	EmailDigest     *bool `json:"email_digest,omitempty"`
	PushEnabled     *bool `json:"push_enabled,omitempty"`
	SmsEnabled      *bool `json:"sms_enabled,omitempty"`
	MarketingEmails *bool `json:"marketing_emails,omitempty"`
}

// PrivacyPreferences represents user privacy preferences
type PrivacyPreferences struct {
	ProfileVisibility *string `json:"profile_visibility,omitempty"` // "public", "private", "org"
	DataSharing       *bool   `json:"data_sharing,omitempty"`
	AnalyticsOptOut   *bool   `json:"analytics_opt_out,omitempty"`
}

// DisplayPreferences represents user display preferences
type DisplayPreferences struct {
	Theme       *string `json:"theme,omitempty"`        // "light", "dark", "auto"
	Language    *string `json:"language,omitempty"`     // ISO language code
	Timezone    *string `json:"timezone,omitempty"`     // IANA timezone
	DateFormat  *string `json:"date_format,omitempty"`  // date format preference
	TimeFormat  *string `json:"time_format,omitempty"`  // 12h or 24h
}

// Repository defines the interface for profile data access
type Repository interface {
	// GetProfile retrieves a user's profile
	GetProfile(ctx context.Context, userID string) (*authModel.UserProfile, error)

	// UpdateProfile updates a user's profile
	UpdateProfile(ctx context.Context, userID string, profile *authModel.UserProfile) error

	// UpdatePreferences updates user preferences
	UpdatePreferences(ctx context.Context, userID string, prefs *PreferencesUpdate) error

	// DeleteProfile deletes a user profile
	DeleteProfile(ctx context.Context, userID string) error

	// OptOut opts a user out of data collection
	OptOut(ctx context.Context, userID string, req *OptOutRequest) error

	// Anonymize anonymizes user data for privacy
	Anonymize(ctx context.Context, userID string) (*AnonymizeResponse, error)

	// RequestDeletion requests deletion of user data
	RequestDeletion(ctx context.Context, userID string, req *DeleteRequest) (*DeleteResponse, error)
}
