package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
)

func ParseTweetText() {
	// <space_name> | <date> | <time> | <time_zone>
	strSlice := strings.SplitAfter("Kamran's Space | Jan 28, 2022 | 6:43 PM | IST", "|")
	if len(strSlice) != 4 {
		log.Warning("invalid string format")
		os.Exit(1)
	}

	spaceName, err := processSpaceName(strSlice[0])
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	dateValue, err := processDateValue(strSlice[1])
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	timeValue, err := processTimeString(strSlice[2])
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	timeZoneValue, err := processTimeZoneValue(strSlice[3])
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("Space Name: %s, Date: %s, Time: %s, TimeZone: %s\n", spaceName, dateValue, timeValue, timeZoneValue)

	ianaName, err := getIanaName("IST")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("IANA Name:", ianaName)
}

func processSpaceName(spaceName string) (string, error) {
	if strings.Contains(spaceName, " |") {
		// fmt.Println("Need to split including space")
		spaceName = strings.Split(spaceName, " |")[0]
	} else {
		// fmt.Println("No need to split including space")
		spaceName = strings.Split(spaceName, "|")[0]
	}
	return spaceName, nil
}

func processDateValue(dateValue string) (string, error) {
	if strings.Contains(dateValue, " |") {
		// fmt.Println("Need to split including space")
		dateValue = strings.Split(dateValue, " |")[0]
	} else {
		// fmt.Println("No need to split including space")
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
		// fmt.Println("Need to split including space")
		timeValue = strings.Split(timeValue, " |")[0]
		timeValue = strings.ReplaceAll(timeValue, " ", "")
	} else {
		// fmt.Println("No need to split including space")
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
		// fmt.Println("Need to split including space")
		timeZoneValue = strings.Split(timeZoneValue, " |")[0]
		timeZoneValue = strings.TrimSpace(timeZoneValue)
	} else {
		// fmt.Println("No need to split including space")
		timeZoneValue = strings.Split(timeZoneValue, "|")[0]
		timeZoneValue = strings.TrimSpace(timeZoneValue)
	}

	return timeZoneValue, nil
}
