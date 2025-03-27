package helper

import "strings"

func ErrorIs(err error, target string) bool {

	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), target)
}
