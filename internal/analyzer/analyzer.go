package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/priyansh-dimri/argus/pkg/logger"
	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type AIClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
	CountTokens(ctx context.Context, text string) (int, error)
	GetMaxTokens() int
}

type Analyzer struct {
	client  AIClient
	marshal func(v any) ([]byte, error)
}

func NewAnalyzer(c AIClient) *Analyzer {
	logger.Info("Initializing new Analyzer instance", "client_type", fmt.Sprintf("%T", c))
	return &Analyzer{
		client:  c,
		marshal: json.Marshal,
	}
}

func (analyzer *Analyzer) Analyze(ctx context.Context, req protocol.AnalysisRequest) (protocol.AnalysisResponse, error) {
	logger.Info("Starting security analysis",
		"component", "analyzer",
		"log_length", len(req.Log),
	)

	maxTokens := analyzer.client.GetMaxTokens()
	logger.Info("Retrieved max tokens from AI client",
		"component", "analyzer",
		"max_tokens", maxTokens,
	)

	originalLogLen := len(req.Log)
	req.Log = TruncateLog(ctx, analyzer.client, req.Log, maxTokens)

	logger.Info("Log truncation done",
		"component", "analyzer",
		"original_length", originalLogLen,
		"truncated_length", len(req.Log),
	)

	reqJSON, err := analyzer.marshal(req)

	if err != nil {
		logger.Error("Failed to marshal analysis request to JSON", err,
			"component", "analyzer",
		)
		return protocol.AnalysisResponse{}, fmt.Errorf("JSON marshal request error: %v", err)
	}

	logger.Info("Successfully marshaled request",
		"component", "analyzer",
		"json_size", len(reqJSON),
	)

	prompt := strings.Replace(SecurityAnalysisPrompt, "{{REQUEST_JSON}}", string(reqJSON), 1)
	logger.Info("Generated AI prompt",
		"component", "analyzer",
		"prompt_length", len(prompt),
	)

	output, err := analyzer.client.Generate(ctx, prompt)

	if err != nil {
		logger.Error("AI generation failed", err,
			"component", "analyzer",
			"prompt_length", len(prompt),
		)
		return protocol.AnalysisResponse{}, ErrAIGenerateFailed
	}

	logger.Info("Received AI response",
		"component", "analyzer",
		"response_length", len(output),
	)

	var response protocol.AnalysisResponse

	if err := json.Unmarshal([]byte(output), &response); err != nil {
		logger.Error("Failed to parse AI response JSON", err,
			"component", "analyzer",
			"raw_output", output,
			"output_length", len(output),
		)
		return protocol.AnalysisResponse{}, ErrMalformedAIResponse
	}

	if response.IsThreat == nil || response.Reason == nil || response.Confidence == nil {
		logger.Warn("AI returned incomplete response with missing fields",
			"component", "analyzer",
			"has_is_threat", response.IsThreat != nil,
			"has_reason", response.Reason != nil,
			"has_confidence", response.Confidence != nil,
			"raw_output", output,
		)
		return protocol.AnalysisResponse{}, ErrMalformedAIResponse
	}

	logger.Info("Analysis completed successfully",
		"component", "analyzer",
		"is_threat", *response.IsThreat,
		"confidence", *response.Confidence,
	)

	return response, nil
}
