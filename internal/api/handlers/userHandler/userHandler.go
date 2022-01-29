package handler

import (
	"net/http"

	"github.com/dghubble/oauth1"

	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/service"
)

type UserHandler interface {
	TwitterOAuthLogin(c *gin.Context)
	HandleTwitterOAuthCallback(c *gin.Context)
}

func NewUserHandler(userService service.UserService, config *config.Config) UserHandler {
	return &userHandler{
		service: userService,
		config:  config,
	}
}

type userHandler struct {
	service service.UserService
	config  *config.Config
}

var (
	twitterOAuthConfig *oauth1.Config
)

func (uh *userHandler) TwitterOAuthLogin(c *gin.Context) {
	authorizationRedirectURL, err := uh.service.LoginWithTwitter()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.Redirect(http.StatusFound, authorizationRedirectURL)
	return
}

func (uh *userHandler) HandleTwitterOAuthCallback(c *gin.Context) {
	var requestPayload TwitterOAuthCallbackPayload
	err := c.BindJSON(&requestPayload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	requestToken := requestPayload.OAuthToken
	verifier := requestPayload.OAuthVerifier
	// for some reason, we don't require the requestSecret for getting the accessToken, hence keeping it empty as of now
	requestSecret := ""

	token, err := uh.service.FetchTwitterOAuthToken(requestToken, requestSecret, verifier)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	user, err := uh.service.GetUserDetailsFromTwitter(token)
	
	// TODO: check the twitter_id in the database and create an account for the user if it doesn't exist

	c.JSON(http.StatusOK, gin.H{
		"twitter_id":  user.ID,
		"screen_name": user.ScreenName,
	})
	return
}
