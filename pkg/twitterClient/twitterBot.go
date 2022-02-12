package twitterClient

import (
	"errors"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
)

func NewTwitterBotClient(accessToken, accessTokenSecret, consumerKey, consumerSecret string) (TwitterBotClient, error) {
	log.Info("Create Gcal Event Twitter Bot")

	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken(accessToken, accessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	// we can retrieve the user and verify if the credentials
	// we have used successfully allow us to log in!
	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	return &twitterBot{
		client: client,
	}, nil
}

type twitterBot struct {
	client *twitter.Client
}

func (tb *twitterBot) FetchTweetMentions() ([]twitter.Tweet, error) {
	// fetch the latest tweet id
	mentionTimelineParams := &twitter.MentionTimelineParams{
		Count:     10,
		TweetMode: "extended",
		// SinceID:   latestTweetID,
	}
	tweets, httpResponse, err := tb.client.Timelines.MentionTimeline(mentionTimelineParams)
	if err != nil {
		errMsg := fmt.Sprintf("unable to fetch tweets from the timeline, error: %v, httpResponse: %+v", err, httpResponse)
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	return tweets, nil
}
