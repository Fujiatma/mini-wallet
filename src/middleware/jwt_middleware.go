package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Read secret key from environment variable
		secretKey := os.Getenv("SECRET_KEY")
		if secretKey == "" {
			fmt.Println("SECRET_KEY is not set in .env file")
		}

		tokenString := strings.Split(authHeader, " ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add the token to the request context
		ctx := context.WithValue(r.Context(), "token", token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
