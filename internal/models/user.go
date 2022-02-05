package models

import (
	"gorm.io/gorm"
)

// User : represents the schema of `users` table
type User struct {
	gorm.Model
	TwitterID         string `gorm:"uniqueIndex;not null" json:"twitter_id"`
	TwitterScreenName string `gorm:"not null" json:"twitter_screen_name"`

	// has one relationship
	GoogleCalendarToken Token
}
