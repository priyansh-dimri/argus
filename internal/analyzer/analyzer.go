package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type AIClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type Analyzer struct {
	client AIClient
}

func NewAnalyzer(c AIClient) *Analyzer {
	return &Analyzer{client: c}
}

func (analyzer *Analyzer) Analyze(ctx context.Context, req protocol.AnalysisRequest) (protocol.AnalysisResponse, error) {
	reqJSON, err := json.Marshal(req)

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

	return response, nil
}
