package argus

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

func TestClient_SendAnalysis(t *testing.T) {
	reqPayload := protocol.AnalysisRequest{
		Log: "test log",
		IP:  "11.1.2.3",
	}

	t.Run("send analysis successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("Expected POST request, got %s", r.Method)
			}
			if r.URL.Path != "/analyze" {
				t.Errorf("Expected /analyze URL path, got %s", r.URL.Path)
			}

			var received protocol.AnalysisRequest
			if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
				t.Fatalf("Failed to decode request on server: %v", err)
			}
			if received.Log != reqPayload.Log {
				t.Errorf("Expected log %q, got %q", reqPayload.Log, received.Log)
			}

			isThreat := true
			reason := "test threat"
			confidence := 0.9
			json.NewEncoder(w).Encode(protocol.AnalysisResponse{
				IsThreat:   &isThreat,
				Reason:     &reason,
				Confidence: &confidence,
			})
		}))
		defer server.Close()

		client := NewClient(server.URL, "fake-api-key", 2*time.Second)

		resp, err := client.SendAnalysis(reqPayload)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !*resp.IsThreat {
			t.Error("Expected IsThreat=true")
		}
		if *resp.Reason != "test threat" {
			t.Errorf("Expected reason 'test threat', got %q", *resp.Reason)
		}
	})

	t.Run("detect timeout error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewClient(server.URL, "fake-api-key", 10*time.Millisecond)

		_, err := client.SendAnalysis(reqPayload)

		if err == nil {
			t.Fatal("Expected timeout error, got nil")
		}
	})

	t.Run("detect server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient(server.URL, "fake-api-key", 2*time.Second)

		_, err := client.SendAnalysis(reqPayload)

		if err == nil {
			t.Fatal("Expected error for 500 status, got nil")
		}
	})

	t.Run("detect server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient(server.URL, "fake-api-key", 2*time.Second)

		_, err := client.SendAnalysis(reqPayload)

		if err == nil {
			t.Fatal("Expected error for 500 status, got nil")
		}
	})

	t.Run("detect json marshalling error", func(t *testing.T) {
		client := NewClient("http://localhost", "key", time.Second)

		client.marshal = func(v any) ([]byte, error) {
			return nil, errors.New("forced marshal error")
		}

		_, err := client.SendAnalysis(reqPayload)

		if err == nil {
			t.Fatal("Expected marshal error, got nil")
		}
		if err.Error() != "failed to marshal request: forced marshal error" {
			t.Errorf("Unexpected error msg: %v", err)
		}
	})

	t.Run("detect http request creation error", func(t *testing.T) {
		invalidURL := "http://some url with spaces"

		client := NewClient(invalidURL, "key", time.Second)

		_, err := client.SendAnalysis(reqPayload)

		if err == nil {
			t.Fatal("Expected NewRequest error, got nil")
		}

		if !strings.Contains(err.Error(), "failed to create http request") {
			t.Errorf("Expected error to contain 'failed to create http request', got: %v", err)
		}
	})

	t.Run("detects decode error in response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{bad json`))
		}))
		defer server.Close()

		client := NewClient(server.URL, "fake-api-key", 2*time.Second)

		_, err := client.SendAnalysis(reqPayload)

		if err == nil {
			t.Fatal("Expected decode error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to decode response") {
			t.Errorf("Expected error to contain 'failed to decode response', got: %v", err)
		}
	})
}
