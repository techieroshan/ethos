package model

import "time"

// Organization represents an organization/tenant in the system
type Organization struct {
	ID          string
	Name        string
	Domain      string
	OwnerID     string
	Description string
	Status      string // "active", "suspended", "trial"
	Plan        string // "trial", "pro", "enterprise"
	MaxUsers    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SuspendedAt *time.Time
}

// OrganizationResponse represents an organization for API responses
type OrganizationResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Domain      string    `json:"domain"`
	OwnerID     string    `json:"owner_id"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	Plan        string    `json:"plan"`
	MaxUsers    int       `json:"max_users"`
	UserCount   int       `json:"user_count"`
	AdminCount  int       `json:"admin_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateOrganizationRequest represents a request to create an organization
type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Domain      string `json:"domain" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"max=500"`
	Plan        string `json:"plan" binding:"required,oneof=trial pro enterprise"`
}

// UpdateOrganizationRequest represents a request to update an organization
type UpdateOrganizationRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=255"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Status      string `json:"status" binding:"omitempty,oneof=active suspended"`
}

// OrganizationMember represents a member of an organization
type OrganizationMember struct {
	ID             string
	OrganizationID string
	UserID         string
	Role           string // "admin", "moderator", "user"
	JoinedAt       time.Time
	LastActiveAt   *time.Time
}

// OrganizationMemberResponse represents a member for API responses
type OrganizationMemberResponse struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	UserName     string     `json:"user_name"`
	UserEmail    string     `json:"user_email"`
	Role         string     `json:"role"`
	Status       string     `json:"status"` // "active", "suspended", "pending"
	JoinedAt     time.Time  `json:"joined_at"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
}

// AddMemberRequest represents a request to add a member to an organization
type AddMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=admin moderator user"`
}

// UpdateMemberRequest represents a request to update a member's role
type UpdateMemberRequest struct {
	Role string `json:"role" binding:"required,oneof=admin moderator user"`
}

// OrganizationSettings represents organization-wide settings
type OrganizationSettings struct {
	ID                       string
	OrganizationID           string
	RequireEmailVerification bool
	AllowPublicProfiles      bool
	EnableModeration         bool
	RequireApproval          bool
	DataRetentionDays        int
	UpdatedAt                time.Time
}

// UpdateSettingsRequest represents a request to update organization settings
type UpdateSettingsRequest struct {
	RequireEmailVerification *bool `json:"require_email_verification"`
	AllowPublicProfiles      *bool `json:"allow_public_profiles"`
	EnableModeration         *bool `json:"enable_moderation"`
	RequireApproval          *bool `json:"require_approval"`
	DataRetentionDays        *int  `json:"data_retention_days"`
}

// ADMIN MODELS - Platform-wide administration

// UserAdminResponse represents detailed user information for admin operations
type UserAdminResponse struct {
	ID            string            `json:"id"`
	Email         string            `json:"email"`
	Name          string            `json:"name"`
	EmailVerified bool              `json:"email_verified"`
	Status        string            `json:"status"` // "active", "suspended", "banned"
	Roles         []UserRole        `json:"roles"`
	LastLogin     *time.Time        `json:"last_login,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     *time.Time        `json:"updated_at,omitempty"`
	Organization  *OrganizationInfo `json:"organization,omitempty"`
}

// UserRole represents a role assigned to a user
type UserRole struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Permissions map[string]interface{} `json:"permissions"`
	AssignedAt  time.Time              `json:"assigned_at"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	IsActive    bool                   `json:"is_active"`
}

// OrganizationInfo represents basic organization information
type OrganizationInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SystemAnalytics represents system-wide analytics
type SystemAnalytics struct {
	TotalUsers       int64 `json:"total_users"`
	ActiveUsers      int64 `json:"active_users"`
	TotalOrganizations int64 `json:"total_organizations"`
	ActiveOrganizations int64 `json:"active_organizations"`
	TotalFeedback    int64 `json:"total_feedback"`
	PendingModeration int64 `json:"pending_moderation"`
	SystemHealth     string `json:"system_health"` // "healthy", "warning", "critical"
}

// UserAnalytics represents user-related analytics
type UserAnalytics struct {
	UserGrowth         []TimeSeriesPoint `json:"user_growth"`
	UserRetention      []RetentionPoint  `json:"user_retention"`
	UserActivity       []ActivityPoint   `json:"user_activity"`
	GeographicDistribution []GeoPoint    `json:"geographic_distribution"`
}

// ContentAnalytics represents content-related analytics
type ContentAnalytics struct {
	FeedbackGrowth     []TimeSeriesPoint `json:"feedback_growth"`
	ContentModeration  []ModerationStats `json:"content_moderation"`
	PopularCategories  []CategoryStats   `json:"popular_categories"`
	EngagementMetrics  EngagementStats   `json:"engagement_metrics"`
}

// TimeSeriesPoint represents a data point in a time series
type TimeSeriesPoint struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

// RetentionPoint represents user retention data
type RetentionPoint struct {
	Cohort   string  `json:"cohort"`
	Day0     int64   `json:"day_0"`
	Day7     int64   `json:"day_7"`
	Day30    int64   `json:"day_30"`
	Day90    int64   `json:"day_90"`
}

// ActivityPoint represents user activity data
type ActivityPoint struct {
	Date         string `json:"date"`
	ActiveUsers  int64  `json:"active_users"`
	NewUsers     int64  `json:"new_users"`
	ReturningUsers int64 `json:"returning_users"`
}

// GeoPoint represents geographic distribution data
type GeoPoint struct {
	Country string  `json:"country"`
	Users   int64   `json:"users"`
	Percent float64 `json:"percent"`
}

// ModerationStats represents content moderation statistics
type ModerationStats struct {
	Date            string `json:"date"`
	PendingContent  int64  `json:"pending_content"`
	ApprovedContent int64  `json:"approved_content"`
	RejectedContent int64  `json:"rejected_content"`
	EscalatedContent int64 `json:"escalated_content"`
}

// CategoryStats represents popular category statistics
type CategoryStats struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
	Percent  float64 `json:"percent"`
}

// EngagementStats represents engagement metrics
type EngagementStats struct {
	AverageLikes      float64 `json:"average_likes"`
	AverageComments   float64 `json:"average_comments"`
	AverageBookmarks  float64 `json:"average_bookmarks"`
	EngagementRate    float64 `json:"engagement_rate"`
}

// AuditLogEntry represents an audit log entry
type AuditLogEntry struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	ResourceID string   `json:"resource_id"`
	Details   string    `json:"details"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Timestamp time.Time `json:"timestamp"`
}

// SystemSettings represents system-wide settings
type SystemSettings struct {
	ID                           string                 `json:"id"`
	RequireEmailVerification    bool                   `json:"require_email_verification"`
	AllowPublicProfiles         bool                   `json:"allow_public_profiles"`
	EnableGlobalModeration      bool                   `json:"enable_global_moderation"`
	MaxFeedbackPerDay           int                    `json:"max_feedback_per_day"`
	MaxCommentsPerHour          int                    `json:"max_comments_per_hour"`
	DataRetentionDays           int                    `json:"data_retention_days"`
	EnableAnalytics             bool                   `json:"enable_analytics"`
	MaintenanceMode             bool                   `json:"maintenance_mode"`
	CustomSettings              map[string]interface{} `json:"custom_settings"`
	UpdatedAt                   time.Time              `json:"updated_at"`
	UpdatedBy                   string                 `json:"updated_by"`
}

// BulkOperationResult represents the result of a bulk operation
type BulkOperationResult struct {
	TotalRequested int                    `json:"total_requested"`
	Successful     int                    `json:"successful"`
	Failed         int                    `json:"failed"`
	Errors         []BulkOperationError   `json:"errors,omitempty"`
}

// BulkOperationError represents an error in a bulk operation
type BulkOperationError struct {
	UserID string `json:"user_id"`
	Error  string `json:"error"`
}
// ORGANIZATION ADMIN MODELS

// OrganizationAnalytics represents organization-specific analytics
type OrganizationAnalytics struct {
	TotalUsers       int64 `json:"total_users"`
	ActiveUsers      int64 `json:"active_users"`
	TotalFeedback    int64 `json:"total_feedback"`
	PendingModeration int64 `json:"pending_moderation"`
	OpenIncidents    int64 `json:"open_incidents"`
	UserGrowth       []TimeSeriesPoint `json:"user_growth"`
	ActivityLevel    string `json:"activity_level"` // low, medium, high
}

// OrganizationUserAnalytics represents organization user analytics
type OrganizationUserAnalytics struct {
	UserRetention      []RetentionPoint  `json:"user_retention"`
	UserActivity       []ActivityPoint   `json:"user_activity"`
	RoleDistribution   []RoleStats       `json:"role_distribution"`
	EngagementMetrics  EngagementStats   `json:"engagement_metrics"`
}

// OrganizationContentAnalytics represents organization content analytics
type OrganizationContentAnalytics struct {
	ContentGrowth      []TimeSeriesPoint `json:"content_growth"`
	ContentModeration  []ModerationStats `json:"content_moderation"`
	PopularCategories  []CategoryStats   `json:"popular_categories"`
	EngagementRate     float64           `json:"engagement_rate"`
}

// OrganizationUser represents a user within an organization context
type OrganizationUser struct {
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Status      string    `json:"status"` // active, suspended, invited
	JoinedAt    time.Time `json:"joined_at"`
	LastActive  *time.Time `json:"last_active,omitempty"`
	FeedbackCount int64   `json:"feedback_count"`
}

// OrganizationAuditEntry represents an audit entry for organization actions
type OrganizationAuditEntry struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	ResourceID string   `json:"resource_id"`
	Details   string    `json:"details"`
	IPAddress string    `json:"ip_address"`
	Timestamp time.Time `json:"timestamp"`
}

// OrganizationIncident represents an incident within an organization
type OrganizationIncident struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"` // open, investigating, resolved, closed
	Priority    string     `json:"priority"` // low, medium, high, critical
	Category    string     `json:"category"`
	AssignedTo  *string    `json:"assigned_to,omitempty"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	Resolution  *string    `json:"resolution,omitempty"`
}

// RoleStats represents role distribution statistics
type RoleStats struct {
	Role  string  `json:"role"`
	Count int64   `json:"count"`
	Percent float64 `json:"percent"`
}
