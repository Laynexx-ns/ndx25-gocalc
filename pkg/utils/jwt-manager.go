package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"ndx/pkg/config"
	"ndx/pkg/logger"
	"time"
)

var jwtSecret []byte

func NewToken(email string) string {
	if len(jwtSecret) == 0 {
		jwtSecret = []byte(config.NewConfig().JwtSecretKey)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(jwtSecret)
	if err != nil {
		logger.L().Logf(0, "error getting token string | err: %v", err)
		return ""
	}
	return tokenString
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidId
	}

	return token, nil
}
