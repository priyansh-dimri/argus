package argus

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

var (
	benchmarkMiddlewareResponse *httptest.ResponseRecorder
)

type benchmarkWAF struct {
	blockRequest bool
}

func (w *benchmarkWAF) Check(r *http.Request) (bool, error) {
	return w.blockRequest, nil
}

type benchmarkSender struct {
	response protocol.AnalysisResponse
}

func (s *benchmarkSender) SendAnalysis(req protocol.AnalysisRequest) (protocol.AnalysisResponse, error) {
	return s.response, nil
}

func BenchmarkMiddleware_LatencyFirst_WAFSafe_SmallBody(b *testing.B) {
	waf := &benchmarkWAF{blockRequest: false}
	isThreat := false
	sender := &benchmarkSender{response: protocol.AnalysisResponse{IsThreat: &isThreat}}
	config := Config{Mode: LatencyFirst}
	mw := NewMiddleware(sender, waf, config)

	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	body := strings.NewReader("small body")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/test", body)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		benchmarkMiddlewareResponse = rec
		body.Seek(0, io.SeekStart)
	}
}

func BenchmarkMiddleware_LatencyFirst_WAFBlock(b *testing.B) {
	waf := &benchmarkWAF{blockRequest: true}
	isThreat := true
	sender := &benchmarkSender{response: protocol.AnalysisResponse{IsThreat: &isThreat}}
	config := Config{Mode: LatencyFirst}
	mw := NewMiddleware(sender, waf, config)

	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	body := strings.NewReader("malicious payload")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/attack", body)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		benchmarkMiddlewareResponse = rec
		body.Seek(0, io.SeekStart)
	}
}

func BenchmarkMiddleware_SmartShield_WAFSafe(b *testing.B) {
	waf := &benchmarkWAF{blockRequest: false}
	isThreat := false
	sender := &benchmarkSender{response: protocol.AnalysisResponse{IsThreat: &isThreat}}
	config := Config{Mode: SmartShield}
	mw := NewMiddleware(sender, waf, config)

	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	body := strings.NewReader("clean request")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/data", body)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		benchmarkMiddlewareResponse = rec
		body.Seek(0, io.SeekStart)
	}
}

func BenchmarkMiddleware_SmartShield_WAFThreatAISafe(b *testing.B) {
	waf := &benchmarkWAF{blockRequest: true}
	isThreat := false
	sender := &benchmarkSender{response: protocol.AnalysisResponse{IsThreat: &isThreat}}
	config := Config{Mode: SmartShield}
	mw := NewMiddleware(sender, waf, config)

	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	body := strings.NewReader("suspicious but safe")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/check", body)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		benchmarkMiddlewareResponse = rec
		body.Seek(0, io.SeekStart)
	}
}

func BenchmarkMiddleware_Paranoid_AllRequests(b *testing.B) {
	waf := &benchmarkWAF{blockRequest: false}
	isThreat := false
	sender := &benchmarkSender{response: protocol.AnalysisResponse{IsThreat: &isThreat}}
	config := Config{Mode: Paranoid}
	mw := NewMiddleware(sender, waf, config)

	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	body := strings.NewReader("paranoid check")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/secure", body)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		benchmarkMiddlewareResponse = rec
		body.Seek(0, io.SeekStart)
	}
}

func BenchmarkMiddleware_LargeBody_1KB(b *testing.B) {
	waf := &benchmarkWAF{blockRequest: false}
	isThreat := false
	sender := &benchmarkSender{response: protocol.AnalysisResponse{IsThreat: &isThreat}}
	config := Config{Mode: LatencyFirst}
	mw := NewMiddleware(sender, waf, config)

	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	largeBody := strings.Repeat("x", 1024)
	body := strings.NewReader(largeBody)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/upload", body)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		benchmarkMiddlewareResponse = rec
		body.Seek(0, io.SeekStart)
	}
}

func BenchmarkMiddleware_LargeBody_10KB(b *testing.B) {
	waf := &benchmarkWAF{blockRequest: false}
	isThreat := false
	sender := &benchmarkSender{response: protocol.AnalysisResponse{IsThreat: &isThreat}}
	config := Config{Mode: LatencyFirst}
	mw := NewMiddleware(sender, waf, config)

	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	largeBody := strings.Repeat("data", 2560)
	body := strings.NewReader(largeBody)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/bulk", body)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		benchmarkMiddlewareResponse = rec
		body.Seek(0, io.SeekStart)
	}
}

func BenchmarkMiddleware_ComplexHeaders(b *testing.B) {
	waf := &benchmarkWAF{blockRequest: false}
	isThreat := false
	sender := &benchmarkSender{response: protocol.AnalysisResponse{IsThreat: &isThreat}}
	config := Config{Mode: LatencyFirst}
	mw := NewMiddleware(sender, waf, config)

	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	body := strings.NewReader("test")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/complex", body)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer token123")
		req.Header.Set("X-Request-ID", "req-12345")
		req.Header.Set("X-Forwarded-For", "192.168.1.1")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Connection", "keep-alive")

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		benchmarkMiddlewareResponse = rec
		body.Seek(0, io.SeekStart)
	}
}

func BenchmarkMiddleware_NoBody(b *testing.B) {
	waf := &benchmarkWAF{blockRequest: false}
	isThreat := false
	sender := &benchmarkSender{response: protocol.AnalysisResponse{IsThreat: &isThreat}}
	config := Config{Mode: LatencyFirst}
	mw := NewMiddleware(sender, waf, config)

	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/status", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		benchmarkMiddlewareResponse = rec
	}
}

func BenchmarkMiddleware_BuildPayload_Isolated(b *testing.B) {
	waf := &benchmarkWAF{blockRequest: false}
	sender := &benchmarkSender{}
	config := Config{Mode: LatencyFirst}
	mw := NewMiddleware(sender, waf, config)

	req := httptest.NewRequest("POST", "/api/test", bytes.NewReader([]byte("test body")))
	req.Header.Set("User-Agent", "TestAgent")
	req.Header.Set("Content-Type", "application/json")
	body := []byte("test body")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mw.buildPayload(req, body, false)
	}
}

func BenchmarkMiddleware_BodyRead_Restoration(b *testing.B) {
	testData := []byte(strings.Repeat("test data ", 100))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/test", bytes.NewReader(testData))

		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
}

func BenchmarkMiddleware_Parallel_LatencyFirst(b *testing.B) {
	waf := &benchmarkWAF{blockRequest: false}
	isThreat := false
	sender := &benchmarkSender{response: protocol.AnalysisResponse{IsThreat: &isThreat}}
	config := Config{Mode: LatencyFirst}
	mw := NewMiddleware(sender, waf, config)

	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		body := strings.NewReader("parallel test")
		for pb.Next() {
			req := httptest.NewRequest("POST", "/api/parallel", body)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			body.Seek(0, io.SeekStart)
		}
	})
}

func BenchmarkMiddleware_Parallel_Paranoid(b *testing.B) {
	waf := &benchmarkWAF{blockRequest: false}
	isThreat := false
	sender := &benchmarkSender{response: protocol.AnalysisResponse{IsThreat: &isThreat}}
	config := Config{Mode: Paranoid}
	mw := NewMiddleware(sender, waf, config)

	handler := mw.Protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		body := strings.NewReader("paranoid parallel")
		for pb.Next() {
			req := httptest.NewRequest("POST", "/api/secure", body)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			body.Seek(0, io.SeekStart)
		}
	})
}

func BenchmarkMiddleware_NoWrapper_Baseline(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	body := strings.NewReader("baseline test")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/baseline", body)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		benchmarkMiddlewareResponse = rec
		body.Seek(0, io.SeekStart)
	}
}
