package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func GenerateToken(claims jwt.Claims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}

	return tokenString, nil
}

func ParseToken(token string, secretKey string, claims jwt.Claims) error {
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to parse token")
	}

	if err = jwtToken.Claims.Valid(); err != nil {
		return errors.Wrap(err, "failed to validate token")
	}

	return nil
}
