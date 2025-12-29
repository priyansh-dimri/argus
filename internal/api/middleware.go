package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/priyansh-dimri/argus/pkg/logger"
)

type AuthStore interface {
	GetProjectIDByKey(ctx context.Context, apiKey string) (string, error)
}

type Middleware struct {
	Store     AuthStore
	JWTSecret string
}

func NewMiddleware(store AuthStore) *Middleware {
	logger.Info("Initializing middleware", "component", "middleware")

	secret := os.Getenv("SUPABASE_JWT_SECRET")
	if secret == "" {
		logger.Warn("SUPABASE_JWT_SECRET not set; dashboard auth will fail",
			"component", "middleware",
		)
	} else {
		logger.Info("JWT secret loaded successfully",
			"component", "middleware",
			"secret_length", len(secret),
		)
	}

	logger.Info("Middleware initialized successfully", "component", "middleware")
	return &Middleware{
		Store:     store,
		JWTSecret: secret,
	}
}

func (m *Middleware) AuthSDK(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Info("SDK authentication started",
			"component", "middleware",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Warn("SDK auth failed: missing authorization header",
				"component", "middleware",
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
			http.Error(w, "Unauthorized: Missing API Key", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("SDK auth failed: invalid header format",
				"component", "middleware",
				"path", r.URL.Path,
				"header_parts", len(parts),
				"remote_addr", r.RemoteAddr,
			)
			http.Error(w, "Unauthorized: Invalid Header Format", http.StatusUnauthorized)
			return
		}

		apiKey := parts[1]
		logger.Info("Validating API key",
			"component", "middleware",
			"key_length", len(apiKey),
		)

		projectID, err := m.Store.GetProjectIDByKey(r.Context(), apiKey)
		if err != nil {
			logger.Error("SDK auth failed: invalid API key", err,
				"component", "middleware",
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
			http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
			return
		}

		logger.Info("SDK authentication successful",
			"component", "middleware",
			"project_id", projectID,
			"duration_ms", time.Since(start).Milliseconds(),
			"path", r.URL.Path,
		)

		ctx := WithProjectID(r.Context(), projectID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *Middleware) AuthDashboard(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Info("Dashboard authentication started",
			"component", "middleware",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Warn("Dashboard auth failed: missing token",
				"component", "middleware",
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
			http.Error(w, "Unauthorized: Missing Token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		logger.Info("Parsing JWT token",
			"component", "middleware",
			"token_length", len(tokenString),
		)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Warn("Unexpected JWT signing method",
					"component", "middleware",
					"method", token.Header["alg"],
				)
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			logger.Error("Dashboard auth failed: invalid token", err,
				"component", "middleware",
				"token_valid", token != nil && token.Valid,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
			http.Error(w, "Unauthorized: Invalid Token", http.StatusUnauthorized)
			return
		}

		logger.Info("JWT token parsed successfully",
			"component", "middleware",
			"token_valid", token.Valid,
		)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(string); ok {
				logger.Info("Dashboard authentication successful",
					"component", "middleware",
					"user_id", sub,
					"duration_ms", time.Since(start).Milliseconds(),
					"path", r.URL.Path,
				)
				ctx := WithUserID(r.Context(), sub)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			logger.Warn("JWT token missing subject claim",
				"component", "middleware",
				"claims", claims,
			)
		} else {
			logger.Warn("Failed to extract JWT claims",
				"component", "middleware",
			)
		}

		http.Error(w, "Unauthorized: Token missing subject", http.StatusUnauthorized)
	}
}

func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		logger.Info("CORS middleware processing request",
			"component", "middleware",
			"method", r.Method,
			"path", r.URL.Path,
			"origin", origin,
			"remote_addr", r.RemoteAddr,
		)

		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			logger.Info("Set CORS origin header",
				"component", "middleware",
				"origin", origin,
			)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Vary", "Origin")

		if r.Method == http.MethodOptions {
			logger.Info("Handling OPTIONS preflight request",
				"component", "middleware",
				"path", r.URL.Path,
				"origin", origin,
			)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		logger.Info("CORS headers set, continuing to handler",
			"component", "middleware",
			"method", r.Method,
			"path", r.URL.Path,
		)

		next.ServeHTTP(w, r)
	})
}
