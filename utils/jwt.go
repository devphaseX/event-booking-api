package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const superSecret = "supersecret_key"

func GenerateJwtToken(userId int64, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"userId": userId,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(superSecret))
}

func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method ")
		}

		return []byte(superSecret), nil
	})

	if err != nil {
		return 0, err
	}

	valid := parsedToken.Valid

	if !valid {
		return 0, errors.New("invalid or expired token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid token claim")
	}

	userId := int64(claims["userId"].(float64))
	return userId, nil
}
