package util

import (
	"errors"
	"go-restfull-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.Env.JWTSecret) // store securely, e.g., in env var

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // expires in 3 days
		"iat":     time.Now().Unix(),                     // issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return jwtSecret, nil
	})
}

