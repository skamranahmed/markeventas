package handler

import (
	"net/http"

	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/token"
	"gorm.io/gorm"

	"github.com/dghubble/oauth1"

	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/api/middlewares"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/service"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
)

type UserHandler interface {
	TwitterOAuthLogin(c *gin.Context)
	HandleTwitterOAuthCallback(c *gin.Context)
	SaveGoogleCalendarRefreshToken(c *gin.Context)
	GetUserProfile(c *gin.Context)
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
	c.JSON(http.StatusOK, gin.H{"redirect_url": authorizationRedirectURL})
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

		// generate jwt token
		userToken, err := uh.service.CreateToken(user.ID, user.TwitterID)
		if err != nil {
			log.Errorf("unable to create user token, error: %s", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		resp := &TwitterOAuthCallbackResponse{
			AccessToken: userToken,
			TwitterID:   user.TwitterID,
			ScreenName:  user.TwitterScreenName,
		}

		c.JSON(http.StatusOK, resp)
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

	log.Infof("✅ successfully created new user record for twitterID: %s, userRecord: %+v", twitterUser.IDStr, user)

	// generate jwt token
	userToken, err := uh.service.CreateToken(user.ID, user.TwitterID)
	if err != nil {
		log.Errorf("unable to create user token, error: %s", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resp := &TwitterOAuthCallbackResponse{
		AccessToken:              userToken,
		TwitterID:                user.TwitterID,
		ScreenName:               user.TwitterScreenName,
		IsGoogleOauthTokenActive: user.IsGcalOauthTokenActive,
	}

	c.JSON(http.StatusOK, resp)
	return
}

func (uh *userHandler) SaveGoogleCalendarRefreshToken(c *gin.Context) {
	// extract the payload from the context that was set by the AuthMiddleware
	authPayload, exists := c.Get(middlewares.AuthorizationPayloadKey)
	if !exists {
		log.Errorf("authorization payload is not present")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authToken, ok := authPayload.(*token.Payload)
	if !ok {
		log.Errorf("invalid authorization payload, got payload: %+v", authPayload)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var requestPayload GoogleApiCodePayload
	err := c.ShouldBindJSON(&requestPayload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check whehter this user has an account with us
	user, err := uh.service.FindByTwitterID(authToken.TwitterID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warningf("user does not have an account with us")
		} else {
			log.Warningf("unable to fetch user from the db, userTwitterID: %s, err: %v", authToken.TwitterID, err)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = uh.service.NewGoogleService(user.ID, requestPayload.Code)
	if err != nil {
		log.Errorf("unable to init google service for the userID: %d, in SaveGoogleCalendarRefreshToken, err: %v", user.ID, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// after successful creation of google service - mark the oauth token field of users table to true
	if !user.IsGcalOauthTokenActive {
		user.IsGcalOauthTokenActive = true
		err = uh.service.Save(user)
		if err != nil {
			log.Warningf("unable to update the google oauth token status for the userID: %s from false -> true, error:%s", user.ID, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusOK)
	return
}

func (uh *userHandler) GetUserProfile(c *gin.Context) {
	// extract the payload from the context that was set by the AuthMiddleware
	authPayload, exists := c.Get(middlewares.AuthorizationPayloadKey)
	if !exists {
		log.Errorf("authorization payload is not present")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authToken, ok := authPayload.(*token.Payload)
	if !ok {
		log.Errorf("invalid authorization payload, got payload: %+v", authPayload)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// fetch the user profile from the twitterID present in the authToken
	user, err := uh.service.FindByTwitterID(authToken.TwitterID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warningf("user does not have an account with us")
		} else {
			log.Warningf("unable to fetch user from the db, userTwitterID: %s, err: %v", authToken.TwitterID, err)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := &UserProfileResponse{
		ID:                       user.ID,
		TwitterID:                user.TwitterID,
		TwitterScreenName:        user.TwitterScreenName,
		IsGoogleOauthTokenActive: user.IsGcalOauthTokenActive,
	}

	c.JSON(http.StatusOK, gin.H{"user_profile": response})
	return
}
