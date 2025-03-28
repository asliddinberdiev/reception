package helper

import "strconv"

func ParseInt(value string, defaultValue int) int {
	if parsedValue, err := strconv.Atoi(value); err == nil {
		return parsedValue
	}
	return defaultValue
}
