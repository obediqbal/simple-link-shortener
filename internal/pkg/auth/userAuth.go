package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserClaim struct {
	jwt.StandardClaims
	Username string              `json:"username"`
	History  []map[string]string `json:"history"`
}

func GenerateJWT(username string) (string, error) {
	// db := database.New()
	var history []map[string]string // TO BE IMPLEMENTED

	claims := UserClaim{
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
		Username: username,
		History:  history,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
}
