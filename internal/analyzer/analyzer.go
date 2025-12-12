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
	count, err := analyzer.client.CountTokens(ctx, req.Log)
	maxTokens := analyzer.client.GetMaxTokens()

	if err != nil {
		logger.Warn("Failed to count tokens", "error", err)
	} else if count > maxTokens {
		safeLength := len(req.Log) / 2 // TODO: use a better strategy
		if safeLength > 0 {
			req.Log = req.Log[:safeLength] + "...[TRUNCATED]"
			logger.Info("Log truncated due to token limit", "original_tokens", count, "limit", maxTokens)
		}
	}

	reqJSON, err := analyzer.marshal(req)

	if err != nil {
		return protocol.AnalysisResponse{}, fmt.Errorf("JSON marshal request error: %v", err)
	}

	prompt := strings.Replace(SecurityAnalysisPrompt, "{{REQUEST_JSON}}", string(reqJSON), 1)

	output, err := analyzer.client.Generate(ctx, prompt)

	if err != nil {
		return protocol.AnalysisResponse{}, ErrAIGenerateFailed
	}

	var response protocol.AnalysisResponse

	if err := json.Unmarshal([]byte(output), &response); err != nil {
		return protocol.AnalysisResponse{}, ErrMalformedAIResponse
	}

	if response.IsThreat == nil || response.Reason == nil || response.Confidence == nil {
		return protocol.AnalysisResponse{}, ErrMalformedAIResponse
	}

	return response, nil
}
