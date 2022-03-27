package log

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/mgutz/ansi"
)

// CallingDepth : Depth of stack to go back to get proper line and file function called from
const (
	CallingDepth = 1

	// Background colors to use
	ErrorColor   = "red"
	WarningColor = "yellow"
	InfoColor    = "green"
	FatalColor   = "magenta+b"

	// Base error string
	BaseString = ` * `
)

// Info : Wrapper around glog Info
func Info(args ...interface{}) {
	glog.InfoDepth(CallingDepth, getColorString(InfoColor), args)
}

// Info : Wrapper around glog Info
func Infof(formatString string, args ...interface{}) {
	msg := fmt.Sprintf(formatString, args...)
	glog.InfoDepth(CallingDepth, getColorString(InfoColor), msg)
}

// Warning : Wrapper around glog Warning
func Warning(args ...interface{}) {
	glog.WarningDepth(CallingDepth, getColorString(WarningColor), args)
}

// Warning : Wrapper around glog Warning
func Warningf(formatString string, args ...interface{}) {
	msg := fmt.Sprintf(formatString, args...)
	glog.WarningDepth(CallingDepth, getColorString(WarningColor), msg)
}

// Error : Wrapper around glog Error
func Error(args ...interface{}) {
	colorString := getColorString(ErrorColor)
	glog.ErrorDepth(CallingDepth, colorString, args)
}

// Errorf : Wrapper around glog Error
func Errorf(formatString string, args ...interface{}) {
	msg := fmt.Sprintf(formatString, args...)
	glog.ErrorDepth(CallingDepth, getColorString(ErrorColor), msg)
}

// Fatal : Wrapper around glog Fatal
func Fatal(args ...interface{}) {
	glog.FatalDepth(CallingDepth, getColorString(FatalColor), args)
}

func getColorString(backgroundColor string) string {
	baseStr := ` * `
	str := fmt.Sprintf(`|%s| `, ansi.Color(baseStr, "white:"+backgroundColor))
	return str
}
