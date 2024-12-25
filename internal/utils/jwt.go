package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

const (
	tokenExpiration = time.Hour * 8
	secretKey       = "cabbageofathing"
)

type Claims struct {
	jwt.RegisteredClaims
	SessionID string `json:"session_id"`
}

func NewJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
		},
		SessionID: uuid.New().String(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func fromString(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func GetSessionID(tokenString string) (string, error) {
	claims, err := fromString(tokenString)
	if err != nil {
		return "", err
	}

	return claims.SessionID, nil
}

func IsExpired(tokenString string) bool {
	claims, err := fromString(tokenString)
	if err != nil {
		return false
	}

	return claims.ExpiresAt.Before(time.Now())
}
