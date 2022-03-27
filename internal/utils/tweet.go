package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/skamranahmed/markeventas/pkg/log"
	"github.com/tkuchiki/go-timezone"
)

var (
	ErrInvalidDateOrTime error = errors.New("invalid date or time")
	ErrInvalidMonth      error = errors.New("invalid month")
)

type ParsedUserTweetContent struct {
	SpaceName           string
	StartDateTimeString string
	EndDateTimeString   string
	TimeZoneIanaName    string
}

func ParseTweetText(tweetText string) (*ParsedUserTweetContent, error) {
	/*
		@markeventas <space_name> | <date> | <time> | <time_zone>

		<date> is of the format <month_name date, year>, eg: September 15, 2022
	*/

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
		log.Error("error in processing date value, got: %s, err: %v", strSlice[1], err)
		return nil, err
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

	startDateTimeString, endDateTimeString, err := processDateTimeValue(dateTimeValue)
	if err != nil {
		return nil, err
	}

	data := &ParsedUserTweetContent{
		SpaceName:           spaceName,
		StartDateTimeString: startDateTimeString,
		EndDateTimeString:   endDateTimeString,
		TimeZoneIanaName:    timeZoneIanaName,
	}
	return data, nil
}

func processTweetText(tweetText string) string {
	// eg: "@markeventas Kamran's Space | Feb 28, 2022 | 6:43 PM | IST"
	pattern := regexp.MustCompile(`(@markeventas).*`)
	matches := pattern.FindAllString(tweetText, -1)
	matchedText := matches[0] // [@markeventas Kamran's Space | Feb 28, 2022 | 6:43 PM | IST]

	tweetContents := strings.Split(matchedText, "@markeventas") // [@markeventas, Kamran's Space | Feb 28, 2022 | 6:43 PM | IST]
	tweetText = strings.TrimLeft(tweetContents[1], " ")         // Kamran's Space | Feb 28, 2022 | 6:43 PM | IST
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

	monthName := strings.TrimSpace(digitPrefix(dateValue))
	formattedMonthName, err := processMonthName(monthName)
	if err != nil {
		return "", err
	}

	// replace the user entered monthName with the formattedMonthName in the dateValue string
	dateValue = strings.Replace(dateValue, monthName, formattedMonthName, -1)

	return dateValue, nil
}

func processMonthName(monthName string) (string, error) {
	// change the month name into the format which is accepted by go standard package
	formattedMonthName, ok := months[strings.ToUpper(monthName)]
	if !ok {
		// if the month name is not present in the map
		// probably the user has entered some invalid month name
		return "", ErrInvalidMonth
	}
	return formattedMonthName, nil
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

func processDateTimeValue(dateTimeValue string) (string, string, error) {
	const longFormTimeLayout = "Jan 2, 2006 at 3:04pm MST"

	startTime, err := time.Parse(longFormTimeLayout, dateTimeValue)
	if err != nil {
		log.Error("error in processing date-time value, got: %s, err: %v", dateTimeValue, err)
		return "", "", ErrInvalidDateOrTime
	}

	// start time of the event
	startDateTimeString := startTime.Format(time.RFC3339)

	// end time of the event = start time + 30 mins
	endTime := startTime.Add(30 * time.Minute)
	endDateTimeString := endTime.Format(time.RFC3339)

	return startDateTimeString, endDateTimeString, nil
}

// digitPrefix : returns the string before the occurence of the first digit in the passed string
func digitPrefix(s string) string {
	for i, r := range s {
		if unicode.IsDigit(r) {
			return s[:i]
		}
	}
	return s
}
