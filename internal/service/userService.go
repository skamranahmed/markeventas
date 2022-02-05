package service

import (
	"errors"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/repo"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/token"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/twitterClient"
	"gorm.io/gorm"
)

// NewUserService : returns a userService struct that implements the UserService interface
func NewUserService(userRepo repo.UserRepository, tokenRepo repo.TokenRepository, config *config.Config, tokenMaker token.Maker) UserService {
	return &userService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		config:     config,
		tokenMaker: tokenMaker,
	}
}

type userService struct {
	userRepo   repo.UserRepository
	tokenRepo  repo.TokenRepository
	config     *config.Config
	tokenMaker token.Maker
}

func (us *userService) GetUserCalendarService(userID uint, code string) (GoogleService, error) {
	googleService, err := NewGoogleService(userID, code, us.tokenRepo, us.config)
	if err != nil {
		log.Errorf("unable to init google service for userID: %d, error: %v", userID, err)
		return nil, err
	}
	return googleService, nil
}

func (us *userService) Create(u *models.User) error {
	return us.userRepo.Create(u)
}

func (us *userService) SaveGoogleToken(u *models.Token) error {
	return us.tokenRepo.Save(u)
}

func (us *userService) CreateToken(userID uint, twitterID string) (string, error) {
	return us.tokenMaker.CreateToken(userID, twitterID)
}

func (us *userService) Save(u *models.User) error {
	return us.userRepo.Save(u)
}

func (us *userService) DoesUserAlreadyExist(twitterID string) (bool, *models.User, error) {
	user, err := us.userRepo.FindByTwitterID(twitterID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warningf("user with twitterID: %s, does not exist in our db", twitterID)
			return false, nil, nil
		}
		errMsg := fmt.Sprintf("unable to query db for finding user with twitterID: %s, error: %v", twitterID, err)
		return false, nil, errors.New(errMsg)
	}
	return true, user, nil
}

func (us *userService) LoginWithTwitter() (string, error) {
	twitterClient := twitterClient.NewTwitterOAuthClient(
		us.config.TwitterGcalEventLoginAppApiKey,
		us.config.TwitterGcalEventLoginAppApiKeySecret,
		"http://localhost:8080/twitter/callback",
	)
	requestToken, err := twitterClient.GetRequestToken()
	if err != nil {
		return "", err
	}
	authorizationRedirectURL, err := twitterClient.GetAuthorizationRedirectURL(requestToken)
	if err != nil {
		return "", err
	}
	return authorizationRedirectURL, nil
}

func (us *userService) FetchTwitterOAuthToken(requestToken, requestSecret, verifier string) (*oauth1.Token, error) {
	twitterClient := twitterClient.NewTwitterOAuthClient(
		us.config.TwitterGcalEventLoginAppApiKey,
		us.config.TwitterGcalEventLoginAppApiKeySecret,
		"http://localhost:8080/twitter/callback",
	)
	token, err := twitterClient.GetToken(requestToken, requestSecret, verifier)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (us *userService) GetUserDetailsFromTwitter(token *oauth1.Token) (*twitter.User, error) {
	twitterClient := twitterClient.NewTwitterOAuthClient(
		us.config.TwitterGcalEventLoginAppApiKey,
		us.config.TwitterGcalEventLoginAppApiKeySecret,
		"http://localhost:8080/twitter/callback",
	)
	user, err := twitterClient.GetUser(token)
	if err != nil {
		return nil, err
	}
	return user, nil
}
