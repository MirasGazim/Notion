package jwt

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"notion/internal/handlers/middleware/ctx"
	"notion/internal/service"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	ErrInvalidToken = errors.New("invalid token")
)

const (
	signingKey = "grkjk#4#%35FSFJlja#4353KSFjH"
)

func AuthMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			auth_header := r.Header.Get("Authorization")

			parts := strings.Split(auth_header, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Error("invalid auth header", slog.String("error", "Unauthorized"))
				http.Error(w, "invalid auth header", http.StatusUnauthorized)
				return
			}
			tokenstring := parts[1]
			claims, err := ParseToken(tokenstring, signingKey)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ctx.UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func ParseToken(tokenstring string, secret string) (*service.TokenClaims, error) {
	claims := &service.TokenClaims{}
	funcc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	}
	token, err := jwt.ParseWithClaims(tokenstring, claims, funcc)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
