package analyzer

import (
	"context"
	"strings"
	"testing"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type mockAIClient struct {
	Resp       string
	Err        error
	PrevPrompt string
}

func (m *mockAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	m.PrevPrompt = prompt
	return m.Resp, m.Err
}

func TestAnalyzer(t *testing.T) {
	t.Run("detect SQLi attack", func(t *testing.T) {
		mockResponse := `{"is_threat": true, "reason": "SQL injection detected", "confidence": 0.95}`

		mockClient := &mockAIClient{Resp: mockResponse}

		analyzer := NewAnalyzer(mockClient)

		logLine := `GET /search?q=' OR 1=1 --`
		req := protocol.AnalysisRequest{
			Log:   logLine,
			IP:    "11.1.2.3",
			Route: "/api/login",
			MetaData: map[string]string{
				"app": "authService",
			},
		}

		res, err := analyzer.Analyze(context.Background(), req)
		assertNoError(t, err)

		if !res.IsThreat {
			t.Fatalf("expected IsThreat=true, got false; resp=%+v", res)
		}

		if strings.TrimSpace(res.Reason) == "" {
			t.Fatalf("expected non-empty Reason")
		}

		if res.Confidence < 0.9 {
			t.Fatalf("expected confidence >= 0.9, got %f", res.Confidence)
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
		req := protocol.AnalysisRequest{
			Log:   logLine,
			IP:    "11.1.2.3",
			Route: "/api/login",
			MetaData: map[string]string{
				"app": "authService",
			},
		}

		_, err := analyzer.Analyze(context.Background(), req)
		assertError(t, err, ErrMalformedAIResponse)
	})

	t.Run("detect ai generate failed", func(t *testing.T) {
		mockClient := &mockAIClient{Err: ErrAIGenerateFailed}
		analyzer := NewAnalyzer(mockClient)

		logLine := `GET /search?q=' OR 1=1 --`
		req := protocol.AnalysisRequest{
			Log:   logLine,
			IP:    "11.1.2.3",
			Route: "/api/login",
			MetaData: map[string]string{
				"app": "authService",
			},
		}

		_, err := analyzer.Analyze(context.Background(), req)
		assertError(t, err, ErrAIGenerateFailed)
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