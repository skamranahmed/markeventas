package models

import (
	"gorm.io/gorm"
)

// User : represents the schema of `users` table
type User struct {
	gorm.Model
	TwitterID              string `gorm:"uniqueIndex;not null" json:"twitter_id"`
	TwitterScreenName      string `gorm:"not null" json:"twitter_screen_name"`
	IsGcalOauthTokenActive bool   `gorm:"not null;default:false" json:"is_gcal_oauth_token_active"` // to maintain the status of google calendar oauth token

	// has one relationship
	GoogleCalendarToken Token
}
