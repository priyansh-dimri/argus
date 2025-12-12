//go:build integration
package integration

import (
	"context"
	"os"
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

	t.Run("Generate", func(t *testing.T) {
		resp, err := client.Generate(ctx, "Say 'test' and nothing else.")
		if err != nil {
			t.Errorf("Generate failed: %v", err)
		}
		if resp == "" {
			t.Error("Expected non-empty response")
		}
		t.Logf("Gemini Response: %s", resp)
	})
}