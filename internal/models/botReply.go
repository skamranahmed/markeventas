package models

import "gorm.io/gorm"

type BotReply struct {
	gorm.Model
	ToTweetID             int64  `gorm:"not null" json:"to_tweet_id"`
	UserTwitterID         string `gorm:"not null" json:"user_twitter_id"`
	UserTwitterScreenName string `gorm:"not null" json:"user_twitter_screen_name"`
	Body                  string `gorm:"type:text;default:null" json:"body"`
	Response              string `gorm:"type:text;default:null" json:"response"`
	StatusCode            int    `json:"status_code"`
	TypeCode              int    `json:"type_code"`
}
