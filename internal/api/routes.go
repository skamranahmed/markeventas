package api

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
	handler "github.com/skamranahmed/twitter-create-gcal-event-api/internal/api/handlers/userHandler"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/api/middlewares"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/repo"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/service"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/token"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/utils"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/twitterClient"
	"google.golang.org/api/calendar/v3"
	"gorm.io/gorm"
)

var (
	tokenMaker token.Maker
)

type repos struct {
	userRepo     repo.UserRepository
	tokenRepo    repo.TokenRepository
	botReplyRepo repo.BotReplyRepository
}

type services struct {
	userService     service.UserService
	botReplyService service.BotReplyService
}

type handlers struct {
	userHandler handler.UserHandler
}

func startTwitterBot(config *config.Config, userService service.UserService, botReplyService service.BotReplyService) {
	twitterBot, err := twitterClient.NewTwitterBotClient(
		config.TwitterBotAccessToken,
		config.TwitterBotAccessTokenSecret,
		config.TwitterBotApiKey,
		config.TwitterBotApiKeySecret,
	)
	if err != nil {
		log.Errorf("unable to init twitter bot client, error: %v", err)
		os.Exit(1)
	}

	// create a ticker that will fetch the twitter timeline mentions and then send a reply
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for range ticker.C {
			sinceTweetID, err := botReplyService.GetLatestRepliedToTweetID()
			if err != nil {
				log.Errorf("unable to fetch latestRepliedToTweetID from the db, err: %v", err)
				sinceTweetID = 0
			}

			tweets, err := twitterBot.FetchTweetMentions(sinceTweetID)
			if err != nil {
				log.Error(err)
				continue
			}

			for _, tweet := range tweets {
				// assume that the user has an account with us
				userAccountExists := true

				// grab the twitterID of the user who mentioned us in the tweet
				userTwitterID := tweet.InReplyToUserIDStr

				// grab the statusID
				statusID := tweet.InReplyToStatusIDStr

				// latestTweetID - required for SinceID
				userTweetID := tweet.ID

				// userTwitterScreenName
				userTwitterScreenName := tweet.InReplyToScreenName

				// check whehter this user has an account with us
				user, err := userService.FindByTwitterID(userTwitterID)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						log.Warningf("user does not have an account with us")
						userAccountExists = false
					} else {
						log.Warningf("unable to fetch user from the db, userTwitterID: %s, err: %v", userTwitterID, err)
						continue
					}
				}

				botReplyRecord := &models.BotReply{
					ToTweetID:             userTweetID,
					UserTwitterID:         userTwitterID,
					UserTwitterScreenName: userTwitterScreenName,
				}

				botReplyRecord, err = botReplyService.Insert(botReplyRecord)
				if err != nil {
					log.Errorf("unable to insert bot reply record: %+v in db, err: %v", botReplyRecord, err)
					continue
				}

				if !userAccountExists {
					// user account does not exist
					// send user doesn't have an account with us reply to the user
					body := fmt.Sprintf(utils.UserAccountDoesNotExistReply.Body, userTwitterScreenName)
					_, responseBody, statusCode, err := twitterBot.ReplyToTweet(userTweetID, body)

					// update the botReply record
					botReplyRecord.Body = body
					botReplyRecord.Response = responseBody
					botReplyRecord.StatusCode = statusCode
					botReplyRecord.TypeCode = utils.UserAccountDoesNotExistReply.Code

					err = botReplyService.Save(botReplyRecord)
					if err != nil {
						log.Errorf("unable to save the updated botReplyRecord: %+v, err: %v", botReplyRecord, err)
					}

					continue
				}

				
				// checking the google calendar consent
				// the code argument would be empty as be would use the refresh token of the user saved in the db
				googleService, err := userService.NewGoogleService(user.ID, "")
				if err != nil {
					log.Warningf("unable to init google service for the userID: %d, err: %v", user.ID, err)

					// send google calendar consent absent reply to the user
					body := fmt.Sprintf(utils.UserGoogleCalendarConsentAbsentReply.Body, userTwitterScreenName)
					_, responseBody, statusCode, err := twitterBot.ReplyToTweet(userTweetID, body)

					// update the botReply record
					botReplyRecord.Body = body
					botReplyRecord.Response = responseBody
					botReplyRecord.StatusCode = statusCode
					botReplyRecord.TypeCode = utils.UserGoogleCalendarConsentAbsentReply.Code

					err = botReplyService.Save(botReplyRecord)
					if err != nil {
						log.Errorf("unable to save the updated botReplyRecord: %+v, err: %v", botReplyRecord, err)
					}

					continue
				}

				tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", userTwitterScreenName, statusID)

				calendarService := googleService.CalendarService()
				event := &calendar.Event{
					Summary:     "Twitter Space",
					Description: tweetURL,
					Start: &calendar.EventDateTime{
						DateTime: "2022-01-26T09:00:00-07:00",
						TimeZone: "Asia/Calcutta",
					},
					ColorId: "2",
					End: &calendar.EventDateTime{
						DateTime: "2022-01-26T10:00:00-07:00",
						TimeZone: "Asia/Calcutta",
					},
				}

				calendarId := "primary"
				event, err = calendarService.Events.Insert(calendarId, event).Do()
				if err != nil {
					log.Errorf("Unable to create event. %v\n", err)
					continue
				}

				body := fmt.Sprintf(utils.UserGoogleCalendarEventCreatedReply.Body, userTwitterScreenName, event.HtmlLink)
				_, responseBody, statusCode, err := twitterBot.ReplyToTweet(userTweetID, body)

				// update the botReply record
				botReplyRecord.Body = body
				botReplyRecord.Response = responseBody
				botReplyRecord.StatusCode = statusCode
				botReplyRecord.TypeCode = utils.UserGoogleCalendarEventCreatedReply.Code

				err = botReplyService.Save(botReplyRecord)
				if err != nil {
					log.Errorf("unable to save the updated botReplyRecord: %+v, err: %v", botReplyRecord, err)
				}
			}
		}
	}()
}

func InitRoutes(db *gorm.DB, config *config.Config) *gin.Engine {
	tokenMaker = token.NewJwtTokenMaker(config.TokenSecretSigningKey)
	_, services, handlers := setDependencies(db, config)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Hello": "World",
		})
		return
	})
	router.GET("/api/login", handlers.userHandler.TwitterOAuthLogin)
	router.POST("/api/twitter/callback", handlers.userHandler.HandleTwitterOAuthCallback)

	authorized := router.Group("/api/google/calendar")
	authorized.Use(middlewares.AuthMiddleware(tokenMaker))
	{
		authorized.POST("/auth-code", handlers.userHandler.SaveGoogleCalendarRefreshToken)
	}

	// run the twitter bot in background
	go startTwitterBot(config, services.userService, services.botReplyService)

	return router
}

func setDependencies(db *gorm.DB, config *config.Config) (*repos, *services, *handlers) {
	repos := &repos{}
	repos.setDependencies(db)

	services := &services{}
	services.setDependencies(repos, config)

	handlers := &handlers{}
	handlers.setDependencies(services, config)

	return repos, services, handlers
}

func (r *repos) setDependencies(db *gorm.DB) {
	userRepo := repo.NewUserRepository(db)
	r.userRepo = userRepo

	tokenRepo := repo.NewTokenRepository(db)
	r.tokenRepo = tokenRepo

	botReplyRepo := repo.NewBotReplyRepository(db)
	r.botReplyRepo = botReplyRepo
}

func (s *services) setDependencies(repos *repos, config *config.Config) {
	userService := service.NewUserService(repos.userRepo, repos.tokenRepo, config, tokenMaker)
	s.userService = userService

	botReplyService := service.NewBotReplyService(repos.botReplyRepo)
	s.botReplyService = botReplyService
}

func (h *handlers) setDependencies(services *services, config *config.Config) {
	userHandler := handler.NewUserHandler(services.userService, config)
	h.userHandler = userHandler
}
