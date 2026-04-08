package models

import "time"

type PasswordReset struct {
	ID        uint64     `json:"id,omitempty"`
	UserID    uint64     `json:"user_id,omitempty"`
	TokenHash string     `json:"token_hash,omitempty"`
	ExpiresAt time.Time  `json:"expires_at,omitempty"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
}
