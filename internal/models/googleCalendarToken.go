package models

import (
	"time"

	"gorm.io/gorm"
)

// GoogleCalendarToken : represents the schema of `google_calendar_tokens` table
type GoogleCalendarToken struct {
	gorm.Model
	UserID       uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	AccessToken  string    `gorm:"type:text;not null" json:"access_token"`
	RefreshToken string    `gorm:"type:text;not null" json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}