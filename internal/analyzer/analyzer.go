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
	return &Analyzer{
		client:  c,
		marshal: json.Marshal,
	}
}

func (analyzer *Analyzer) Analyze(ctx context.Context, req protocol.AnalysisRequest) (protocol.AnalysisResponse, error) {
	maxTokens := analyzer.client.GetMaxTokens()

	req.Log = TruncateLog(ctx, analyzer.client, req.Log, maxTokens)

	reqJSON, err := analyzer.marshal(req)

	if err != nil {
		logger.Error("Failed to marshal request for AI", err)
		return protocol.AnalysisResponse{}, fmt.Errorf("JSON marshal request error: %v", err)
	}

	prompt := strings.Replace(SecurityAnalysisPrompt, "{{REQUEST_JSON}}", string(reqJSON), 1)

	output, err := analyzer.client.Generate(ctx, prompt)

	if err != nil {
		logger.Error("AI Generation failed", err)
		return protocol.AnalysisResponse{}, ErrAIGenerateFailed
	}

	var response protocol.AnalysisResponse

	if err := json.Unmarshal([]byte(output), &response); err != nil {
		logger.Error("Failed to parse AI response", err, "raw_output", output)
		return protocol.AnalysisResponse{}, ErrMalformedAIResponse
	}

	if response.IsThreat == nil || response.Reason == nil || response.Confidence == nil {
		logger.Warn("AI returned empty reason", "raw_output", output)
		return protocol.AnalysisResponse{}, ErrMalformedAIResponse
	}

	return response, nil
}
