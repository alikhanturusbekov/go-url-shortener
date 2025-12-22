package authorization

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const userIDContextKey contextKey = "userID"

const (
	cookieName = "auth_token"
	tokenTTL   = 30 * 24 * time.Hour
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)

	return userID, ok
}

func AuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cookie, err := r.Cookie(cookieName); err == nil {
				if claims, err := parseToken(cookie.Value, jwtKey); err == nil && claims.UserID != "" {
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

func parseToken(tokenStr string, jwtKey []byte) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Expires:  time.Now().Add(tokenTTL),
		MaxAge:   int(tokenTTL.Seconds()),
		HttpOnly: true,
	})
}
