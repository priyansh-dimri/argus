package api

import (
	"context"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type mockAnalyzer struct {
	Response    protocol.AnalysisResponse
	Err         error
	PrevRequest protocol.AnalysisRequest
}

func (m *mockAnalyzer) Analyze(ctx context.Context, req protocol.AnalysisRequest) (protocol.AnalysisResponse, error) {
	m.PrevRequest = req
	return m.Response, m.Err
}

func newMockAnalyzer(response protocol.AnalysisResponse, err error) *mockAnalyzer {
	return &mockAnalyzer{Response: response, Err: err}
}

type mockStore struct {
	Saved      bool
	Req        protocol.AnalysisRequest
	Res        protocol.AnalysisResponse
	Err        error
	SaveSignal chan struct{}
}

func (m *mockStore) SaveThreat(ctx context.Context, req protocol.AnalysisRequest, res protocol.AnalysisResponse) error {
	m.Saved = true
	m.Req = req
	m.Res = res

	select {
	case m.SaveSignal <- struct{}{}:
	default:
	}

	return m.Err
}
