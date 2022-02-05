package models

import (
	"time"

	"gorm.io/gorm"
)

// Token : represents the schema of `tokens` table
type Token struct {
	gorm.Model
	UserID       uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	AccessToken  string    `gorm:"type:text;not null" json:"access_token"`
	RefreshToken string    `gorm:"type:text;not null" json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}