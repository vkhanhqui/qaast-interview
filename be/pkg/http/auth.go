package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(jwtKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				JSON(w, http.StatusUnauthorized, map[string]string{"error": "missing token"})
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtKey), nil
			})
			if err != nil {
				JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
				return
			}

			uid, ok := claims["user_id"].(string)
			if !ok {
				JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(w http.ResponseWriter, r *http.Request) string {
	id, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		JSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return ""
	}
	return id
}
