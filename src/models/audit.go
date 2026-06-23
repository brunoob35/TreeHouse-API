package models

import "time"

type AuditLogFieldChange struct {
	Field  string      `json:"field"`
	Before interface{} `json:"before,omitempty"`
	After  interface{} `json:"after,omitempty"`
}

type AuditLogItem struct {
	ID          uint64                 `json:"id"`
	UserID      *uint64                `json:"user_id,omitempty"`
	UserName    string                 `json:"user_name"`
	TableName   string                 `json:"table_name"`
	RecordID    string                 `json:"record_id"`
	Action      string                 `json:"action"`
	Description string                 `json:"description,omitempty"`
	Summary     string                 `json:"summary"`
	ChangesCount int                   `json:"changes_count"`
	Changes     []AuditLogFieldChange  `json:"changes"`
	Before      map[string]interface{} `json:"before,omitempty"`
	After       map[string]interface{} `json:"after,omitempty"`
	IPOrigin    string                 `json:"ip_origin,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

type AuditLogPage struct {
	Items      []AuditLogItem `json:"items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      uint64         `json:"total"`
	TotalPages int            `json:"total_pages"`
}
