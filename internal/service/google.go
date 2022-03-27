package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/skamranahmed/markeventas/config"
	"github.com/skamranahmed/markeventas/internal/models"
	"github.com/skamranahmed/markeventas/internal/repo"
	"github.com/skamranahmed/markeventas/pkg/log"
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

func NewGoogleService(userID uint, code string, tokenRepo repo.TokenRepository) (GoogleService, error) {
	ctx := context.Background()
	clientSecret := []byte(config.GoogleAppClientSecret)
	oAuth2Config, err := google.ConfigFromJSON(clientSecret, calendar.CalendarEventsScope)
	if err != nil {
		log.Errorf("unable to parse client secret file to config: %v", err)
		return nil, errors.New(err.Error())
	}

	var service googleService
	service.tokenRepo = tokenRepo
	service.oAuth2Config = oAuth2Config

	var googleOAuth2Token *oauth2.Token

	if code == "" {
		// fetch the token from db
		googleOAuth2Token, _, err = service.GetUserTokenFromDB(userID)
		if err != nil {
			return nil, err
		}
	} else {
		// get token from web and then create a new token record in the db for future use
		googleOAuth2Token, err = service.GetUserTokenFromWeb(userID, code)
		if err != nil {
			return nil, err
		}

		// TODO: handle token save error
		_ = service.SaveUserToken(googleOAuth2Token, userID)
	}

	// init the http client for googleService
	service.client = oAuth2Config.Client(context.Background(), googleOAuth2Token)

	// init the calendar service
	calendarService, err := calendar.NewService(ctx, option.WithHTTPClient(service.client))
	if err != nil {
		log.Errorf("unable to retrieve Calendar client: %v", err)
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

func (g *googleService) GetUserTokenFromDB(userID uint) (*oauth2.Token, *models.Token, error) {
	token, err := g.tokenRepo.GetUserToken(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Infof("google token for userID: %d, not found in DB", userID)
			return nil, nil, err
		}
		log.Errorf("unable to get google token for userID: %d from the DB, error:%v", userID, err)
		return nil, nil, err
	}

	googleOAuth2Token := &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: token.RefreshToken,
		Expiry:       token.ExpiresAt,
	}

	return googleOAuth2Token, token, nil
}

func (g *googleService) SaveUserToken(googleOAuth2Token *oauth2.Token, userID uint) error {
	// fetch the token from db
	_, existingTokenRecord, err := g.GetUserTokenFromDB(userID)
	if err != nil {
		// if no existing token found, then create a new token record
		if err == gorm.ErrRecordNotFound {
			newTokenRecord := &models.Token{
				UserID:       userID,
				AccessToken:  googleOAuth2Token.AccessToken,
				RefreshToken: googleOAuth2Token.RefreshToken,
				ExpiresAt:    googleOAuth2Token.Expiry,
			}
			log.Infof("google token not found in DB for userID: %d | creating a new token record: %+v", userID, newTokenRecord)
			err := g.tokenRepo.Create(newTokenRecord)
			if err != nil {
				log.Errorf("unable to create google token record in db: %+v, for userID: %d, error: %v", newTokenRecord, userID, err)
				return err
			}
			return nil
		}
		log.Errorf("unable to query db for token record of userID: %d, error: %v", userID, err)
		return err
	}

	log.Infof("google token found in DB for userID: %d, tokenRecord: %+v", userID, existingTokenRecord)

	// if a token already exists in db then update the record
	existingTokenRecord.UpdatedAt = time.Now()
	existingTokenRecord.AccessToken = googleOAuth2Token.AccessToken
	existingTokenRecord.RefreshToken = googleOAuth2Token.RefreshToken
	existingTokenRecord.ExpiresAt = googleOAuth2Token.Expiry

	err = g.tokenRepo.Save(existingTokenRecord)
	if err != nil {
		log.Errorf("unable to update google token record in db: %+v, for userID: %d, error: %v", existingTokenRecord, userID, err)
		return err
	}

	log.Infof("✅ successfully updated the existing google token record: %+v, for userID: %d", existingTokenRecord, userID)
	return nil
}
