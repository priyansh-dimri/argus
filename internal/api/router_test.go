package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/priyansh-dimri/argus/pkg/protocol"
)

func TestRouter(t *testing.T) {
	mock := newMockAnalyzer(protocol.AnalysisResponse{}, nil)
	store := &mockStore{
		MockProjectID: "project_123",
	}
	apiHandler := NewAPI(mock, store)

	testSecret := "test-secret"
	mw := &Middleware{
		Store:     store,
		JWTSecret: testSecret,
	}

	router := NewRouter(apiHandler, mw)

	validJWT, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user_123",
		"exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte(testSecret))

	tests := []struct {
		name           string
		method         string
		path           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Valid + Authenticated POST /analyze",
			method:         http.MethodPost,
			path:           "/analyze",
			authHeader:     "Bearer argus_valid_key",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid + Unauthenticated POST /analyze",
			method:         http.MethodPost,
			path:           "/analyze",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Method GET /analyze",
			method:         http.MethodGet,
			path:           "/analyze",
			authHeader:     "Bearer argus_valid_key",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Unknown Route POST /random",
			method:         http.MethodPost,
			path:           "/random",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Valid + Authenticated GET /projects",
			method:         http.MethodGet,
			path:           "/projects",
			authHeader:     "Bearer " + validJWT,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid + Authenticated POST /projects",
			method:         http.MethodPost,
			path:           "/projects",
			authHeader:     "Bearer " + validJWT,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid + Unauthenticated GET /projects",
			method:         http.MethodGet,
			path:           "/projects",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid + Unauthenticated POST /projects",
			method:         http.MethodPost,
			path:           "/projects",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body *strings.Reader
			if tc.method == http.MethodPost {
				switch tc.path {
				case "/analyze":
					body = strings.NewReader(`{"log": "test"}`)
				case "/projects":
					body = strings.NewReader(`{"name": "test project"}`)
				default:
					body = strings.NewReader("")
				}
			} else {
				body = strings.NewReader("")
			}
			req := httptest.NewRequest(tc.method, tc.path, body)
			req.Header.Set("Content-Type", "application/json")

			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tc.expectedStatus {
				t.Errorf("expected status %d for %s %s, got %d", tc.expectedStatus, tc.method, tc.path, recorder.Code)
			}
		})
	}
}
