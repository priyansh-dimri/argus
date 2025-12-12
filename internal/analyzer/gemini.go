package analyzer

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	model  string
}

func NewGeminiClient(ctx context.Context, apiKey string, modelName string) (*GeminiClient, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	return &GeminiClient{
		client: client,
		model:  modelName,
	}, nil
}

func (g *GeminiClient) Generate(ctx context.Context, prompt string) (string, error) {
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	}

	response, err := g.client.Models.GenerateContent(ctx, g.model, genai.Text(prompt), config)

	if err != nil {
		return "", fmt.Errorf("gemini response generation failed: %w", err)
	}

	if len(response.Candidates) == 0 || response.Candidates[0].Content == nil || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned empty response")
	}

	var sb string
	for _, part := range response.Candidates[0].Content.Parts {
		sb += part.Text
	}

	if sb == "" {
		return "", fmt.Errorf("gemini response contained no text")
	}

	return sb, nil
}

func (g *GeminiClient) CountTokens(ctx context.Context, text string) (int, error) {
	resp, err := g.client.Models.CountTokens(ctx, g.model, genai.Text(text), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to count tokens: %w", err)
	}
	return int(resp.TotalTokens), nil
}

func (g *GeminiClient) GetMaxTokens() int {
	return 30_000
}
