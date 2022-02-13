package repo

import (
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
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

type BotReplyRepository interface {
	Create(botReply *models.BotReply) (*models.BotReply, error)
	Save(b *models.BotReply) error
	GetLatestRepliedToTweetID() (int64, error)
}
