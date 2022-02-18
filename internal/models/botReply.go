package models

import "gorm.io/gorm"

type BotLog struct {
	gorm.Model
	ToTweetID             int64  `gorm:"not null" json:"to_tweet_id"`
	UserTwitterID         string `gorm:"not null" json:"user_twitter_id"`
	UserTwitterScreenName string `gorm:"not null" json:"user_twitter_screen_name"`
	UserTweet             string `gorm:"type:text;default:null" json:"user_tweet"`
	UserParsedTweet       string `gorm:"type:text;default:null" json:"user_parsed_tweet"`
	Reply                 string `gorm:"type:text;default:null" json:"reply"`
	HttpResponse          string `gorm:"type:text;default:null" json:"http_response"`
	HttpStatusCode        int    `json:"http_status_code"`
	ReplyTypeCode         int    `json:"reply_type_code"`
}
