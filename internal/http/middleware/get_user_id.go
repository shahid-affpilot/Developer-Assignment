package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/config"
)

func GetUserID(r *http.Request) (string, error) {
	// Get auth cookie
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return "", ErrNoAuthCookie
	}

	if cookie.Value == "" {
		return "", ErrNoAuthCookie
	}

	cnf := config.GetConfig()

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		fmt.Println("Token method:", token.Method.Alg())
		return []byte(cnf.JWT.Secret), nil
	})

	//fmt.Println("Claims after parse:", claims)

	if err != nil {
		return "", ErrInvalidToken
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	userType, ok := claims["user_id"].(string)
	if !ok || userType == "" {
		return "", ErrNoUserType
	}

	return userType, nil
}
