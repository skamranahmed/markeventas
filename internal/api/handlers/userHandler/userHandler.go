package handler

import (
	"net/http"

	"github.com/dghubble/oauth1"

	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/service"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
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

	// use the token to get user details from twitter
	twitterUser, err := uh.service.GetUserDetailsFromTwitter(token)
	if err != nil {
		log.Errorf("unable to get user details from twitter, error:%v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userAlreadyExists, user, err := uh.service.DoesUserAlreadyExist(twitterUser.IDStr)
	if userAlreadyExists {
		log.Infof("user with twitterID: %s already exists in our db, userID: %d", twitterUser.IDStr, user.ID)

		// look for twitter screename change and update db
		if twitterUser.ScreenName != user.TwitterScreenName {
			log.Infof(
				"twitterScreenName changed for userID: %d, earlier: %s, now: %s, updating db",
				user.ID, user.TwitterScreenName, twitterUser.ScreenName,
			)
			user.TwitterScreenName = twitterUser.ScreenName
			err = uh.service.Save(user)
			if err != nil {
				log.Warningf(
					"unable to update the twitterScreenName of the userID: %s from %s -> %s, error:%s",
					user.ID, user.TwitterScreenName, twitterUser.ScreenName, err,
				)
			}
		}

		// TODO: generate jwt token and embed the user details in the token payload
		c.JSON(http.StatusOK, gin.H{
			"twitter_id":  twitterUser.ID,
			"screen_name": twitterUser.ScreenName,
		})
		return
	}

	// user doesn't exist in our db, hence create a new user
	log.Infof("👨‍💻 creating new user record for twitterID: %s", twitterUser.IDStr)
	user = &models.User{
		TwitterID:         twitterUser.IDStr,
		TwitterScreenName: twitterUser.ScreenName,
	}

	err = uh.service.Create(user)
	if err != nil {
		log.Errorf("unable to create user record in db, error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Infof("✅ successfuly created new user record for twitterID: %s, userRecord: %+v", twitterUser.IDStr, user)

	c.JSON(http.StatusOK, gin.H{
		"twitter_id":  user.TwitterID,
		"screen_name": user.TwitterScreenName,
	})
	return
}
