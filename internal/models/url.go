package models

import (
	"time"
)

type URL struct {
	ID          int64      `db:"id" json:"id"`
	UserID      string     `db:"user_id" json:"user_id"`
	ShortCode   string     `db:"short_code" json:"short_code"`
	OriginalURL string     `db:"original_url" json:"original_url"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	ClickCount  int64      `db:"click_count" json:"click_count"`
	ExpiresAt   *time.Time `db:"expires_at" json:"expires_at,omitempty"`
}
