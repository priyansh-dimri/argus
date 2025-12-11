package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type mockAnalyzer struct {
	Response    protocol.AnalysisResponse
	Err         error
	PrevRequest protocol.AnalysisRequest
}

func newMockAnalyzer(response protocol.AnalysisResponse, err error) *mockAnalyzer {
	return &mockAnalyzer{Response: response, Err: err}
}

func (m *mockAnalyzer) Analyze(ctx context.Context, req protocol.AnalysisRequest) (protocol.AnalysisResponse, error) {
	m.PrevRequest = req
	return m.Response, m.Err
}

func TestHandlers(t *testing.T) {
	t.Run("return threat response", func(t *testing.T) {
		isThreat := true
		reason := "SQLi attack"
		confidence := 0.99

		response := protocol.AnalysisResponse{
			IsThreat:   &isThreat,
			Reason:     &reason,
			Confidence: &confidence,
		}

		mock := newMockAnalyzer(response, nil)
		api := &API{Analyzer: mock, Store: nil}

		request_body := map[string]string{"log": `GET /search?q=' OR 1=1 --`}
		body, err := json.Marshal(request_body)

		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		api.HandleAnalyze(recorder, req)

		resp := recorder.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 OK, got %d", resp.StatusCode)
		}

		var got protocol.AnalysisResponse
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		if got.IsThreat == nil || !*got.IsThreat {
			t.Fatalf("expected is_threat=true in response; got %+v", got)
		}
	})

	t.Run("return error for non-POST request", func(t *testing.T) {
		response := protocol.AnalysisResponse{}
		mock := newMockAnalyzer(response, nil)
		api := &API{Analyzer: mock, Store: nil}

		request_body := map[string]string{"log": `GET /search?q=' OR 1=1 --`}
		body, err := json.Marshal(request_body)

		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/analyze", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		api.HandleAnalyze(recorder, req)

		resp := recorder.Result()
		defer resp.Body.Close()

		assertStatusCode(t, resp.StatusCode, http.StatusMethodNotAllowed)
	})

	t.Run("return error for invalid JSON", func(t *testing.T) {
		response := protocol.AnalysisResponse{}
		mock := newMockAnalyzer(response, nil)
		api := &API{Analyzer: mock, Store: nil}

		body := strings.NewReader(`{bad json}`)
		req := httptest.NewRequest(http.MethodPost, "/analyze", body)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		api.HandleAnalyze(recorder, req)

		resp := recorder.Result()
		defer resp.Body.Close()

		assertStatusCode(t, resp.StatusCode, http.StatusBadRequest)
	})

	t.Run("return error for analysis error", func(t *testing.T) {
		response := protocol.AnalysisResponse{}
		mock := newMockAnalyzer(response, errors.New("analyzer failure"))
		api := &API{Analyzer: mock, Store: nil}

		request_body := map[string]string{"log": "test"}
		body, err := json.Marshal(request_body)

		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		api.HandleAnalyze(recorder, req)

		resp := recorder.Result()
		defer resp.Body.Close()

		assertStatusCode(t, resp.StatusCode, http.StatusInternalServerError)

		if mock.PrevRequest.Log != "test" {
			t.Fatalf("analyzer was not called with expected log %q; got %q", request_body["log"], mock.PrevRequest.Log)
		}
	})
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %d, wanted %d", got, want)
	}
}
