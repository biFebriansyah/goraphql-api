package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	Id    string `json:"_id"`
	Admin bool   `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJwt(userId string, role bool) (string, error) {
	claims := jwtClaims{
		Id:    userId,
		Admin: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "goraphql",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute + 1)),
		},
	}

	secrets := os.Getenv("JWT_SECRETS")
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return tokens.SignedString([]byte(secrets))
}

func ParseJwt(token string) (*jwtClaims, error) {
	secrets := os.Getenv("JWT_SECRETS")
	jwtData, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secrets), nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail parse jwt: %w", err)
	}

	claimData := jwtData.Claims.(*jwtClaims)
	return claimData, nil
}
