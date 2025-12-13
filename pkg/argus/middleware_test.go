package argus

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type MockWAF struct {
	BlockRequest bool
	Err          error
}

func (m *MockWAF) Check(r *http.Request) (bool, error) {
	return m.BlockRequest, m.Err
}

type MockSender struct {
	Response   protocol.AnalysisResponse
	Err        error
	SentReq    protocol.AnalysisRequest
	CallSignal chan struct{}
}

func (m *MockSender) SendAnalysis(req protocol.AnalysisRequest) (protocol.AnalysisResponse, error) {
	m.SentReq = req
	if m.CallSignal != nil {
		select {
		case m.CallSignal <- struct{}{}:
		default:
		}
	}
	return m.Response, m.Err
}

func TestLatencyFirstMiddleware(t *testing.T) {
	t.Run("WAF Safe -> Pass immediately + Async Log", func(t *testing.T) {
		waf := &MockWAF{BlockRequest: false}
		sender := &MockSender{CallSignal: make(chan struct{}, 1)}
		config := Config{Mode: LatencyFirst}
		mw := NewMiddleware(sender, waf, config)

		handlerCalled := false
		handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("POST", "/api", strings.NewReader("clean"))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if !handlerCalled {
			t.Error("LatencyFirst: Clean request should reach handler")
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", rec.Code)
		}

		select {
		case <-sender.CallSignal:
		case <-time.After(50 * time.Millisecond):
			t.Error("LatencyFirst: Failed to send async log")
		}
	})

	t.Run("WAF Threat -> Block immediately + Async Log", func(t *testing.T) {
		waf := &MockWAF{BlockRequest: true}
		sender := &MockSender{CallSignal: make(chan struct{}, 1)}
		config := Config{Mode: LatencyFirst}
		mw := NewMiddleware(sender, waf, config)

		handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Handler should NOT be called")
		}))

		req := httptest.NewRequest("POST", "/api", strings.NewReader("attack"))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected 403, got %d", rec.Code)
		}

		select {
		case <-sender.CallSignal:
		case <-time.After(50 * time.Millisecond):
			t.Error("LatencyFirst: Failed to send async log on block")
		}
	})
}

func TestSmartShieldMiddleware(t *testing.T) {
	t.Run("WAF Safe -> Pass immediately + Async Log", func(t *testing.T) {
		waf := &MockWAF{BlockRequest: false}
		sender := &MockSender{CallSignal: make(chan struct{}, 1)}
		config := Config{Mode: SmartShield}
		mw := NewMiddleware(sender, waf, config)

		handlerCalled := false
		handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		}))

		req := httptest.NewRequest("POST", "/api", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if !handlerCalled {
			t.Error("SmartShield: Clean WAF request should reach handler")
		}

		select {
		case <-sender.CallSignal:
		case <-time.After(50 * time.Millisecond):
			t.Error("SmartShield: Failed to send async log")
		}
	})

	t.Run("WAF Threat + AI Safe -> Pass", func(t *testing.T) {
		waf := &MockWAF{BlockRequest: true}
		isThreat := false
		sender := &MockSender{
			Response:   protocol.AnalysisResponse{IsThreat: &isThreat},
			CallSignal: make(chan struct{}, 1),
		}
		config := Config{Mode: SmartShield}
		mw := NewMiddleware(sender, waf, config)

		handlerCalled := false
		handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		}))

		req := httptest.NewRequest("POST", "/api", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if !handlerCalled {
			t.Error("SmartShield: AI Safe verdict should overrule WAF Threat")
		}

		select {
		case <-sender.CallSignal:
		default:
			t.Error("SmartShield: AI should be called synchronously")
		}
	})

	t.Run("WAF Threat + AI Threat -> Block", func(t *testing.T) {
		waf := &MockWAF{BlockRequest: true}
		isThreat := true
		sender := &MockSender{
			Response: protocol.AnalysisResponse{IsThreat: &isThreat},
		}
		config := Config{Mode: SmartShield}
		mw := NewMiddleware(sender, waf, config)

		handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Handler should NOT be called")
		}))

		req := httptest.NewRequest("POST", "/api", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected 403, got %d", rec.Code)
		}
	})
}

func TestParanoidMiddleware(t *testing.T) {
	t.Run("WAF Result Ignored + AI Safe -> Pass", func(t *testing.T) {
		waf := &MockWAF{BlockRequest: true}
		isThreat := false

		sender := &MockSender{
			Response: protocol.AnalysisResponse{IsThreat: &isThreat},
		}
		config := Config{Mode: Paranoid}
		mw := NewMiddleware(sender, waf, config)

		handlerCalled := false
		handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		}))

		req := httptest.NewRequest("POST", "/api", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if !handlerCalled {
			t.Error("Paranoid: AI Safe verdict should allow request")
		}

		if sender.SentReq.MetaData["waf_result"] != "BLOCK" {
			t.Errorf("Paranoid: Expected WAF result in metadata, got %v", sender.SentReq.MetaData)
		}
	})

	t.Run("AI Threat -> Block", func(t *testing.T) {
		waf := &MockWAF{BlockRequest: false}
		isThreat := true

		sender := &MockSender{
			Response: protocol.AnalysisResponse{IsThreat: &isThreat},
		}
		config := Config{Mode: Paranoid}
		mw := NewMiddleware(sender, waf, config)

		handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Handler should NOT be called")
		}))

		req := httptest.NewRequest("POST", "/api", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected 403, got %d", rec.Code)
		}
	})
}

func TestHeaderParsingInMiddleware(t *testing.T) {
	waf := &MockWAF{BlockRequest: false}
	sender := &MockSender{CallSignal: make(chan struct{}, 1)}
	config := Config{Mode: LatencyFirst}
	mw := NewMiddleware(sender, waf, config)

	req := httptest.NewRequest("POST", "/test", nil)

	req.Header["User-Agent"] = []string{"Go-Test-Client"}
	req.Header["Empty-Slice-Header"] = []string{}
	req.Header["Empty-String-Header"] = []string{""}

	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	select {
	case <-sender.CallSignal:
		headers := sender.SentReq.Headers

		if val, ok := headers["User-Agent"]; !ok || val != "Go-Test-Client" {
			t.Errorf("Expected User-Agent header to be captured, got: %s", val)
		}

		if _, ok := headers["Empty-Slice-Header"]; ok {
			t.Error("Expected Empty-Slice-Header to be skipped (len=0), but it was present")
		}

		if val, ok := headers["Empty-String-Header"]; !ok || val != "" {
			t.Errorf("Expected Empty-String-Header to be captured as empty string, got: %q", val)
		}

	case <-time.After(50 * time.Millisecond):
		t.Fatal("Timeout waiting for async log")
	}
}
