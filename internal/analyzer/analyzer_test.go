package analyzer

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type mockAIClient struct {
	Resp             string
	Err              error
	PrevPrompt       string
	CountTokensFunc  func(ctx context.Context, text string) (int, error)
	GetMaxTokensFunc func() int
}

func (m *mockAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	m.PrevPrompt = prompt
	return m.Resp, m.Err
}

func (m *mockAIClient) CountTokens(ctx context.Context, text string) (int, error) {
	if m.CountTokensFunc != nil {
		return m.CountTokensFunc(ctx, text)
	}
	return 100, nil
}

func (m *mockAIClient) GetMaxTokens() int {
	if m.GetMaxTokensFunc != nil {
		return m.GetMaxTokensFunc()
	}
	return 4000
}

func TestAnalyzer(t *testing.T) {
	t.Run("detect SQLi attack", func(t *testing.T) {
		mockResponse := `{"is_threat": true, "reason": "SQL injection detected", "confidence": 0.95}`

		mockClient := &mockAIClient{Resp: mockResponse}

		analyzer := NewAnalyzer(mockClient)

		logLine := `GET /search?q=' OR 1=1 --`
		req := newTestRequest(logLine)

		res, err := analyzer.Analyze(context.Background(), req)
		assertNoError(t, err)

		if !(*res.IsThreat) {
			t.Fatalf("expected IsThreat=true, got false; resp=%+v", res)
		}

		if strings.TrimSpace(*res.Reason) == "" {
			t.Fatalf("expected non-empty Reason")
		}

		if *res.Confidence < 0.9 {
			t.Fatalf("expected confidence >= 0.9, got %f", *res.Confidence)
		}

		if !strings.Contains(mockClient.PrevPrompt, logLine) {
			t.Fatalf("expected prompt to include log content; prompt=%q", mockClient.PrevPrompt)
		}
	})

	t.Run("detect malformed response", func(t *testing.T) {
		mockResponse := `{"is_threat": true, "reason": "SQL injection detected", "confidence": 0.95`

		mockClient := &mockAIClient{Resp: mockResponse}

		analyzer := NewAnalyzer(mockClient)

		logLine := `GET /search?q=' OR 1=1 --`
		req := newTestRequest(logLine)

		_, err := analyzer.Analyze(context.Background(), req)
		assertError(t, err, ErrMalformedAIResponse)
	})

	t.Run("detect ai generate failed", func(t *testing.T) {
		mockClient := &mockAIClient{Err: ErrAIGenerateFailed}
		analyzer := NewAnalyzer(mockClient)

		logLine := `GET /search?q=' OR 1=1 --`
		req := newTestRequest(logLine)

		_, err := analyzer.Analyze(context.Background(), req)
		assertError(t, err, ErrAIGenerateFailed)
	})

	t.Run("return json marshal error", func(t *testing.T) {
		mockClient := &mockAIClient{}
		analyzer := &Analyzer{
			client: mockClient,
			marshal: func(v any) ([]byte, error) {
				return nil, fmt.Errorf("forced marshal error")
			},
		}

		logLine := `GET /search?q=' OR 1=1 --`
		req := newTestRequest(logLine)

		_, err := analyzer.Analyze(context.Background(), req)
		if err == nil || !strings.Contains(err.Error(), "JSON marshal request error") {
			t.Fatalf("expected JSON marshal request error, got %v", err)
		}
	})

	t.Run("detect empty AI response", func(t *testing.T) {
		mockResponse := ""
		mockClient := &mockAIClient{Resp: mockResponse}
		analyzer := NewAnalyzer(mockClient)

		logLine := `GET /foo/bar`
		req := newTestRequest(logLine)

		res, err := analyzer.Analyze(context.Background(), req)
		assertError(t, err, ErrMalformedAIResponse)

		if (res != protocol.AnalysisResponse{}) {
			t.Errorf("expected empty AnalysisResponse, got %+v", res)
		}
	})

	t.Run("detect missing isThreat field in response", func(t *testing.T) {
		mockResponse := `{"reason": "SQL injection detected", "confidence": 0.95}`
		mockClient := &mockAIClient{Resp: mockResponse}
		analyzer := NewAnalyzer(mockClient)

		logLine := `GET /search?q=' OR 1=1 --`
		req := newTestRequest(logLine)

		_, err := analyzer.Analyze(context.Background(), req)
		assertError(t, err, ErrMalformedAIResponse)
	})
	t.Run("detect missing reason field in response", func(t *testing.T) {
		mockResponse := `{"is_threat": true, "confidence": 0.95}`
		mockClient := &mockAIClient{Resp: mockResponse}
		analyzer := NewAnalyzer(mockClient)

		logLine := `GET /search?q=' OR 1=1 --`
		req := newTestRequest(logLine)

		_, err := analyzer.Analyze(context.Background(), req)
		assertError(t, err, ErrMalformedAIResponse)
	})
	t.Run("detect missing confidence field in response", func(t *testing.T) {
		mockResponse := `{"is_threat": true, "reason": "SQL injection detected"}`
		mockClient := &mockAIClient{Resp: mockResponse}
		analyzer := NewAnalyzer(mockClient)

		logLine := `GET /search?q=' OR 1=1 --`
		req := newTestRequest(logLine)

		_, err := analyzer.Analyze(context.Background(), req)
		assertError(t, err, ErrMalformedAIResponse)
	})

	t.Run("detect multiple missing fields in response", func(t *testing.T) {
		mockResponse := `{"confidence": 0.95}`
		mockClient := &mockAIClient{Resp: mockResponse}
		analyzer := NewAnalyzer(mockClient)

		logLine := `GET /search?q=' OR 1=1 --`
		req := newTestRequest(logLine)

		_, err := analyzer.Analyze(context.Background(), req)
		assertError(t, err, ErrMalformedAIResponse)
	})

	t.Run("truncate logs exceeding token count limit", func(t *testing.T) {
		mockClient := &mockAIClient{
			Resp: `{"is_threat": false, "reason": "safe", "confidence": 1.0}`,
			CountTokensFunc: func(ctx context.Context, text string) (int, error) {
				return 5000, nil
			},
			GetMaxTokensFunc: func() int {
				return 4000
			},
		}

		analyzer := NewAnalyzer(mockClient)

		hugeLog := strings.Repeat("A", 10000) // 10 thousand 'A's
		request := newTestRequest(hugeLog)

		_, err := analyzer.Analyze(context.Background(), request)
		assertNoError(t, err)

		if strings.Contains(mockClient.PrevPrompt, hugeLog) {
			t.Fatal("expected truncated log, found original log")
		}

		if !strings.Contains(mockClient.PrevPrompt, "[TRUNCATED]") {
			t.Fatal("Prompt does not contain [TRUNCATED]")
		}
	})

	t.Run("proceed even if CountTokens fails", func(t *testing.T) {
		mockClient := &mockAIClient{
			Resp: `{"is_threat": false, "reason": "safe", "confidence": 1.0}`,
			CountTokensFunc: func(ctx context.Context, text string) (int, error) {
				return 0, fmt.Errorf("some error")
			},
		}

		analyzer := NewAnalyzer(mockClient)
		req := newTestRequest("log")

		_, err := analyzer.Analyze(context.Background(), req)
		assertNoError(t, err)

		if mockClient.PrevPrompt == "" {
			t.Fatal("expected Generate to be called even if CountTokens fails")
		}
	})
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Fatalf("got error %q, wanted %q", got, want)
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()

	if got != nil {
		t.Fatalf("got error %q when not wanted", got)
	}
}

func newTestRequest(logLine string) protocol.AnalysisRequest {
	return protocol.AnalysisRequest{
		Log:   logLine,
		IP:    "11.1.2.3",
		Route: "/api/login",
		MetaData: map[string]string{
			"app": "authService",
		},
	}
}
