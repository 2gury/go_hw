package session

import (
	"fmt"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/models"
	"lectures-2022-1/05_web_app/99_hw/redditclone/tools/response"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secret = []byte("ASF9FA003NKFSA9")

func NewSession(usr models.User) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
		"user": response.Body{
			"username": usr.Username,
			"id":       usr.ID,
		},
	})
	strToken, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}
	return strToken, nil
}

func CheckSession(tokenValue string) (*models.User, error) {
	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unsupported signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userRawData := claims["user"].(map[string]interface{})
		return &models.User{
			ID:       userRawData["id"].(string),
			Username: userRawData["username"].(string),
		}, nil
	}
	return nil, fmt.Errorf("jwt token is not valid")
}

func GetTokenValue(rawToken string) string {
	return strings.Split(rawToken, " ")[1]
}
