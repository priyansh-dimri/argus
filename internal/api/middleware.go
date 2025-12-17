package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/priyansh-dimri/argus/pkg/logger"
)

type AuthStore interface {
	GetProjectIDByKey(ctx context.Context, apiKey string) (string, error)
}

type Middleware struct {
	Store AuthStore
	JWTSecret string
}

func NewMiddleware(store AuthStore) *Middleware {
	secret := os.Getenv("SUPABASE_JWT_SECRET")
	if secret == "" {
		logger.Warn("SUPABASE_JWT_SECRET not set; dashboard auth will fail")
	}
	return &Middleware{
		Store:     store,
		JWTSecret: secret,
	}
}

func (m *Middleware) AuthSDK(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing API Key", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized: Invalid Header Format", http.StatusUnauthorized)
			return
		}

		apiKey := parts[1]
		projectID, err := m.Store.GetProjectIDByKey(r.Context(), apiKey)
		if err != nil {
			logger.Error("AuthSDK failed", err)
			http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
			return
		}

		ctx := WithProjectID(r.Context(), projectID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *Middleware) AuthDashboard(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing Token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			logger.Error("Dashboard Auth Failed", err)
			http.Error(w, "Unauthorized: Invalid Token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(string); ok {
				ctx := WithUserID(r.Context(), sub)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		http.Error(w, "Unauthorized: Token missing subject", http.StatusUnauthorized)
	}
}
