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
	ID            int          `json:"id"`
	TrackingCode  string       `json:"tracking_code"`
	Village       *string      `json:"village"`
	SubDistrictId *int         `json:"sub_district_id"`
	Fullarea      string       `json:"full_area"`
	NameStatus    string       `json:"name_status"`
	Files         []ReportFile `json:"files"`
	Details       string       `json:"details"`
	Status        string       `json:"status"`
	CreatedAt     *time.Time   `json:"created_at"`
	UpdatedAt     *time.Time   `json:"updated_at"`
}

type SendReportRequest struct {
	Details       string `json:"details"`
	SubDistrictId int    `json:"sub_district_id"`
	Village       string `json:"village"`
}
