package handler

import (
	"net/http"

	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/token"

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
		AccessToken: userToken,
		TwitterID:   user.TwitterID,
		ScreenName:  user.TwitterScreenName,
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

	userID := authToken.UserID

	_, err = uh.service.NewGoogleService(userID, requestPayload.Code)
	// TODO: error Handling
	// googleService, err := uh.service.NewGoogleService(userID, "")
	// calendarService := googleService.CalendarService()

	// t := time.Now().Format(time.RFC3339)
	// events, err := calendarService.Events.List("primary").ShowDeleted(false).
	// 	SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	// if err != nil {
	// 	log.Errorf("Unable to retrieve next ten of the user's events: %v", err)
	// }
	// fmt.Println("Upcoming events:")
	// if len(events.Items) == 0 {
	// 	fmt.Println("No upcoming events found.")
	// } else {
	// 	for _, item := range events.Items {
	// 		date := item.Start.DateTime
	// 		if date == "" {
	// 			date = item.Start.Date
	// 		}
	// 		fmt.Printf("%v (%v)\n", item.Summary, date)
	// 	}
	// }

	c.Status(http.StatusOK)
	return
}
