// Package utils provides utility functions for the application.
// It contains helpers for JWT authentication, rendering templates, and other common operations.
package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(id string, name string) (string, error) {
	segredo := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":   id,
			"name": name,
			"exp":  time.Now().Add(time.Hour * 1).Unix(),
		})
	tokenString, err := token.SignedString([]byte(segredo))
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func VerifyJwtToken(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return err
	}
	return nil
}

func GetTokenClaims(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
