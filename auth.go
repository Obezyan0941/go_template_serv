package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret-key")

func generateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,                                // Идентификатор пользователя
		"exp": time.Now().Add(24 * time.Hour).Unix(), // Срок действия (24 часа)
	})

	return token.SignedString(jwtSecret)
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
}
