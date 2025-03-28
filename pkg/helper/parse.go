package helper

import (
	"errors"
	"strconv"
	"time"
)

func ParseInt(value string, defaultValue int) int {
	if parsedValue, err := strconv.Atoi(value); err == nil {
		return parsedValue
	}
	return defaultValue
}

func StringToTime(timeStr, format string) (time.Time, error) {
	timeFormat, err := time.Parse(format, timeStr)
	if err != nil {
		return time.Time{}, errors.New("invalid time format")
	}

	return timeFormat, nil
}
