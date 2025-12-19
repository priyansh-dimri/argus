package analyzer

import (
	"context"
	"strings"
	"testing"
)

var benchmarkTruncatorResult string

type mockTruncatorAIClient struct {
	CountTokensFunc func(ctx context.Context, text string) (int, error)
	GenerateFunc    func(ctx context.Context, prompt string) (string, error)
	MaxTokens       int
}

func (m *mockTruncatorAIClient) CountTokens(ctx context.Context, text string) (int, error) {
	if m.CountTokensFunc != nil {
		return m.CountTokensFunc(ctx, text)
	}
	return len(text) / 4, nil
}

func (m *mockTruncatorAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(ctx, prompt)
	}
	return "", nil
}

func (m *mockTruncatorAIClient) GetMaxTokens() int {
	if m.MaxTokens > 0 {
		return m.MaxTokens
	}
	return 30000
}

func BenchmarkTruncateLog_SmallLog(b *testing.B) {
	mock := &mockTruncatorAIClient{
		CountTokensFunc: func(ctx context.Context, text string) (int, error) {
			return len(text) / 4, nil
		},
	}

	input := strings.Repeat("short log entry ", 50)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result := TruncateLog(ctx, mock, input, 1000)
		benchmarkTruncatorResult = result
	}
}

func BenchmarkTruncateLog_LargeLog_NeedsTruncation(b *testing.B) {
	mock := &mockTruncatorAIClient{
		CountTokensFunc: func(ctx context.Context, text string) (int, error) {
			return len(text) / 4, nil
		},
	}

	input := strings.Repeat("A", 10000)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result := TruncateLog(ctx, mock, input, 1000)
		benchmarkTruncatorResult = result
	}
}

func BenchmarkTruncateLog_VeryLargeLog(b *testing.B) {
	mock := &mockTruncatorAIClient{
		CountTokensFunc: func(ctx context.Context, text string) (int, error) {
			return len(text) / 4, nil
		},
	}

	input := strings.Repeat("LOG ENTRY ", 10000)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result := TruncateLog(ctx, mock, input, 5000)
		benchmarkTruncatorResult = result
	}
}

func BenchmarkTruncateLog_APIFailure_SafeMode(b *testing.B) {
	mock := &mockTruncatorAIClient{
		CountTokensFunc: func(ctx context.Context, text string) (int, error) {
			return 0, context.DeadlineExceeded
		},
	}

	input := strings.Repeat("B", 5000)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result := TruncateLog(ctx, mock, input, 1000)
		benchmarkTruncatorResult = result
	}
}

func BenchmarkTruncateLog_EdgeCase_SingleChar(b *testing.B) {
	mock := &mockTruncatorAIClient{
		CountTokensFunc: func(ctx context.Context, text string) (int, error) {
			return 2, nil
		},
	}

	input := "X"
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result := TruncateLog(ctx, mock, input, 1)
		benchmarkTruncatorResult = result
	}
}

func BenchmarkTruncateLog_ExactLimit(b *testing.B) {
	mock := &mockTruncatorAIClient{
		CountTokensFunc: func(ctx context.Context, text string) (int, error) {
			return 1000, nil
		},
	}

	input := strings.Repeat("A", 4000)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result := TruncateLog(ctx, mock, input, 1000)
		benchmarkTruncatorResult = result
	}
}

func BenchmarkTruncateLog_Parallel(b *testing.B) {
	mock := &mockTruncatorAIClient{
		CountTokensFunc: func(ctx context.Context, text string) (int, error) {
			return len(text) / 4, nil
		},
	}

	input := strings.Repeat("parallel log data ", 500)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result := TruncateLog(ctx, mock, input, 1000)
			benchmarkTruncatorResult = result
		}
	})
}

func BenchmarkTruncateLog_RealWorld_JSONPayload(b *testing.B) {
	mock := &mockTruncatorAIClient{
		CountTokensFunc: func(ctx context.Context, text string) (int, error) {
			return len(text) / 3, nil
		},
	}

	jsonLog := `{"timestamp":"2024-01-15T10:30:00Z","level":"ERROR","request":{"method":"POST","path":"/api/users","body":` +
		strings.Repeat(`{"field":"value","nested":{"key":"data"}},`, 100) +
		`"headers":{"User-Agent":"Mozilla/5.0","Content-Type":"application/json"}},"error":"Internal Server Error"}`

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result := TruncateLog(ctx, mock, jsonLog, 2000)
		benchmarkTruncatorResult = result
	}
}

func BenchmarkTruncateLog_NoTruncationNeeded(b *testing.B) {
	mock := &mockTruncatorAIClient{
		CountTokensFunc: func(ctx context.Context, text string) (int, error) {
			return 50, nil
		},
	}

	input := "Small log that doesn't need truncation"
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result := TruncateLog(ctx, mock, input, 1000)
		benchmarkTruncatorResult = result
	}
}
