package analyzer

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestTruncator(t *testing.T) {
	t.Run("does not shorten small logs", func(t *testing.T) {
		mock := &mockAIClient{
			CountTokensFunc: func(ctx context.Context, text string) (int, error) {
				return 100, nil
			},
		}
		input := "short log"
		output := TruncateLog(context.Background(), mock, input, 1000)
		if output != input {
			t.Errorf("Expected same log, got truncated: %s", output)
		}
	})

	t.Run("does not shorten logs with token count within limit", func(t *testing.T) {
		mock := &mockAIClient{
			CountTokensFunc: func(ctx context.Context, text string) (int, error) {
				return 999, nil
			},
		}
		input := strings.Repeat("A", 1005)
		output := TruncateLog(context.Background(), mock, input, 1000)
		if output != input {
			t.Errorf("Expected same log, got truncated: %s", output)
		}
	})

	t.Run("truncates big logs using ratio", func(t *testing.T) {
		input := strings.Repeat("A", 10000)

		mock := &mockAIClient{
			CountTokensFunc: func(ctx context.Context, text string) (int, error) {
				return len(text), nil
			},
		}

		output := TruncateLog(context.Background(), mock, input, 5000)

		if len(output) >= len(input) {
			t.Errorf("Log was not truncated. Got length: %d", len(output))
		}

		if !strings.Contains(output, "[TRUNCATED]") {
			t.Error("Missing [TRUNCATED] marker")
		}

		if len(output) > 5000 {
			t.Errorf("Truncation insufficient. Got length: %d", len(output))
		}
	})

	t.Run("falls back to hard limit on API failure", func(t *testing.T) {
		mock := &mockAIClient{
			CountTokensFunc: func(ctx context.Context, text string) (int, error) {
				return 0, errors.New("api down")
			},
		}

		input := strings.Repeat("B", 1000)

		output := TruncateLog(context.Background(), mock, input, 100)

		if !strings.Contains(output, "[TRUNCATED_SAFE_MODE]") {
			t.Error("Missing [TRUNCATED_SAFE_MODE] marker")
		}

		if len(output) > 350 {
			t.Errorf("Safe mode insufficient. Got Length: %d", len(output))
		}
	})

	t.Run("return same small log even on API failure", func(t *testing.T) {
		mock := &mockAIClient{
			CountTokensFunc: func(ctx context.Context, text string) (int, error) {
				return 0, errors.New("api failure")
			},
		}

		input := strings.Repeat("A", 150)

		output := TruncateLog(context.Background(), mock, input, 100)

		// Should NOT contain TRUNCATED_SAFE_MODE because it fits
		if output != input {
			t.Error("Expected log to be returned as-is when API fails but log is small")
		}
	})

	t.Run("truncate even if length is one if token count exceeds", func(t *testing.T) {
		mock := &mockAIClient{
			CountTokensFunc: func(ctx context.Context, text string) (int, error) {
				return 2, nil
			},
		}
		input := "A"
		output := TruncateLog(context.Background(), mock, input, 1)

		if !strings.Contains(output, "[TRUNCATED]") {
			t.Errorf("Expected marker even for 1 char string, got %q", output)
		}
	})
}
