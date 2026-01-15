package main

import (
	"context" // to store userid inside the request
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

// authenticate is a middleware the validates the JWT token
// It wraps a standard http.HandlerFunc and returns a new http.HandlerFunc
func (app *application) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// removing bearer prefix to just get the token string
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := headerParts[1]

		// parse and validate the token
		// we will pass a callback function that return the secret key to the parser
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// validating that signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// checking if parsing failed
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// getting the user id from claims
		userIdFloat, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}
		userId := int(userIdFloat)

		// adding user id to request context
		ctx := context.WithValue(r.Context(), "userId", userId)

		next(w, r.WithContext(ctx))
	}
}