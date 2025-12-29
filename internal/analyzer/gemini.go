package analyzer

import (
	"context"
	"fmt"

	"github.com/priyansh-dimri/argus/pkg/logger"
	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	model  string
}

func NewGeminiClient(ctx context.Context, apiKey string, modelName string) (*GeminiClient, error) {
	logger.Info("Initializing Gemini AI client.",
		"component", "gemini",
		"model", modelName,
	)

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		logger.Error("Failed to create Gemini client", err,
			"component", "gemini",
			"model", modelName,
		)
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	logger.Info("Gemini client initialized successfully",
		"component", "gemini",
		"model", modelName,
	)

	return &GeminiClient{
		client: client,
		model:  modelName,
	}, nil
}

func (g *GeminiClient) Generate(ctx context.Context, prompt string) (string, error) {
	logger.Info("Starting content generation request",
		"component", "gemini",
		"model", g.model,
		"prompt_length", len(prompt),
	)

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	}

	response, err := g.client.Models.GenerateContent(ctx, g.model, genai.Text(prompt), config)

	if err != nil {
		logger.Error("Gemini API content generation failed", err,
			"component", "gemini",
			"model", g.model,
			"prompt_length", len(prompt),
		)
		return "", fmt.Errorf("gemini response generation failed: %w", err)
	}
	logger.Info("Received response from Gemini API",
		"component", "gemini",
		"candidates_count", len(response.Candidates),
	)

	if len(response.Candidates) == 0 || response.Candidates[0].Content == nil || len(response.Candidates[0].Content.Parts) == 0 {
		logger.Error("Gemini returned empty or malformed response", fmt.Errorf("empty response"),
			"component", "gemini",
			"model", g.model,
			"candidates_count", len(response.Candidates),
		)
		return "", fmt.Errorf("gemini returned empty response")
	}

	var sb string
	for i, part := range response.Candidates[0].Content.Parts {
		logger.Info("Processing response part",
			"component", "gemini",
			"part_index", i,
			"part_length", len(part.Text),
		)
		sb += part.Text
	}

	if sb == "" {
		logger.Error("Gemini response contained no text content", fmt.Errorf("empty text"),
			"component", "gemini",
			"model", g.model,
			"parts_count", len(response.Candidates[0].Content.Parts),
		)
		return "", fmt.Errorf("gemini response contained no text")
	}

	logger.Info("Content generation completed successfully",
		"component", "gemini",
		"response_length", len(sb),
		"parts_processed", len(response.Candidates[0].Content.Parts),
	)

	return sb, nil
}

func (g *GeminiClient) CountTokens(ctx context.Context, text string) (int, error) {
	logger.Info("Starting token count request",
		"component", "gemini",
		"model", g.model,
		"text_length", len(text),
	)

	resp, err := g.client.Models.CountTokens(ctx, g.model, genai.Text(text), nil)
	if err != nil {
		logger.Error("Failed to count tokens", err,
			"component", "gemini",
			"model", g.model,
			"text_length", len(text),
		)
		return 0, fmt.Errorf("failed to count tokens: %w", err)
	}

	tokenCount := int(resp.TotalTokens)
	logger.Info("Token count completed",
		"component", "gemini",
		"text_length", len(text),
		"token_count", tokenCount,
	)

	return tokenCount, nil
}

func (g *GeminiClient) GetMaxTokens() int {
	maxTokens := 30_000
	logger.Info("Retrieved max tokens configuration",
		"component", "gemini",
		"max_tokens", maxTokens,
	)
	return maxTokens
}
