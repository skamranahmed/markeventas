package handler

import (
	"fmt"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	twitterAuth "github.com/dghubble/oauth1/twitter"

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
	twitterOAuthConfig = &oauth1.Config{
		ConsumerKey:    uh.config.TwitterGcalEventLoginAppApiKey,
		ConsumerSecret: uh.config.TwitterGcalEventLoginAppApiKeySecret,
		CallbackURL:    "http://localhost:8080/twitter/callback",
		Endpoint:       twitterAuth.AuthenticateEndpoint,
	}

	requestToken, requestSecret, err := twitterOAuthConfig.RequestToken()
	fmt.Println("Request Token: ", requestToken)
	fmt.Println("Request Secret: ", requestSecret)
	fmt.Println("Error: ", err)

	authorizationURL, err := twitterOAuthConfig.AuthorizationURL(requestToken)
	fmt.Println("Error: ", err)

	c.Redirect(http.StatusFound, authorizationURL.String())
	return
}

func (uh *userHandler) HandleTwitterOAuthCallback(c *gin.Context) {
	var requestPayload TwitterOAuthCallbackPayload
	err := c.BindJSON(&requestPayload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// requestToken, verifier, err := oauth1.ParseAuthorizationCallback(c.Request)
	// fmt.Println("Error: ", err)

	// We will get the requestToken (i.e auth_token) and verifier from the frontend, so we will skip the above two lines

	requestToken := requestPayload.OAuthToken
	verifier := requestPayload.OAuthVerifier
	requestSecret := ""

	accessToken, accessSecret, err := twitterOAuthConfig.AccessToken(requestToken, requestSecret, verifier)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := twitterOAuthConfig.Client(oauth1.NoContext, token)
	twitterClient := twitter.NewClient(httpClient)
	accountVerifyParams := &twitter.AccountVerifyParams{
		IncludeEntities: twitter.Bool(false),
		SkipStatus:      twitter.Bool(true),
		IncludeEmail:    twitter.Bool(false),
	}
	user, resp, err := twitterClient.Accounts.VerifyCredentials(accountVerifyParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	fmt.Printf("Resp: %+v\n", resp)
	// // fmt.Fprintf(w, `Hello %s, your Twitter ID is: %d`, user.ScreenName, user.ID)

	// TODO: check the twitter_id in the database and create an account for the user if it doesn't exist

	c.JSON(http.StatusOK, gin.H{
		"twitter_id":  user.ID,
		"screen_name": user.ScreenName,
	})

	return

}
