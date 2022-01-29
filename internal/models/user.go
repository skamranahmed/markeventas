package models

import (
	"gorm.io/gorm"
)

// User : represents the db model of `users` table
type User struct {
	gorm.Model
	TwitterID         string `gorm:"uniqueIndex" json:"twitter_id"`
	TwitterScreenName string `gorm:"uniqueIndex" json:"twitter_screen_name"`
}
