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
		{
			name:           "Valid + Authenticated PATCH /projects",
			method:         http.MethodPatch,
			path:           "/projects",
			authHeader:     "Bearer " + validJWT,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid + Unauthenticated PATCH /projects",
			method:         http.MethodPatch,
			path:           "/projects",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Method PUT /projects (for update route)",
			method:         http.MethodPut,
			path:           "/projects",
			authHeader:     "Bearer " + validJWT,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Valid + Authenticated DELETE /projects",
			method:         http.MethodDelete,
			path:           "/projects",
			authHeader:     "Bearer " + validJWT,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid + Unauthenticated DELETE /projects",
			method:         http.MethodDelete,
			path:           "/projects",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Method PATCH /projects as delete (wrong method)",
			method:         http.MethodPatch,
			path:           "/projects",
			authHeader:     "Bearer " + validJWT,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid + Authenticated POST /rotate-key",
			method:         http.MethodPost,
			path:           "/rotate-key",
			authHeader:     "Bearer " + validJWT,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid + Unauthenticated POST /rotate-key",
			method:         http.MethodPost,
			path:           "/rotate-key",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Method GET /rotate-key",
			method:         http.MethodGet,
			path:           "/rotate-key",
			authHeader:     "Bearer " + validJWT,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body *strings.Reader

			switch {
			case tc.method == http.MethodPost && tc.path == "/analyze":
				body = strings.NewReader(`{"log": "test"}`)
			case tc.method == http.MethodPost && tc.path == "/projects":
				body = strings.NewReader(`{"name": "test project"}`)
			case tc.method == http.MethodPatch && tc.path == "/projects":
				body = strings.NewReader(`{"id": "proj_1", "name": "updated"}`)
			case tc.method == http.MethodDelete && tc.path == "/projects":
				body = strings.NewReader(`{"id": "proj_1"}`)
			case tc.method == http.MethodPost && tc.path == "/rotate-key":
				body = strings.NewReader(`{"id": "proj_1"}`)
			default:
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
