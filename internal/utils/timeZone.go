package utils

import "github.com/pkg/errors"

const longFormTimeLayout = "Jan 2, 2006 at 3:04pm MST"

type TimeZone struct {
	IANAName string
}

var availableTimeZones = map[string]TimeZone{
	"IST": {
		IANAName: "Asia/Calcutta",
	},
}

func getIanaName(timeZone string) (string, error) {
	if availableTimeZones[timeZone].IANAName == "" {
		return "", errors.Errorf("cant find IANA Name for timezone %s", timeZone)
	}
	return availableTimeZones[timeZone].IANAName, nil
}
