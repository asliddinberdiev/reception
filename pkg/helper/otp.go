package helper

import "math/rand"

var numberRunes = []rune("0123456789")

func RandNumberStringRunes(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = numberRunes[rand.Intn(len(numberRunes))]
	}
	return string(b)
}
