package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

func TestRouter(t *testing.T) {
	mock := newMockAnalyzer(protocol.AnalysisResponse{}, nil)
	store := &mockStore{}
	api := NewAPI(mock, store)

	router := NewRouter(api)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Valid POST /analyze",
			method:         http.MethodPost,
			path:           "/analyze",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Method GET /analyze",
			method:         http.MethodGet,
			path:           "/analyze",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Unknown Route POST /random",
			method:         http.MethodPost,
			path:           "/random",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body := strings.NewReader(`{"log": "test"}`)
			req := httptest.NewRequest(tc.method, tc.path, body)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tc.expectedStatus {
				t.Errorf("expected status %d for %s %s, got %d", tc.expectedStatus, tc.method, tc.path, recorder.Code)
			}
		})
	}
}
