package main

import (
	"net/http"
	"strings"

	"github.com/NoobsBucket/iqra-initi/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Fields(r.Header.Get("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			adminAuthError(w, http.StatusUnauthorized, "bearer token required")
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(app.config.jwt.secret), nil
		}, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil || token == nil || !token.Valid {
			adminAuthError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		userID, err := claims.GetSubject()
		if err != nil || userID == "" {
			adminAuthError(w, http.StatusUnauthorized, "token subject is missing")
			return
		}

		user, err := app.authStore.GetUserByID(r.Context(), userID)
		if err != nil {
			adminAuthError(w, http.StatusUnauthorized, "user is not authorized")
			return
		}
		if user.Role != "admin" {
			adminAuthError(w, http.StatusForbidden, "admin access required")
			return
		}

		r = r.WithContext(auth.WithUserID(r.Context(), user.ID))
		next.ServeHTTP(w, r)
	})
}

func adminAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"` + message + `"}`))
}
