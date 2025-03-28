package models

import "github.com/golang-jwt/jwt"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type PatientClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

type OtpClaims struct {
	PhoneNumber string `json:"phone_number"`
	jwt.StandardClaims
}

type VerifyInput struct {
	Otp string `json:"otp" validate:"required,numeric,len=8"`
}

type RegisterInput struct {
	PhoneNumber string `json:"phone_number" validate:"required,numeric,len=12"`
}

type LoginInput struct {
	PhoneNumber string `json:"phone_number" validate:"required,numeric,len=12"`
	Password    string `json:"password" validate:"required,min=6"`
}
