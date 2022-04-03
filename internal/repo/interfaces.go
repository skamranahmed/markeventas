package repo

import (
	"github.com/skamranahmed/markeventas/internal/models"
)

// UserRepository
type UserRepository interface {
	Create(u *models.User) error
	Save(u *models.User) error
	FindByTwitterID(twitterID string) (*models.User, error)
}

// GoogleCalendarTokenRepository
type GoogleCalendarTokenRepository interface {
	Create(u *models.GoogleCalendarToken) error
	GetUserToken(userID uint) (*models.GoogleCalendarToken, error)
	Save(u *models.GoogleCalendarToken) error
}

// BotLogRepository
type BotLogRepository interface {
	Create(botLog *models.BotLog) (*models.BotLog, error)
	Save(b *models.BotLog) error
	GetLatestRepliedToTweetID() (int64, error)
}
