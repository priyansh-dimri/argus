package api

import (
	"context"
	"time"

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
	Saved           bool
	ProjectID       string
	Req             protocol.AnalysisRequest
	Res             protocol.AnalysisResponse
	Err             error
	SaveSignal      chan struct{}
	MockProject     *protocol.Project
	MockProjectList []protocol.Project
	MockProjectID   string
}

func (m *mockStore) SaveThreat(ctx context.Context, projectID string, req protocol.AnalysisRequest, res protocol.AnalysisResponse) error {
	m.Saved = true
	m.Req = req
	m.Res = res
	m.ProjectID = projectID

	if m.SaveSignal != nil {
		select {
		case m.SaveSignal <- struct{}{}:
		default:
		}
	}

	return m.Err
}

func (m *mockStore) CreateProject(ctx context.Context, userID string, name string) (*protocol.Project, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if m.MockProject != nil {
		return m.MockProject, nil
	}
	return &protocol.Project{
		ID:        "mock_proj_id",
		UserID:    userID,
		Name:      name,
		APIKey:    "argus_mock_key",
		CreatedAt: time.Now(),
	}, nil
}

func (m *mockStore) GetProjectsByUser(ctx context.Context, userID string) ([]protocol.Project, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.MockProjectList != nil {
		return m.MockProjectList, nil
	}
	return []protocol.Project{}, nil
}

func (m *mockStore) GetProjectIDByKey(ctx context.Context, apiKey string) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	if m.MockProjectID != "" {
		return m.MockProjectID, nil
	}
	return "mock_project_id", nil
}
