package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/skamranahmed/markeventas/pkg/log"
	"github.com/tkuchiki/go-timezone"
)

type ParsedUserTweetContent struct {
	SpaceName           string
	StartDateTimeString string
	EndDateTimeString   string
	TimeZoneIanaName    string
}

func ParseTweetText(tweetText string) (*ParsedUserTweetContent, error) {
	// @creategcalevent <space_name> | <date> | <time> | <time_zone>
	tweetText = processTweetText(tweetText)

	strSlice := strings.SplitAfter(tweetText, "|")
	if len(strSlice) != 4 {
		return nil, errors.New("invalid tweet text format")
	}

	spaceName, err := processSpaceName(strSlice[0])
	if err != nil {
		log.Error(err)
		return nil, errors.New("invalid space name")
	}

	dateValue, err := processDateValue(strSlice[1])
	if err != nil {
		log.Error(err)
		return nil, errors.New("invalid date format")
	}

	timeValue, err := processTimeString(strSlice[2])
	if err != nil {
		log.Error(err)
		return nil, errors.New("invalid time format")
	}

	timeZoneAbbr, err := processTimeZoneValue(strSlice[3])
	if err != nil {
		log.Error(err)
		return nil, errors.New("invalid time zone format")
	}

	log.Infof("Space Name: %s, Date: %s, Time: %s, TimeZone: %s\n", spaceName, dateValue, timeValue, timeZoneAbbr)

	tz := timezone.New()
	timezones, err := tz.GetTimezones(timeZoneAbbr)
	if err != nil {
		log.Error(err)
		return nil, errors.New("time zone not available")
	}

	timeZoneIanaName := timezones[0]

	// eg: Feb 2, 2022 at 6:54pm IST
	dateTimeValue := fmt.Sprintf("%s at %s %s", dateValue, timeValue, timeZoneAbbr)

	const longFormTimeLayout = "Jan 2, 2006 at 3:04pm MST"

	startTime, err := time.Parse(longFormTimeLayout, dateTimeValue)
	if err != nil {
		log.Error(err)
		return nil, errors.New("invalid time format")
	}
	startDateTimeString := startTime.Format(time.RFC3339)

	endTime := startTime.Add(30 * time.Minute)
	endDateTimeString := endTime.Format(time.RFC3339)

	data := &ParsedUserTweetContent{
		SpaceName:           spaceName,
		StartDateTimeString: startDateTimeString,
		EndDateTimeString:   endDateTimeString,
		TimeZoneIanaName:    timeZoneIanaName,
	}
	return data, nil
}

func processTweetText(tweetText string) string {
	// eg: "@creategcalevent Kamran's Space | Feb 28, 2022 | 6:43 PM | IST"
	pattern := regexp.MustCompile(`(@creategcalevent).*`)
	matches := pattern.FindAllString(tweetText, -1)
	matchedText := matches[0] // [@creategcalevent Kamran's Space | Feb 28, 2022 | 6:43 PM | IST]

	tweetContents := strings.Split(matchedText, "@creategcalevent") // [@creategcalevent, Kamran's Space | Feb 28, 2022 | 6:43 PM | IST]
	tweetText = strings.TrimLeft(tweetContents[1], " ")             // Kamran's Space | Feb 28, 2022 | 6:43 PM | IST
	tweetText = strings.TrimSpace(tweetText)

	return tweetText
}

func processSpaceName(spaceName string) (string, error) {
	if strings.Contains(spaceName, " |") {
		// Need to split including space
		spaceName = strings.Split(spaceName, " |")[0]
	} else {
		// No need to split including space
		spaceName = strings.Split(spaceName, "|")[0]
	}
	return spaceName, nil
}

func processDateValue(dateValue string) (string, error) {
	dateValue = strings.TrimSpace(dateValue)
	if strings.Contains(dateValue, " |") {
		// Need to split including space
		dateValue = strings.Split(dateValue, " |")[0]
	} else {
		// No need to split including space
		dateValue = strings.Split(dateValue, "|")[0]
	}
	return dateValue, nil
}

func processTimeString(timeString string) (string, error) {
	timeValue := strings.TrimSpace(timeString)
	timeValue = strings.ToLower(timeValue)

	if (!strings.Contains(timeValue, "pm")) && (!strings.Contains(timeValue, "am")) {
		// incorrect time value provided
		return "", errors.New("incorrect time format")
	}

	if strings.Contains(timeValue, " |") {
		// Need to split including space
		timeValue = strings.Split(timeValue, " |")[0]
		timeValue = strings.ReplaceAll(timeValue, " ", "")
	} else {
		// No need to split including space
		timeValue = strings.Split(timeValue, "|")[0]
		timeValue = strings.ReplaceAll(timeValue, " ", "")
	}

	amOrPm := timeValue[len(timeValue)-2:]
	if !strings.Contains(timeValue, ":") {
		numericTimeValue := strings.Split(timeValue, timeValue[len(timeValue)-2:])[0]
		timeValue = fmt.Sprintf("%s:%s%s", numericTimeValue, "00", amOrPm)
	}

	return timeValue, nil
}

func processTimeZoneValue(timeZoneValue string) (string, error) {
	if strings.Contains(timeZoneValue, " |") {
		// Need to split including space
		timeZoneValue = strings.Split(timeZoneValue, " |")[0]
		timeZoneValue = strings.TrimSpace(timeZoneValue)
	} else {
		// No need to split including space
		timeZoneValue = strings.Split(timeZoneValue, "|")[0]
		timeZoneValue = strings.TrimSpace(timeZoneValue)
	}

	return timeZoneValue, nil
}
