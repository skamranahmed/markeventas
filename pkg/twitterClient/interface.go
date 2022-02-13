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

type TwitterBotClient interface {
	FetchTweetMentions(sinceTweetID int64) ([]twitter.Tweet, error)
	ReplyToTweet(tweetID int64, replyBody string) (*twitter.Tweet, string, int, error)
}
