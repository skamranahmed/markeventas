package models

import (
	"github.com/skamranahmed/markeventas/internal/utils"
	"gorm.io/gorm"
)

type BotReplyType struct {
	gorm.Model
	Code int    `gorm:"uniqueIndex;not null" json:"code"`
	Name string `gorm:"not null" json:"name"`
	Body string `gorm:"not null" json:"body"`
}

var (
	// UserAccountDoesNotExistReply
	UserAccountDoesNotExistReply = BotReplyType{
		Code: 1,
		Name: string(utils.UserAccountDoesNotExist),
		Body: utils.UserAccountDoesNotExist.Body(),
	}

	// UserGoogleCalendarConsentAbsentReply
	UserGoogleCalendarConsentAbsentReply = BotReplyType{
		Code: 2,
		Name: string(utils.UserGoogleCalendarConsentAbsent),
		Body: utils.UserGoogleCalendarConsentAbsent.Body(),
	}

	// UserGoogleCalendarEventCreatedReply
	UserGoogleCalendarEventCreatedReply = BotReplyType{
		Code: 3,
		Name: string(utils.UserGoogleCalendarEventCreated),
		Body: utils.UserGoogleCalendarEventCreated.Body(),
	}

	BotReplyTypes = []BotReplyType{
		UserAccountDoesNotExistReply,
		UserGoogleCalendarConsentAbsentReply,
		UserGoogleCalendarEventCreatedReply,
	}
)
