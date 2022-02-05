package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/repo"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type googleService struct {
	tokenRepo    repo.TokenRepository
	oAuth2Config *oauth2.Config
	client       *http.Client
	calendar     *calendar.Service
}

func NewGoogleService(userID uint, code string, tokenRepo repo.TokenRepository, config *config.Config) (GoogleService, error) {
	ctx := context.Background()
	clientSecret := []byte(config.GoogleAppClientSecret)
	oAuth2Config, err := google.ConfigFromJSON(clientSecret, calendar.CalendarScope)
	if err != nil {
		log.Errorf("Unable to parse client secret file to config: %v", err)
		return nil, errors.New(err.Error())
	}

	var service googleService
	service.tokenRepo = tokenRepo
	service.oAuth2Config = oAuth2Config

	var token *oauth2.Token

	if code == "" {
		// fetch the token from db
		token, err = service.GetUserTokenFromDB(userID)
		if err != nil {
			return nil, err
		}
	} else {
		// get token from web and then create a new token record in the db for future use
		token, err = service.GetUserTokenFromWeb(userID, code)
		if err != nil {
			return nil, err
		}

		// TODO: handle token save error
		_ = service.SaveUserToken(token, userID)
	}

	// init the http client for googleService
	service.client = oAuth2Config.Client(context.Background(), token)

	// init the calendar service
	calendarService, err := calendar.NewService(ctx, option.WithHTTPClient(service.client))
	if err != nil {
		log.Errorf("Unable to retrieve Calendar client: %v", err)
		return nil, err
	}

	service.calendar = calendarService
	return service, nil
}

func (g googleService) CalendarService() *calendar.Service {
	return g.calendar
}

func (g *googleService) GetUserTokenFromWeb(userID uint, code string) (*oauth2.Token, error) {
	token, err := g.oAuth2Config.Exchange(context.TODO(), code)
	if err != nil {
		log.Errorf("unable to get google user tokens for the userID: %d, code: %s, error: %v", userID, code, err)
		return nil, err
	}
	return token, nil
}

func (g *googleService) GetUserTokenFromDB(userID uint) (*oauth2.Token, error) {
	token, err := g.tokenRepo.GetUserToken(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Infof("google token for userID: %d, not found in DB", userID)
			return nil, err
		}
		log.Errorf("unable to get google token for userID: %d from the DB, error:%v", userID, err)
		return nil, err
	}

	googleOAuth2Token := &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: token.RefreshToken,
		Expiry:       token.ExpiresAt,
	}

	return googleOAuth2Token, nil
}

func (g *googleService) SaveUserToken(googleToken *oauth2.Token, userID uint) error {
	token := &models.Token{
		UserID:       userID,
		AccessToken:  googleToken.AccessToken,
		RefreshToken: googleToken.RefreshToken,
		ExpiresAt:    googleToken.Expiry,
	}
	err := g.tokenRepo.Create(token)
	if err != nil {
		log.Errorf("unable to save google token: %+v, for userID: %d, error: %v", token, userID, err)
		return err
	}
	return nil
}