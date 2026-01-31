package models

import "time"

type ReportFile struct {
	ID        int64     `json:"id"`
	FileName  string    `json:"file_name"`
	FileType  string    `json:"file_type"`
	FileSize  int64     `json:"file_size"`
	ObjectKey string    `json:"object_key"`
	CreatedAt time.Time `json:"created_at"`
	StreamURL string    `json:"stream_url,omitempty"`
}

type NacorticsReport struct {
	ID           int          `json:"id"`
	TrackingCode string       `json:"tracking_code"`
	NameStatus   string       `json:"name_status"`
	Files        []ReportFile `json:"files"`
	Details      string       `json:"details"`
	Status       string       `json:"status"`
	CreatedAt    *time.Time   `json:"created_at"`
	UpdatedAt    *time.Time   `json:"updated_at"`
}
