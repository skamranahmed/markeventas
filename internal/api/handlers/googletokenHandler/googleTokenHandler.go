package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/api/middlewares"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/service"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/token"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
	"gorm.io/gorm"
)

type GoogleTokenHandler interface {
	SaveRefreshToken(c *gin.Context)
}

func NewGoogleTokenHandler(userService service.UserService) GoogleTokenHandler {
	return &googleTokenHandler{
		userService: userService,
	}
}

type googleTokenHandler struct {
	userService service.UserService
}

func (gh *googleTokenHandler) SaveRefreshToken(c *gin.Context) {
	// extract the payload from the context that was set by the AuthMiddleware
	jwtAuthPayload, exists := c.Get(middlewares.AuthorizationPayloadKey)
	if !exists {
		log.Errorf("authorization payload is not present")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	jwtAuthToken, ok := jwtAuthPayload.(*token.Payload)
	if !ok {
		log.Errorf("invalid authorization payload, got payload: %+v", jwtAuthToken)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var requestPayload RefreshTokenCodeRequest
	err := c.ShouldBindJSON(&requestPayload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check whehter this user has an account with us
	user, err := gh.userService.FindByTwitterID(jwtAuthToken.TwitterID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warningf("user does not have an account with us")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		log.Warningf("unable to fetch user from the db, userTwitterID: %s, err: %v", jwtAuthToken.TwitterID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = gh.userService.NewGoogleService(user.ID, requestPayload.Code)
	if err != nil {
		log.Errorf("unable to init google service for the userID: %d, in SaveGoogleCalendarRefreshToken, err: %v", user.ID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// after successful creation of google service - mark the oauth token field of users table to true
	if !user.IsGcalOauthTokenActive {
		user.IsGcalOauthTokenActive = true
		err = gh.userService.Save(user)
		if err != nil {
			log.Warningf("unable to update the google oauth token status for the userID: %s from false -> true, error:%s", user.ID, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusOK)
	return
}
