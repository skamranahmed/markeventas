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

type TokenRepository interface {
	Create(u *models.Token) error
	GetUserToken(userID uint) (*models.Token, error)
	Save(u *models.Token) error
}

type BotLogRepository interface {
	Create(botLog *models.BotLog) (*models.BotLog, error)
	Save(b *models.BotLog) error
	GetLatestRepliedToTweetID() (int64, error)
}
