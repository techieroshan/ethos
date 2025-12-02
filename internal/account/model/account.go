package model

import "time"

// SecurityEvent represents a security event
type SecurityEvent struct {
	EventID   string    `json:"event_id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	Location  string    `json:"location"`
}

// DataExport represents a data export request
type DataExport struct {
	ExportID    string     `json:"export_id"`
	Status      string     `json:"status"`
	DownloadURL string     `json:"download_url,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

