package twitterClient

import (
	"errors"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	twitterAuth "github.com/dghubble/oauth1/twitter"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
)

var (
	twitterOAuthConfig *oauth1.Config
)

func NewTwitterOAuthClient(consumerKey, consumerSecret, callbackURL string) TwitterOAuthClient {
	return &twitterOAuthClient{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		CallbackURL:    callbackURL,
	}
}

type twitterOAuthClient struct {
	ConsumerKey    string
	ConsumerSecret string
	CallbackURL    string
}

func (t *twitterOAuthClient) GetRequestToken() (string, error) {
	twitterOAuthConfig = &oauth1.Config{
		ConsumerKey:    t.ConsumerKey,
		ConsumerSecret: t.ConsumerSecret,
		CallbackURL:    t.CallbackURL,
		Endpoint:       twitterAuth.AuthenticateEndpoint,
	}

	requestToken, _, err := twitterOAuthConfig.RequestToken()
	if err != nil {
		log.Errorf("unable to get request token from twitter, error: %v", err)
		return "", err
	}

	return requestToken, nil
}

func (t *twitterOAuthClient) GetAuthorizationRedirectURL(requestToken string) (string, error) {
	authorizationRedirectURL, err := twitterOAuthConfig.AuthorizationURL(requestToken)
	if err != nil {
		log.Errorf("unable to get authorization redirect URL from twitter, error: %v", err)
		return "", err
	}

	if authorizationRedirectURL == nil {
		errMsg := fmt.Sprintf("got empty autorization redirect URL")
		log.Error(errMsg)
		return "", errors.New(errMsg)
	}

	return authorizationRedirectURL.String(), nil
}

func (t *twitterOAuthClient) GetToken(requestToken, requestSecret, verifier string) (*oauth1.Token, error) {
	accessToken, accessSecret, err := twitterOAuthConfig.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get accessToken and accessSecret, error: %v", err)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	token := oauth1.NewToken(accessToken, accessSecret)
	if token == nil {
		errMsg := fmt.Sprintf("twitter oAuth token is nil")
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return token, nil
}

func (t *twitterOAuthClient) GetUser(token *oauth1.Token) (*twitter.User, error) {
	httpClient := twitterOAuthConfig.Client(oauth1.NoContext, token)
	twitterClient := twitter.NewClient(httpClient)
	accountVerifyParams := &twitter.AccountVerifyParams{
		IncludeEntities: twitter.Bool(false),
		SkipStatus:      twitter.Bool(true),
		IncludeEmail:    twitter.Bool(false),
	}
	
	user, resp, err := twitterClient.Accounts.VerifyCredentials(accountVerifyParams)
	if err != nil {
		log.Errorf("unable to verify credentials of the user, error: %v, responseCode: %d", err, resp.StatusCode)
		return nil, err
	}

	if user == nil {
		errMsg := fmt.Sprintf("got empty twitter user")
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	return user, nil
}
