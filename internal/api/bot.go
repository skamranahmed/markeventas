package api

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/skamranahmed/markeventas/config"
	"github.com/skamranahmed/markeventas/internal/models"
	"github.com/skamranahmed/markeventas/internal/service"
	"github.com/skamranahmed/markeventas/internal/utils"
	"github.com/skamranahmed/markeventas/pkg/log"
	"github.com/skamranahmed/markeventas/pkg/twitterClient"
	"google.golang.org/api/calendar/v3"
	"gorm.io/gorm"
)

func startTwitterBot(userService service.UserService, botLogService service.BotLogService) {
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
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for t := range ticker.C {
			sinceTweetID, err := botLogService.GetLatestRepliedToTweetID()
			if err != nil {
				log.Errorf("unable to fetch latestRepliedToTweetID from the db, err: %v", err)
				sinceTweetID = 0
			}

			tweets, err := twitterBot.FetchTweetMentions(sinceTweetID)
			if err != nil {
				log.Error(err)
				continue
			}

			log.Infof("Got %d tweets to process at %v", len(tweets), t)

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

				// extract the user tweet text
				userTweetText := tweet.FullText

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

				botLogRecord := &models.BotLog{
					ToTweetID:             userTweetID,
					UserTwitterID:         userTwitterID,
					UserTwitterScreenName: userTwitterScreenName,
					UserTweet:             userTweetText,
				}

				botLogRecord, err = botLogService.Insert(botLogRecord)
				if err != nil {
					log.Errorf("unable to insert bot log record: %+v in db, err: %v", botLogRecord, err)
					continue
				}

				if !userAccountExists {
					// user account does not exist
					// send user doesn't have an account with us reply to the user
					log.Warningf("user account for twitterID: %s does not exist, sending UserAccountDoesNotExistReply", userTwitterID)
					body := fmt.Sprintf(utils.UserAccountDoesNotExistReply.Body, userTwitterScreenName)
					_, responseBody, statusCode, err := twitterBot.ReplyToTweet(userTweetID, body)
					if err != nil {
						log.Errorf("unable to send UserAccountDoesNotExistReply to userTwitterID: %s, err: %v", userTwitterID, err)
					}

					// update the botLog record
					botLogRecord.Reply = body
					botLogRecord.HttpResponse = responseBody
					botLogRecord.HttpStatusCode = statusCode
					botLogRecord.ReplyTypeCode = utils.UserAccountDoesNotExistReply.Code

					err = botLogService.Save(botLogRecord)
					if err != nil {
						log.Errorf("unable to save the updated botLogRecord: %+v, err: %v", botLogRecord, err)
					}

					continue
				}

				// checking the status of user google calendar oauth token
				if !user.IsGcalOauthTokenActive {
					// send a reply that the user has not provided the google calendar consent
					// send google calendar consent absent reply to the user
					log.Warningf("user google calendar consent absent for userID: %d, sending UserGoogleCalendarConsentAbsentReply", user.ID)
					body := fmt.Sprintf(utils.UserGoogleCalendarConsentAbsentReply.Body, userTwitterScreenName)
					_, responseBody, statusCode, err := twitterBot.ReplyToTweet(userTweetID, body)
					if err != nil {
						log.Errorf("unable to send UserGoogleCalendarConsentAbsentReply to userID: %d, err: %v", user.ID, err)
					}

					// update the botLog record
					botLogRecord.Reply = body
					botLogRecord.HttpResponse = responseBody
					botLogRecord.HttpStatusCode = statusCode
					botLogRecord.ReplyTypeCode = utils.UserGoogleCalendarConsentAbsentReply.Code

					err = botLogService.Save(botLogRecord)
					if err != nil {
						log.Errorf("unable to save the updated botLogRecord: %+v, err: %v", botLogRecord, err)
					}

					continue
				}

				// checking the google calendar consent
				// the code argument would be empty as be would use the refresh token of the user saved in the db
				googleService, err := userService.NewGoogleService(user.ID, "")
				if err != nil {
					log.Warningf("unable to init google service for the userID: %d, err: %v", user.ID, err)

					// change the oauth token active status to false
					user.IsGcalOauthTokenActive = false
					err := userService.Save(user)
					if err != nil {
						log.Errorf("unable to update the user gcal oauth token status to false, user: %+v, err: %v", user, err)
					}

					// send google calendar consent absent reply to the user
					log.Warningf("user google calendar consent absent for userID: %d, sending UserGoogleCalendarConsentAbsentReply", user.ID)
					body := fmt.Sprintf(utils.UserGoogleCalendarConsentAbsentReply.Body, userTwitterScreenName)
					_, responseBody, statusCode, err := twitterBot.ReplyToTweet(userTweetID, body)
					if err != nil {
						log.Errorf("unable to send UserGoogleCalendarConsentAbsentReply to userID: %d, err: %v", user.ID, err)
					}

					// update the botLog record
					botLogRecord.Reply = body
					botLogRecord.HttpResponse = responseBody
					botLogRecord.HttpStatusCode = statusCode
					botLogRecord.ReplyTypeCode = utils.UserGoogleCalendarConsentAbsentReply.Code

					err = botLogService.Save(botLogRecord)
					if err != nil {
						log.Errorf("unable to save the updated botLogRecord: %+v, err: %v", botLogRecord, err)
					}

					continue
				}

				// parse the user tweet text
				// _ = "Kamran's Space | Jan 28, 2022 | 6:43 PM | IST"
				userTweetData, err := utils.ParseTweetText(userTweetText)
				if err != nil {
					// send reply to the user saying that the format of the tweet is incorrect
					log.Error(err)
					continue
				}

				tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", userTwitterScreenName, statusID)

				calendarService := googleService.CalendarService()
				event := &calendar.Event{
					Summary:     userTweetData.SpaceName,
					Description: tweetURL,
					Start: &calendar.EventDateTime{
						DateTime: userTweetData.StartDateTimeString,
						TimeZone: userTweetData.TimeZoneIanaName,
					},
					ColorId: "2",
					End: &calendar.EventDateTime{
						DateTime: userTweetData.EndDateTimeString,
						TimeZone: userTweetData.TimeZoneIanaName,
					},
				}

				calendarId := "primary"
				event, err = calendarService.Events.Insert(calendarId, event).Do()
				if err != nil {
					log.Errorf("Unable to create google calendar event for the userID: %d, err: %v", user.ID, err)

					// change the oauth token active status to false
					user.IsGcalOauthTokenActive = false
					err := userService.Save(user)
					if err != nil {
						log.Errorf("unable to update the user gcal oauth token status to false, user: %+v, err: %v", user, err)
					}

					// send google calendar consent absent reply to the user
					// send error reply to the user
					log.Warningf("user google calendar consent absent for userID: %d, sending UserGoogleCalendarConsentAbsentReply", user.ID)
					body := fmt.Sprintf(utils.UserGoogleCalendarConsentAbsentReply.Body, userTwitterScreenName)
					_, responseBody, statusCode, err := twitterBot.ReplyToTweet(userTweetID, body)
					if err != nil {
						log.Errorf("unable to send UserGoogleCalendarConsentAbsentReply to userID: %d, err: %v", user.ID, err)
					}

					// update the botLog record
					botLogRecord.Reply = body
					botLogRecord.HttpResponse = responseBody
					botLogRecord.HttpStatusCode = statusCode
					botLogRecord.ReplyTypeCode = utils.UserGoogleCalendarConsentAbsentReply.Code

					err = botLogService.Save(botLogRecord)
					if err != nil {
						log.Errorf("unable to save the updated botLogRecord: %+v, err: %v", botLogRecord, err)
					}

					continue
				}

				log.Infof("google calendar event created for userID: %d, sending UserGoogleCalendarEventCreatedReply", user.ID)
				body := fmt.Sprintf(utils.UserGoogleCalendarEventCreatedReply.Body, userTwitterScreenName, event.HtmlLink)
				_, responseBody, statusCode, err := twitterBot.ReplyToTweet(userTweetID, body)

				// convert the user tweet data into JSON for saving in db
				userTweetJson, _ := json.Marshal(userTweetData)

				// update the botLog record
				botLogRecord.Reply = body
				botLogRecord.UserParsedTweetData = string(userTweetJson)
				botLogRecord.HttpResponse = responseBody
				botLogRecord.HttpStatusCode = statusCode
				botLogRecord.ReplyTypeCode = utils.UserGoogleCalendarEventCreatedReply.Code

				err = botLogService.Save(botLogRecord)
				if err != nil {
					log.Errorf("unable to save the updated botLogRecord: %+v, err: %v", botLogRecord, err)
				}
			}
		}
	}()
}
