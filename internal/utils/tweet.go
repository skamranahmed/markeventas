package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
)

type ParsedUserTweetContent struct {
	SpaceName string
	DateTimeString string
	TimeZoneIanaName string
}

func ParseTweetText(tweetText string) (*ParsedUserTweetContent, error) {
	// <space_name> | <date> | <time> | <time_zone>
	strSlice := strings.SplitAfter("Kamran's Space | Jan 28, 2022 | 6:43 PM | IST", "|")
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

	timeZoneValue, err := processTimeZoneValue(strSlice[3])
	if err != nil {
		log.Error(err)
		return nil, errors.New("invalid time zone format")
	}

	log.Info("Space Name: %s, Date: %s, Time: %s, TimeZone: %s\n", spaceName, dateValue, timeValue, timeZoneValue)

	timeZoneIanaName, err := getIanaName("IST")
	if err != nil {
		log.Error(err)
		return nil, errors.New("time zone not available")
	}

	// eg: Feb 2, 2022 at 6:54pm IST
	dateTimeValue := fmt.Sprintf("%s at %s %s", dateValue, timeValue, timeZoneValue)

	t, err := time.Parse(longFormTimeLayout, dateTimeValue)
	if err != nil {
		log.Error(err)
		return nil, errors.New("invalid time format")
	}

	dateTimeString := t.Format(time.RFC3339)

	data := &ParsedUserTweetContent{
		SpaceName: spaceName,
		DateTimeString: dateTimeString,
		TimeZoneIanaName: timeZoneIanaName,
	}
	return data, nil
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
