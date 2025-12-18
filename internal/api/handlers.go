package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type Analyzer interface {
	Analyze(ctx context.Context, req protocol.AnalysisRequest) (protocol.AnalysisResponse, error)
}

type Store interface {
	SaveThreat(ctx context.Context, projectID string, req protocol.AnalysisRequest, res protocol.AnalysisResponse) error
	CreateProject(ctx context.Context, userID string, name string) (*protocol.Project, error)
	GetProjectsByUser(ctx context.Context, userID string) ([]protocol.Project, error)
	GetProjectIDByKey(ctx context.Context, apiKey string) (string, error)
	UpdateProjectName(ctx context.Context, userID string, projectID string, newName string) error
	RotateAPIKey(ctx context.Context, userID string, projectID string) (string, error)
	DeleteProject(ctx context.Context, userID string, projectID string) error
}

type API struct {
	Analyzer      Analyzer
	Store         Store
	ErrorReporter func(msg string, args ...any)
}

func NewAPI(analyzer Analyzer, store Store) *API {
	return &API{
		Analyzer:      analyzer,
		Store:         store,
		ErrorReporter: slog.Error,
	}
}

func (api *API) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	projectID, ok := GetProjectID(r.Context())
	if !ok || projectID == "" {
		http.Error(w, "Unauthorized: Missing Project Context", http.StatusUnauthorized)
		return
	}

	var req protocol.AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON decoding error", http.StatusBadRequest)
		return
	}

	res, err := api.Analyzer.Analyze(r.Context(), req)
	if err != nil {
		http.Error(w, "analysis error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(res)

	go func() {
		bgContext := context.Background()
		if err := api.Store.SaveThreat(bgContext, projectID, req, res); err != nil {
			if api.ErrorReporter != nil {
				api.ErrorReporter("Failed to save threat log", "error", err)
			}
		}
	}()
}

func (api *API) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req protocol.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}

	project, err := api.Store.CreateProject(r.Context(), userID, req.Name)
	if err != nil {
		api.ErrorReporter("Failed to create project", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(protocol.CreateProjectResponse{Project: *project})
}

func (api *API) HandleListProjects(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projects, err := api.Store.GetProjectsByUser(r.Context(), userID)
	if err != nil {
		api.ErrorReporter("Failed to list projects", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if projects == nil {
		projects = []protocol.Project{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (api *API) HandleUpdateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req protocol.UpdateProjectNameRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	err := api.Store.UpdateProjectName(r.Context(), userID, req.ID, req.Name)
	if err != nil {
		api.ErrorReporter("Update failed", "error", err)
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *API) HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req protocol.DeleteProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if err := api.Store.DeleteProject(r.Context(), userID, req.ID); err != nil {
		api.ErrorReporter("Delete failed", "error", err)
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *API) HandleRotateKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req protocol.RotateProjectAPIKeyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	newKey, err := api.Store.RotateAPIKey(r.Context(), userID, req.ID)
	if err != nil {
		api.ErrorReporter("Rotation failed", "error", err)
		http.Error(w, "Rotation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(protocol.RotateProjectAPIKeyResponse{APIKey: newKey})
}
