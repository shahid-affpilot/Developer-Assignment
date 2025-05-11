package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoAuthCookie = errors.New("no auth cookie found")
	ErrInvalidToken = errors.New("invalid auth token")
	ErrNoUserType   = errors.New("no user type in token")
)

func GetUserRole(r *http.Request) (string, error) {
	// Get auth cookie
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		fmt.Println("11")
		return "", ErrNoAuthCookie
	}

	if cookie.Value == "" {
		fmt.Println("22")
		return "", ErrNoAuthCookie
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("33")
			return nil, ErrInvalidToken
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println("44")
		return "", ErrInvalidToken
	}

	if !token.Valid {
		fmt.Println("55")
		return "", ErrInvalidToken
	}

	// Extract user type with type assertion
	userType, ok := claims["user_type"].(string)
	if !ok || userType == "" {
		fmt.Println("66")
		return "", ErrNoUserType
	}

	return userType, nil
}
