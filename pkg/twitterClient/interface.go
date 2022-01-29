package twitterClient

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type TwitterOAuthClient interface {
	GetRequestToken() (string, error)
	GetAuthorizationRedirectURL(requestToken string) (string, error)
	GetToken(requestToken, requestSecret, verifier string) (*oauth1.Token, error)
	GetUser(token *oauth1.Token) (*twitter.User, error)
}