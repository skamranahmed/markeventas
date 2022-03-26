package twitterClient

import (
	"errors"
	"fmt"
	"net/http/httputil"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/skamranahmed/markeventas/pkg/log"
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

func (tb *twitterBot) FetchTweetMentions(sinceTweetID int64) ([]twitter.Tweet, error) {
	// fetch the latest tweet id
	mentionTimelineParams := &twitter.MentionTimelineParams{
		Count:     10,
		TweetMode: "extended",
		SinceID:   sinceTweetID,
	}
	tweets, httpResp, err := tb.client.Timelines.MentionTimeline(mentionTimelineParams)
	if err != nil {
		errMsg := fmt.Sprintf("unable to fetch tweets from the timeline, error: %v, httpResponse: %+v", err, httpResp)
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	return tweets, nil
}

func (tb *twitterBot) ReplyToTweet(tweetID int64, replyBody string) (*twitter.Tweet, string, int, error) {
	tweet, httpResp, err := tb.client.Statuses.Update(replyBody, &twitter.StatusUpdateParams{
		InReplyToStatusID: tweetID,
	})
	statusCode := httpResp.StatusCode

	respBody, _ := httputil.DumpResponse(httpResp, true)
	responseBody := string(respBody)

	if err != nil {
		errMsg := fmt.Sprintf("unable to reply to tweetID: %d from the timeline, error: %v, httpResponse: %+v", tweetID, err, httpResp)
		log.Errorf(errMsg)
		return nil, responseBody, statusCode, errors.New(errMsg)
	}
	return tweet, responseBody, statusCode, nil
}
