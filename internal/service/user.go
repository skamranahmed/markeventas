package service

import (
	"errors"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/skamranahmed/markeventas/config"
	"github.com/skamranahmed/markeventas/internal/models"
	"github.com/skamranahmed/markeventas/internal/repo"
	"github.com/skamranahmed/markeventas/internal/token"
	"github.com/skamranahmed/markeventas/pkg/log"
	"github.com/skamranahmed/markeventas/pkg/twitterClient"
	"gorm.io/gorm"
)

// NewUserService : returns a userService struct that implements the UserService interface
func NewUserService(userRepo repo.UserRepository, googleCalendarTokenRepo repo.GoogleCalendarTokenRepository, tokenMaker token.Maker) UserService {
	twitterOAuthClient := twitterClient.NewTwitterOAuthClient(
		config.TwitterLoginAppApiKey,
		config.TwitterLoginAppApiKeySecret,
		config.TwitterLoginAppCallbackUrl,
	)

	log.Infof("Twitter Callback URL: %s", config.TwitterLoginAppCallbackUrl)

	return &userService{
		userRepo:                userRepo,
		googleCalendarTokenRepo: googleCalendarTokenRepo,
		tokenMaker:              tokenMaker,
		twitterOAuthClient:      twitterOAuthClient,
	}
}

type userService struct {
	userRepo                repo.UserRepository
	googleCalendarTokenRepo repo.GoogleCalendarTokenRepository
	tokenMaker              token.Maker
	twitterOAuthClient      twitterClient.TwitterOAuthClient
}

func (us *userService) NewGoogleService(userID uint, code string) (GoogleService, error) {
	googleService, err := NewGoogleService(userID, code, us.googleCalendarTokenRepo)
	if err != nil {
		log.Errorf("unable to init google service for userID: %d, error: %v", userID, err)
		return nil, err
	}
	return googleService, nil
}

func (us *userService) FindByTwitterID(twitterID string) (*models.User, error) {
	user, err := us.userRepo.FindByTwitterID(twitterID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) Create(u *models.User) error {
	return us.userRepo.Create(u)
}

func (us *userService) SaveGoogleToken(u *models.GoogleCalendarToken) error {
	return us.googleCalendarTokenRepo.Save(u)
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
	requestToken, err := us.twitterOAuthClient.GetRequestToken()
	if err != nil {
		return "", err
	}
	authorizationRedirectURL, err := us.twitterOAuthClient.GetAuthorizationRedirectURL(requestToken)
	if err != nil {
		return "", err
	}
	return authorizationRedirectURL, nil
}

func (us *userService) FetchTwitterOAuthToken(requestToken, requestSecret, verifier string) (*oauth1.Token, error) {
	token, err := us.twitterOAuthClient.GetToken(requestToken, requestSecret, verifier)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (us *userService) GetUserDetailsFromTwitter(token *oauth1.Token) (*twitter.User, error) {
	user, err := us.twitterOAuthClient.GetUser(token)
	if err != nil {
		return nil, err
	}
	return user, nil
}
