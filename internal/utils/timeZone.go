package utils

import "github.com/pkg/errors"

type TimeZone struct {
	IANAName string
}

var availableTimeZones = map[string]TimeZone{
	"IST": {
		IANAName: "Asia/Calcutta",
	},
}

func getTimeZoneIANAName(timeZone string) (string, error) {
	if availableTimeZones[timeZone].IANAName == "" {
		return "", errors.Errorf("cant find IANA Name for timezone %s", timeZone)
	}
	return availableTimeZones[timeZone].IANAName, nil
}