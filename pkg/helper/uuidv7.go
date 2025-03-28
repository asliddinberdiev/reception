package helper

import "github.com/google/uuid"

func NewV7ID() string {
	id, _ := uuid.NewV7()
	return id.String()
}
