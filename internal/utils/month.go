package utils

const (
	January   string = "Jan"
	February  string = "Feb"
	March     string = "Mar"
	April     string = "Apr"
	May       string = "May"
	June      string = "Jun"
	July      string = "Jul"
	August    string = "Aug"
	September string = "Sep"
	October   string = "Oct"
	November  string = "Nov"
	December  string = "Dec"
)

var (
	months = map[string]string{
		// January
		"JAN":     January,
		"JANUARY": January,

		// February
		"FEB":      February,
		"FEBRUARY": February,

		// March
		"MAR":   March,
		"MARCH": March,

		// April
		"APR":   April,
		"APRIL": April,

		// May
		"MAY": May,

		// June
		"JUN":  June,
		"JUNE": June,

		// July
		"JUL":  July,
		"JULY": July,

		// August
		"AUG":    August,
		"AGST":   August,
		"AUGST":  August,
		"AUGUST": August,

		// September
		"SEP":       September,
		"SEPT":      September,
		"SEPTEMBER": September,

		// October
		"OCT":     October,
		"OCTOBER": October,

		// November
		"NOV":      November,
		"NOVEMBER": November,

		// December
		"DEC":      December,
		"DECEMBER": December,
	}
)