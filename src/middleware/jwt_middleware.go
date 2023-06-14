package middleware

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

// Define JWTClaimsContextKey as the key to access JWT claims in the context
const JWTClaimsContextKey = "jwtClaims"

// JWTClaims represents the claims in the JWT token
type JWTClaims struct {
	CustomerXID string `json:"customer_xid"`
	jwt.StandardClaims
}

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
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
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

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add the claims to the request context
		ctx := context.WithValue(r.Context(), JWTClaimsContextKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
