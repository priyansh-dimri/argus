package argus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type AnalysisSender interface {
	SendAnalysis(req protocol.AnalysisRequest) (protocol.AnalysisResponse, error)
}

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	marshal    func(v any) ([]byte, error)
}

var _ AnalysisSender = (*Client)(nil) // compile time check

func NewClient(baseURL, apiKey string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey: apiKey,
		marshal: json.Marshal,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) SendAnalysis(req protocol.AnalysisRequest) (protocol.AnalysisResponse, error) {
	bodyBytes, err := c.marshal(req)
	if err != nil {
		return protocol.AnalysisResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := fmt.Sprintf("%s/analyze", c.baseURL)
	httpReq, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return protocol.AnalysisResponse{}, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return protocol.AnalysisResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return protocol.AnalysisResponse{}, fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	var analysisResp protocol.AnalysisResponse
	if err := json.NewDecoder(resp.Body).Decode(&analysisResp); err != nil {
		return protocol.AnalysisResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return analysisResp, nil
}
