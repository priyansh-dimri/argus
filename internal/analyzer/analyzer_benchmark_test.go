package analyzer

import (
	"context"
	"strings"
	"testing"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

var benchmarkAnalyzerResult protocol.AnalysisResponse

type mockBenchmarkAIClient struct {
	response   string
	tokenCount int
	maxTokens  int
}

func (m *mockBenchmarkAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	return m.response, nil
}

func (m *mockBenchmarkAIClient) CountTokens(ctx context.Context, text string) (int, error) {
	if m.tokenCount > 0 {
		return m.tokenCount, nil
	}
	return len(text) / 4, nil
}

func (m *mockBenchmarkAIClient) GetMaxTokens() int {
	if m.maxTokens > 0 {
		return m.maxTokens
	}
	return 30000
}

func BenchmarkAnalyzer_SmallLog(b *testing.B) {
	mock := &mockBenchmarkAIClient{
		response:   `{"is_threat": false, "reason": "No malicious patterns detected", "confidence": 0.95}`,
		tokenCount: 50,
	}

	analyzer := NewAnalyzer(mock)

	req := protocol.AnalysisRequest{
		Log:   `GET /api/users?id=12345`,
		IP:    "192.168.1.1",
		Route: "/api/users",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0",
		},
		MetaData: map[string]string{
			"app": "webService",
		},
	}

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := analyzer.Analyze(ctx, req)
		if err != nil {
			b.Fatalf("Analyze failed: %v", err)
		}
		benchmarkAnalyzerResult = result
	}
}

func BenchmarkAnalyzer_ThreatDetection(b *testing.B) {
	mock := &mockBenchmarkAIClient{
		response:   `{"is_threat": true, "reason": "SQL injection attempt detected in query parameter", "confidence": 0.98}`,
		tokenCount: 100,
	}

	analyzer := NewAnalyzer(mock)

	req := protocol.AnalysisRequest{
		Log:   `GET /search?q=' OR 1=1 --`,
		IP:    "10.0.0.5",
		Route: "/search",
		Headers: map[string]string{
			"User-Agent": "sqlmap/1.0",
		},
		MetaData: map[string]string{
			"waf_result": "BLOCK",
		},
	}

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := analyzer.Analyze(ctx, req)
		if err != nil {
			b.Fatalf("Analyze failed: %v", err)
		}
		benchmarkAnalyzerResult = result
	}
}

func BenchmarkAnalyzer_LargePayload(b *testing.B) {
	mock := &mockBenchmarkAIClient{
		response:   `{"is_threat": false, "reason": "Large payload appears safe", "confidence": 0.85}`,
		tokenCount: 2000,
	}

	analyzer := NewAnalyzer(mock)

	largeBody := `{"data":[` + strings.Repeat(`{"field":"value","nested":{"key":"data"}},`, 100) + `]}`

	req := protocol.AnalysisRequest{
		Log:   largeBody,
		IP:    "172.16.0.10",
		Route: "/api/bulk",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		MetaData: map[string]string{
			"app": "bulkService",
		},
	}

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := analyzer.Analyze(ctx, req)
		if err != nil {
			b.Fatalf("Analyze failed: %v", err)
		}
		benchmarkAnalyzerResult = result
	}
}

func BenchmarkAnalyzer_WithTruncation(b *testing.B) {
	mock := &mockBenchmarkAIClient{
		response:   `{"is_threat": false, "reason": "Truncated log analyzed", "confidence": 0.80}`,
		tokenCount: 35000,
		maxTokens:  30000,
	}

	analyzer := NewAnalyzer(mock)

	hugeLog := strings.Repeat("A", 50000)

	req := protocol.AnalysisRequest{
		Log:   hugeLog,
		IP:    "10.1.1.1",
		Route: "/api/upload",
		MetaData: map[string]string{
			"truncated": "yes",
		},
	}

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := analyzer.Analyze(ctx, req)
		if err != nil {
			b.Fatalf("Analyze failed: %v", err)
		}
		benchmarkAnalyzerResult = result
	}
}

func BenchmarkAnalyzer_ComplexHeaders(b *testing.B) {
	mock := &mockBenchmarkAIClient{
		response:   `{"is_threat": false, "reason": "Headers appear legitimate", "confidence": 0.92}`,
		tokenCount: 150,
	}

	analyzer := NewAnalyzer(mock)

	headers := make(map[string]string)
	for i := range 20 {
		headers["X-Custom-Header-"+string(rune(i+65))] = "value-" + string(rune(i+97))
	}

	req := protocol.AnalysisRequest{
		Log:     `POST /api/data`,
		IP:      "192.168.100.50",
		Route:   "/api/data",
		Headers: headers,
		MetaData: map[string]string{
			"app": "complexService",
		},
	}

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := analyzer.Analyze(ctx, req)
		if err != nil {
			b.Fatalf("Analyze failed: %v", err)
		}
		benchmarkAnalyzerResult = result
	}
}

func BenchmarkAnalyzer_JSONMarshaling(b *testing.B) {
	mock := &mockBenchmarkAIClient{
		response:   `{"is_threat": false, "reason": "Safe", "confidence": 0.99}`,
		tokenCount: 100,
	}

	analyzer := NewAnalyzer(mock)

	req := protocol.AnalysisRequest{
		Log:   strings.Repeat("x", 1000),
		IP:    "1.2.3.4",
		Route: "/test",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   "TestAgent/1.0",
		},
		MetaData: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := analyzer.Analyze(ctx, req)
		if err != nil {
			b.Fatalf("Analyze failed: %v", err)
		}
		benchmarkAnalyzerResult = result
	}
}

func BenchmarkAnalyzer_Parallel(b *testing.B) {
	mock := &mockBenchmarkAIClient{
		response:   `{"is_threat": false, "reason": "Parallel test", "confidence": 0.90}`,
		tokenCount: 75,
	}

	analyzer := NewAnalyzer(mock)

	req := protocol.AnalysisRequest{
		Log:   `GET /api/test`,
		IP:    "192.168.1.100",
		Route: "/api/test",
		Headers: map[string]string{
			"User-Agent": "BenchmarkClient/1.0",
		},
		MetaData: map[string]string{
			"test": "parallel",
		},
	}

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result, err := analyzer.Analyze(ctx, req)
			if err != nil {
				b.Fatalf("Analyze failed: %v", err)
			}
			benchmarkAnalyzerResult = result
		}
	})
}

func BenchmarkAnalyzer_MinimalRequest(b *testing.B) {
	mock := &mockBenchmarkAIClient{
		response:   `{"is_threat": false, "reason": "OK", "confidence": 1.0}`,
		tokenCount: 10,
	}

	analyzer := NewAnalyzer(mock)

	req := protocol.AnalysisRequest{
		Log:      "GET /",
		IP:       "127.0.0.1",
		Route:    "/",
		MetaData: map[string]string{},
	}

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := analyzer.Analyze(ctx, req)
		if err != nil {
			b.Fatalf("Analyze failed: %v", err)
		}
		benchmarkAnalyzerResult = result
	}
}

func BenchmarkAnalyzer_RealWorldMix(b *testing.B) {
	mock := &mockBenchmarkAIClient{
		response:   `{"is_threat": false, "reason": "Normal API request", "confidence": 0.93}`,
		tokenCount: 200,
	}

	analyzer := NewAnalyzer(mock)

	req := protocol.AnalysisRequest{
		Log:   `POST /api/v1/users {"name":"John Doe","email":"john@example.com","role":"user"}`,
		IP:    "203.0.113.45",
		Route: "/api/v1/users",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		},
		MetaData: map[string]string{
			"waf_result": "PASS",
			"app":        "userService",
			"version":    "v1.2.3",
		},
	}

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := analyzer.Analyze(ctx, req)
		if err != nil {
			b.Fatalf("Analyze failed: %v", err)
		}
		benchmarkAnalyzerResult = result
	}
}
