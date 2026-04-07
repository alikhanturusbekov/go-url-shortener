package authorization

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	cookieName = "auth_token"
)

// AuthMiddleware provides JWT-based authentication middleware
func AuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cookie, err := r.Cookie(cookieName); err == nil {
				if claims, err := ParseToken(cookie.Value, jwtKey); err == nil && claims.UserID != "" {
					ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			userID := uuid.NewString()

			token, err := createToken(userID, jwtKey)
			if err != nil {
				http.Error(w, "could not create token", http.StatusInternalServerError)
				return
			}

			setAuthCookie(w, token)

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// createToken generates a signed JWT for the given user ID
func createToken(userID string, jwtKey []byte) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// setAuthCookie sets the authentication cookie
func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Expires:  time.Now().Add(tokenTTL),
		MaxAge:   int(tokenTTL.Seconds()),
		HttpOnly: true,
	})
}
