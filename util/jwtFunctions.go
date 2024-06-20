package util

import (
	"blogbackend/models"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

const (
	tokenTTL = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	IsAdmin bool `json:"is_admin"`
}

func GenerateJWT(user models.User) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}
	signingKey := os.Getenv("SIGNING_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			Issuer:    strings.TrimSpace(strconv.Itoa(int(user.Id))),
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.IsAdmin})
	return token.SignedString([]byte(signingKey))
}

func ParseJWT(cookie string) (string, bool, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	signingKey := os.Getenv("SIGNING_KEY")

	token, err := jwt.ParseWithClaims(cookie, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", false, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", false, errors.New("token claims are not of type *tokenClaims")
	}

	userId, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		return "", false, errors.New("invalid user ID in token claims")
	}

	return strconv.Itoa(userId), claims.IsAdmin, nil
}
