package service

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/skamranahmed/markeventas/internal/models"
	"google.golang.org/api/calendar/v3"
)

type UserService interface {
	Create(u *models.User) error
	CreateToken(userID uint, twitterID string) (string, error)
	Save(u *models.User) error
	LoginWithTwitter() (string, error)
	FetchTwitterOAuthToken(requestToken, requestSecret, verifier string) (*oauth1.Token, error)
	GetUserDetailsFromTwitter(token *oauth1.Token) (*twitter.User, error)
	DoesUserAlreadyExist(twitterID string) (bool, *models.User, error)
	NewGoogleService(userID uint, code string) (GoogleService, error)
	SaveGoogleToken(u *models.Token) error
	FindByTwitterID(twitterID string) (*models.User, error)
}

type GoogleService interface {
	CalendarService() *calendar.Service
}

type BotLogService interface {
	Insert(botLog *models.BotLog) (*models.BotLog, error)
	Save(b *models.BotLog) error
	GetLatestRepliedToTweetID() (int64, error)
}
