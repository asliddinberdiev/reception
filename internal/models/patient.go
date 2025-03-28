package models

import "time"

type Patient struct {
	ID          string    `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	IsVerified  bool      `json:"is_validated"`
	Firstname   string    `json:"first_name"`
	Lastname    string    `json:"last_name"`
	Gender      string    `json:"gender"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PatientShort struct {
	ID        string `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

type PatientCreateInput struct {
	PhoneNumber string `json:"phone_number" validate:"required,numeric,len=12"`
	Firstname   string `json:"first_name" validate:"required,min=2,lowercase"`
	Lastname    string `json:"last_name" validate:"required,min=2,lowercase"`
	Gender      string `json:"gender" validate:"required,oneof=male female unknown"`
	Password    string `json:"password" validate:"required,min=6"`
}
