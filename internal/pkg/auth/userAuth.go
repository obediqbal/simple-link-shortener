package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserClaim struct {
	jwt.StandardClaims
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

func GenerateJWT(username string, displayName string) (string, error) {
	claims := UserClaim{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
		Username:    username,
		DisplayName: displayName,
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
		if claims.ExpiresAt < time.Now().Unix() {
			return nil, fmt.Errorf("Expired token")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
}
