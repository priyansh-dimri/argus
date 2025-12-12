package integration

import (
	"context"
	"os"
	"strings"

	// "strings"
	"testing"

	"github.com/priyansh-dimri/argus/internal/analyzer"
)

func TestGeminiIntegration(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: GEMINI_API_KEY not set")
	}

	ctx := context.Background()
	client, err := analyzer.NewGeminiClient(ctx, apiKey, "gemini-2.5-flash")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("CountTokens", func(t *testing.T) {
		count, err := client.CountTokens(ctx, "Hello Argus")
		if err != nil {
			t.Errorf("CountTokens failed: %v", err)
		}
		if count <= 0 {
			t.Errorf("Expected tokens > 0, got %d", count)
		}
		t.Logf("Token count for 'Hello Argus' is %d", count)
	})

	t.Run("Generate response successfully", func(t *testing.T) {
		if max := client.GetMaxTokens(); max != 30_000 {
			t.Errorf("GetMaxTokens() = %d; want 30000", max)
		}

		count, err := client.CountTokens(ctx, "Hello Argus")
		if err != nil {
			t.Errorf("CountTokens failed: %v", err)
		}

		if count <= 0 {
			t.Errorf("Expected tokens > 0, got %d", count)
		}

		resp, err := client.Generate(ctx, "Say 'test' and nothing else.")
		if err != nil {
			t.Errorf("Generate failed: %v", err)
		}
		if resp == "" {
			t.Error("Expected non-empty response")
		}
		t.Logf("Gemini Response: %s", resp)
	})

	t.Run("detect error in CountTokens due to fake model", func(t *testing.T) {
		client, _ := analyzer.NewGeminiClient(ctx, apiKey, "fake-model")
		_, err := client.CountTokens(ctx, "Hello")
		if err == nil {
			t.Error("Expected count tokens error with bad model, got nil")
		} else if !strings.Contains(err.Error(), "failed to count tokens") {
			t.Errorf("Expected count error wrapper, got: %v", err)
		}
	})
}
