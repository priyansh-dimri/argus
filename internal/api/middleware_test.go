package api_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/priyansh-dimri/argus/internal/api"
)

type mockAuthStore struct {
	ProjectID string
	Err       error
}

func (m *mockAuthStore) GetProjectIDByKey(ctx context.Context, apiKey string) (string, error) {
	return m.ProjectID, m.Err
}

func TestAuthSDK(t *testing.T) {
	store := &mockAuthStore{ProjectID: "proj_123"}
	mw := api.NewMiddleware(store)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := api.GetProjectID(r.Context())
		if !ok || id != "proj_123" {
			t.Errorf("Context project_id mismatch: got %v", id)
		}
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Valid Key", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/analyze", nil)
		req.Header.Set("Authorization", "Bearer argus_valid")
		rec := httptest.NewRecorder()

		mw.AuthSDK(nextHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", rec.Code)
		}
	})

	t.Run("Missing Header", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/analyze", nil)
		rec := httptest.NewRecorder()

		mw.AuthSDK(nextHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", rec.Code)
		}
	})

	t.Run("Invalid Key Format", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/analyze", nil)
		req.Header.Set("Authorization", "Basic user:pass")
		rec := httptest.NewRecorder()

		mw.AuthSDK(nextHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", rec.Code)
		}
	})

	t.Run("Detect database rejected", func(t *testing.T) {
		rejectStore := &mockAuthStore{Err: errors.New("not found")}
		mwReject := api.NewMiddleware(rejectStore)

		req := httptest.NewRequest("POST", "/analyze", nil)
		req.Header.Set("Authorization", "Bearer argus_bad")
		rec := httptest.NewRecorder()

		mwReject.AuthSDK(nextHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", rec.Code)
		}
	})
}

func TestAuthDashboard(t *testing.T) {
	secret := "super-secret-jwt-key"
	os.Setenv("SUPABASE_JWT_SECRET", secret)
	defer os.Unsetenv("SUPABASE_JWT_SECRET")

	store := &mockAuthStore{}
	mw := api.NewMiddleware(store)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := api.GetUserID(r.Context())
		if !ok || id != "user_123" {
			t.Errorf("Context user_id mismatch: got %v", id)
		}
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Valid JWT", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "user_123",
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(secret))

		req := httptest.NewRequest("GET", "/projects", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()

		mw.AuthDashboard(nextHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", rec.Code)
		}
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "user_123",
		})
		tokenString, _ := token.SignedString([]byte("wrong-secret"))

		req := httptest.NewRequest("GET", "/projects", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()

		mw.AuthDashboard(nextHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", rec.Code)
		}
	})

	t.Run("Detect error missing token header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/projects", nil)
		rec := httptest.NewRecorder()

		mw.AuthDashboard(nextHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d", rec.Code)
		}
		if rec.Body.String() != "Unauthorized: Missing Token\n" {
			t.Errorf("Unexpected body: %q", rec.Body.String())
		}
	})

	t.Run("Detect unexpected signing method like RSA", func(t *testing.T) {
		key, _ := rsa.GenerateKey(rand.Reader, 2048)
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"sub": "user_123",
		})
		tokenString, _ := token.SignedString(key)

		req := httptest.NewRequest("GET", "/projects", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()

		mw.AuthDashboard(nextHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d", rec.Code)
		}
	})

	t.Run("Detect missing subject claim", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"role": "admin",
			"exp":  time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(secret))

		req := httptest.NewRequest("GET", "/projects", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rec := httptest.NewRecorder()

		mw.AuthDashboard(nextHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d", rec.Code)
		}
		if rec.Body.String() != "Unauthorized: Token missing subject\n" {
			t.Errorf("Unexpected body: %q", rec.Body.String())
		}
	})
}

func TestCORS(t *testing.T) {
	store := &mockAuthStore{}
	mw := api.NewMiddleware(store)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Preflight OPTIONS returns 204 with CORS headers", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/projects", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")
		req.Header.Set("Access-Control-Request-Headers", "Authorization")

		rec := httptest.NewRecorder()
		mw.CORS(nextHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("Expected 204, got %d", rec.Code)
		}
		if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
			t.Errorf("Missing/wrong Access-Control-Allow-Origin: %q", rec.Header().Get("Access-Control-Allow-Origin"))
		}
		if rec.Header().Get("Access-Control-Allow-Methods") != "GET, POST, OPTIONS, DELETE, PATCH" {
			t.Errorf("Missing/wrong Access-Control-Allow-Methods: %q", rec.Header().Get("Access-Control-Allow-Methods"))
		}
		if rec.Header().Get("Access-Control-Allow-Headers") != "Authorization, Content-Type" {
			t.Errorf("Missing/wrong Access-Control-Allow-Headers: %q", rec.Header().Get("Access-Control-Allow-Headers"))
		}
	})

	t.Run("Non-OPTIONS passes through to next handler", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/projects", nil)
		recorder := httptest.NewRecorder()

		mw.CORS(nextHandler).ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Errorf("Expected 200 from next handler, got %d", recorder.Code)
		}
	})

	t.Run("Sets Vary header", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/projects", nil)
		request.Header.Set("Origin", "http://localhost:3000")
		recorder := httptest.NewRecorder()

		mw.CORS(nextHandler).ServeHTTP(recorder, request)

		if recorder.Header().Get("Vary") != "Origin" {
			t.Errorf("Expected Vary: Origin, got %q", recorder.Header().Get("Vary"))
		}
	})
}
