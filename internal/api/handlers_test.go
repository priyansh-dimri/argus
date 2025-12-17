package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

func TestAnalyzeHandler(t *testing.T) {
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

	t.Run("return unauthorized when project context is missing", func(t *testing.T) {
		api := &API{Analyzer: nil, Store: nil}

		req := httptest.NewRequest(http.MethodPost, "/analyze", nil)
		recorder := httptest.NewRecorder()

		api.HandleAnalyze(recorder, req)

		assertStatusCode(t, recorder.Code, http.StatusUnauthorized)
	})
}

func TestHandleCreateProject(t *testing.T) {
	t.Run("return unauthorized when user_id is missing", func(t *testing.T) {
		api := &API{Store: &mockStore{}}
		req := httptest.NewRequest(http.MethodPost, "/projects", nil)
		recorder := httptest.NewRecorder()

		api.HandleCreateProject(recorder, req)

		assertStatusCode(t, recorder.Code, http.StatusUnauthorized)
	})

	t.Run("return error for invalid JSON body", func(t *testing.T) {
		api := &API{Store: &mockStore{}}
		req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader("{bad json"))
		req = addUserIDContext(req)
		recorder := httptest.NewRecorder()

		api.HandleCreateProject(recorder, req)

		assertStatusCode(t, recorder.Code, http.StatusBadRequest)
	})

	t.Run("return error when project name is empty", func(t *testing.T) {
		api := &API{Store: &mockStore{}}
		body := strings.NewReader(`{"name": ""}`)
		req := httptest.NewRequest(http.MethodPost, "/projects", body)
		req = addUserIDContext(req)
		recorder := httptest.NewRecorder()

		api.HandleCreateProject(recorder, req)

		assertStatusCode(t, recorder.Code, http.StatusBadRequest)
	})

	t.Run("return internal error when storage fails", func(t *testing.T) {
		store := &mockStore{Err: errors.New("db failure")}
		errorChan := make(chan error, 1)
		api := &API{
			Store: store,
			ErrorReporter: func(msg string, args ...any) {
				if len(args) > 1 {
					if err, ok := args[1].(error); ok {
						errorChan <- err
					}
				}
			},
		}

		body := strings.NewReader(`{"name": "My Project"}`)
		req := httptest.NewRequest(http.MethodPost, "/projects", body)
		req = addUserIDContext(req)
		recorder := httptest.NewRecorder()

		api.HandleCreateProject(recorder, req)

		assertStatusCode(t, recorder.Code, http.StatusInternalServerError)

		select {
		case err := <-errorChan:
			if err.Error() != "db failure" {
				t.Errorf("expected 'db failure', got %v", err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error report")
		}
	})

	t.Run("create project successfully", func(t *testing.T) {
		store := &mockStore{}
		api := &API{Store: store}

		body := strings.NewReader(`{"name": "New Project"}`)
		req := httptest.NewRequest(http.MethodPost, "/projects", body)
		req = addUserIDContext(req)
		recorder := httptest.NewRecorder()

		api.HandleCreateProject(recorder, req)

		assertStatusCode(t, recorder.Code, http.StatusOK)

		var resp protocol.CreateProjectResponse
		if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Project.Name != "New Project" {
			t.Errorf("expected project name 'New Project', got %q", resp.Project.Name)
		}
	})
}

func TestHandleListProjects(t *testing.T) {
	t.Run("return unauthorized when user_id is missing", func(t *testing.T) {
		api := &API{Store: &mockStore{}}
		// Request without addUserIDContext
		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		recorder := httptest.NewRecorder()

		api.HandleListProjects(recorder, req)

		assertStatusCode(t, recorder.Code, http.StatusUnauthorized)
	})

	t.Run("return internal error when storage fails", func(t *testing.T) {
		store := &mockStore{Err: errors.New("db failure")}
		errorChan := make(chan error, 1)
		api := &API{
			Store: store,
			ErrorReporter: func(msg string, args ...any) {
				if len(args) > 1 {
					if err, ok := args[1].(error); ok {
						errorChan <- err
					}
				}
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		req = addUserIDContext(req)
		recorder := httptest.NewRecorder()

		api.HandleListProjects(recorder, req)

		assertStatusCode(t, recorder.Code, http.StatusInternalServerError)

		select {
		case err := <-errorChan:
			if err.Error() != "db failure" {
				t.Errorf("expected 'db failure', got %v", err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error report")
		}
	})

	t.Run("return empty list when no projects found", func(t *testing.T) {
		// Mock returns nil list by default if MockProjectList is nil
		store := &mockStore{}
		api := &API{Store: store}

		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		req = addUserIDContext(req)
		recorder := httptest.NewRecorder()

		api.HandleListProjects(recorder, req)

		assertStatusCode(t, recorder.Code, http.StatusOK)

		var projects []protocol.Project
		if err := json.NewDecoder(recorder.Body).Decode(&projects); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if projects == nil || len(projects) != 0 {
			t.Errorf("expected empty list, got %v", projects)
		}
	})

	t.Run("return list of projects successfully", func(t *testing.T) {
		expectedProjects := []protocol.Project{
			{ID: "p1", Name: "Project A"},
			{ID: "p2", Name: "Project B"},
		}
		store := &mockStore{
			MockProjectList: expectedProjects,
		}
		api := &API{Store: store}

		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		req = addUserIDContext(req)
		recorder := httptest.NewRecorder()

		api.HandleListProjects(recorder, req)

		assertStatusCode(t, recorder.Code, http.StatusOK)

		var projects []protocol.Project
		if err := json.NewDecoder(recorder.Body).Decode(&projects); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(projects) != 2 {
			t.Errorf("expected 2 projects, got %d", len(projects))
		}
		if projects[0].ID != "p1" || projects[1].Name != "Project B" {
			t.Errorf("unexpected project data: %v", projects)
		}
	})

	t.Run("return empty list when DB explicitly returns nil", func(t *testing.T) {
		store := &nilProjectStore{mockStore: &mockStore{}}
		api := &API{Store: store}

		req := httptest.NewRequest(http.MethodGet, "/projects", nil)
		req = addUserIDContext(req)
		recorder := httptest.NewRecorder()

		api.HandleListProjects(recorder, req)

		assertStatusCode(t, recorder.Code, http.StatusOK)

		var projects []protocol.Project
		if err := json.NewDecoder(recorder.Body).Decode(&projects); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		// Ensure the API converted nil to empty slice []
		if projects == nil {
			t.Error("expected empty slice [], got nil")
		}
		if len(projects) != 0 {
			t.Errorf("expected empty list, got len %d", len(projects))
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
	req = addAuthContext(req)
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

func addAuthContext(req *http.Request) *http.Request {
	ctx := WithProjectID(req.Context(), "test-project-id")
	return req.WithContext(ctx)
}

func addUserIDContext(req *http.Request) *http.Request {
	ctx := WithUserID(req.Context(), "test-user-id")
	return req.WithContext(ctx)
}
