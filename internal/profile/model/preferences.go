package model

// UserPreferences represents user-specific preferences
type UserPreferences struct {
	UserID        string `json:"user_id"`
	NotifyOnLogin bool   `json:"notify_on_login"`
	Locale        string `json:"locale"`
}

// UpdatePreferencesRequest represents a preferences update request
type UpdatePreferencesRequest struct {
	NotifyOnLogin *bool   `json:"notify_on_login"`
	Locale        *string `json:"locale"`
}

