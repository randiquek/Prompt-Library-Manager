package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get auth header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// auth header format: "Bearer <token>", split to get only token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header fomrat. Use: Bearer <token>", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// parse and validate token

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// return secret key for validation
			secret := os.Getenv("JWT_SECRET")
			return []byte(secret), nil
		})

		if err != nil {
			log.Println("Token parse error", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// check for valid token
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// extract claims (username, expiration, etc)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// add username to request context for handler usage
			username := claims["username"].(string)
			ctx := context.WithValue(r.Context(), "username", username)
			r = r.WithContext(ctx)
			
			log.Printf("Authenticated user: %s", username)
		}

		// Token is valid, call next handler
		next.ServeHTTP(w, r)
	})
}