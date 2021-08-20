package functions

import (
	"fmt"
	"shortener-app/models"
	"time"

	"github.com/golang-jwt/jwt"
)

func NewToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"email":    user.Email,
		"name":     user.Name,
		"expireAt": time.Now().Format(time.RFC3339),
	})

	signedToken, err := token.SignedString([]byte(GetEnv("JWT_SIGN")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func IsValidToken(token string) (bool, jwt.MapClaims) {
	validatedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(GetEnv("JWT_SIGN")), nil
	})

	if err != nil {
		return false, nil
	}

	claims, ok := validatedToken.Claims.(jwt.MapClaims)

	month := (time.Hour * 24) * 30
	expireAt, _ := time.Parse(time.RFC3339, string(claims["expireAt"].(string)))
	timeElapsed := time.Since(expireAt)

	if ok && validatedToken.Valid && timeElapsed < month {
		return true, claims
	}

	return false, nil
}
