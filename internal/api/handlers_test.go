package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

type mockStore struct {
	Saved      bool
	Req        protocol.AnalysisRequest
	Res        protocol.AnalysisResponse
	Err        error
	SaveSignal chan struct{}
}

func (m *mockStore) SaveThreat(ctx context.Context, req protocol.AnalysisRequest, res protocol.AnalysisResponse) error {
	m.Saved = true
	m.Req = req
	m.Res = res

	select {
	case m.SaveSignal <- struct{}{}:
	default:
	}

	return m.Err
}

func TestHandlers(t *testing.T) {
	t.Run("return threat response and SaveThreat asynchronously", func(t *testing.T) {
		response := sampleThreat()

		mock := newMockAnalyzer(response, nil)

		saveChan := make(chan struct{}, 1)
		store := &mockStore{SaveSignal: saveChan}
		api := &API{Analyzer: mock, Store: store}

		request_body := map[string]string{"log": `GET /search?q=' OR 1=1 --`}
		req, recorder := newJSONRequest(t, http.MethodPost, "/analyze", request_body)
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

		select {
		case <-saveChan:
			if !store.Saved {
				t.Error("SaveSignal is received but store.Saved is false")
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for SaveThreat call")
		}
	})

	t.Run("return threat response and return SaveThreat error asynchronously", func(t *testing.T) {
		response := sampleThreat()
		mock := newMockAnalyzer(response, nil)
		store := &mockStore{Err: errors.New("database connection lost")}

		errorChan := make(chan error, 1)
		api := &API{
			Analyzer: mock,
			Store:    store,
			ErrorReporter: func(msg string, args ...any) {
				if len(args) > 1 {
					if err, ok := args[1].(error); ok {
						errorChan <- err
					}
				}
			},
		}

		request_body := map[string]string{"log": `GET /search?q=' OR 1=1 --`}
		req, recorder := newJSONRequest(t, http.MethodPost, "/analyze", request_body)
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

		select {
		case err := <-errorChan:
			if err.Error() != "database connection lost" {
				t.Errorf("expected 'database connection lost', got %v", err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for async SaveThreat error")
		}
	})

	t.Run("return error for non-POST request", func(t *testing.T) {
		response := protocol.AnalysisResponse{}
		mock := newMockAnalyzer(response, nil)
		api := &API{Analyzer: mock, Store: nil}

		request_body := map[string]string{"log": `GET /search?q=' OR 1=1 --`}
		req, recorder := newJSONRequest(t, http.MethodGet, "/analyze", request_body)
		api.HandleAnalyze(recorder, req)

		resp := recorder.Result()
		defer resp.Body.Close()

		assertStatusCode(t, resp.StatusCode, http.StatusMethodNotAllowed)
	})

	t.Run("return error for invalid JSON", func(t *testing.T) {
		response := protocol.AnalysisResponse{}
		mock := newMockAnalyzer(response, nil)
		api := &API{Analyzer: mock, Store: nil}

		req, recorder := newJSONRequest(t, http.MethodPost, "/analyze", `{bad json}`)
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
		req, recorder := newJSONRequest(t, http.MethodPost, "/analyze", request_body)
		api.HandleAnalyze(recorder, req)

		resp := recorder.Result()
		defer resp.Body.Close()

		assertStatusCode(t, resp.StatusCode, http.StatusInternalServerError)

		if mock.PrevRequest.Log != "test" {
			t.Fatalf("analyzer was not called with expected log %q; got %q", request_body["log"], mock.PrevRequest.Log)
		}
	})
}

func TestNewAPI(t *testing.T) {
	mock := newMockAnalyzer(protocol.AnalysisResponse{}, nil)
	store := &mockStore{}

	api := NewAPI(mock, store)

	if api == nil {
		t.Fatal("NewAPI returned nil")
	}

	if api.Analyzer != mock {
		t.Error("NewAPI did not assign Analyzer")
	}

	if api.Store != store {
		t.Error("NewAPI did not assign Store")
	}

	if api.ErrorReporter == nil {
		t.Error("NewAPI did not assign a default ErrorReporter")
	}
}

func newJSONRequest(t testing.TB, method, path string, request_body any) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()
	var buf io.Reader
	if request_body != nil {
		body, err := json.Marshal(request_body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		buf = bytes.NewReader(body)
	}

	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %d, wanted %d", got, want)
	}
}

func sampleThreat() protocol.AnalysisResponse {
	isThreat := true
	reason := "SQLi attack"
	confidence := 0.99
	return protocol.AnalysisResponse{
		IsThreat:   &isThreat,
		Reason:     &reason,
		Confidence: &confidence,
	}
}
