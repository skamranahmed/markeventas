package service

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
)

type UserService interface {
	Create(u *models.User) error
	LoginWithTwitter() (string, error)
	FetchTwitterOAuthToken(requestToken, requestSecret, verifier string) (*oauth1.Token, error)
	GetUserDetailsFromTwitter(token *oauth1.Token) (*twitter.User, error)
}
