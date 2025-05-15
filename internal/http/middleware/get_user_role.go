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
	//fmt.Println("Token string from cookie:", cookie.Value)

	if cookie.Value == "" {
		fmt.Println("22")
		return "", ErrNoAuthCookie
	}

	//fmt.Println("Cookie value:", cookie.Value)
	//fmt.Println("JWT_SECRET in middleware:", os.Getenv("JWT_SECRET"))
	//fmt.Println("Before jwt.ParseWithClaims")

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		fmt.Println("Token method:", token.Method.Alg())
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	//fmt.Println("Claims after parse:", claims)

	if err != nil {
		fmt.Println("JWT parse error:", err)
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
