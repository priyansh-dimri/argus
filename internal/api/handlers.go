package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/priyansh-dimri/argus/pkg/logger"
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
	DeleteUser(ctx context.Context, userID string) error
}

type API struct {
	Analyzer      Analyzer
	Store         Store
	ErrorReporter func(msg string, err error, args ...any)
}

func NewAPI(analyzer Analyzer, store Store) *API {
	logger.Info("Initializing API handlers", "component", "api")
	api := &API{
		Analyzer:      analyzer,
		Store:         store,
		ErrorReporter: logger.Error,
	}
	logger.Info("API handlers initialized successfully", "component", "api")
	return api
}

func (api *API) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger.Info("HandleAnalyze started",
		"component", "handler",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	projectID, ok := GetProjectID(r.Context())
	if !ok || projectID == "" {
		logger.Warn("Analyze request missing project context",
			"component", "handler",
			"remote_addr", r.RemoteAddr,
		)
		http.Error(w, "Unauthorized: Missing Project Context", http.StatusUnauthorized)
		return
	}

	logger.Info("Project context retrieved",
		"component", "handler",
		"project_id", projectID,
	)

	var req protocol.AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode analysis request JSON", err,
			"component", "handler",
			"project_id", projectID,
			"remote_addr", r.RemoteAddr,
		)
		http.Error(w, "JSON decoding error", http.StatusBadRequest)
		return
	}

	logger.Info("Analysis request decoded successfully",
		"component", "handler",
		"project_id", projectID,
		"log_length", len(req.Log),
	)

	aiCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	logger.Info("Starting AI analysis",
		"component", "handler",
		"project_id", projectID,
		"timeout_seconds", 60,
	)

	res, err := api.Analyzer.Analyze(aiCtx, req)
	if err != nil {
		logger.Error("Analysis failed", err,
			"component", "handler",
			"project_id", projectID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		http.Error(w, "analysis error", http.StatusInternalServerError)
		return
	}

	logger.Info("Analysis completed successfully",
		"component", "handler",
		"project_id", projectID,
		"is_threat", res.IsThreat,
		"confidence", res.Confidence,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Error("Failed to encode analysis response", err,
			"component", "handler",
			"project_id", projectID,
		)
	} else {
		logger.Info("Analysis response sent successfully",
			"component", "handler",
			"project_id", projectID,
		)
	}

	go func() {
		saveStart := time.Now()
		logger.Info("Starting background threat save",
			"component", "handler",
			"project_id", projectID,
		)

		bgContext := context.Background()
		if err := api.Store.SaveThreat(bgContext, projectID, req, res); err != nil {
			if api.ErrorReporter != nil {
				api.ErrorReporter("Failed to save threat log", err,
					"component", "handler",
					"project_id", projectID,
					"duration_ms", time.Since(saveStart).Milliseconds(),
				)
			}
		} else {
			logger.Info("Threat log saved successfully",
				"component", "handler",
				"project_id", projectID,
				"duration_ms", time.Since(saveStart).Milliseconds(),
			)
		}
	}()
}

func (api *API) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger.Info("HandleCreateProject started",
		"component", "handler",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	userID, ok := GetUserID(r.Context())
	if !ok || userID == "" {
		logger.Warn("Create project request missing user context",
			"component", "handler",
			"remote_addr", r.RemoteAddr,
		)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req protocol.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode create project request", err,
			"component", "handler",
			"user_id", userID,
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		logger.Warn("Create project request missing project name",
			"component", "handler",
			"user_id", userID,
		)
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}

	logger.Info("Creating new project",
		"component", "handler",
		"user_id", userID,
		"project_name", req.Name,
	)

	project, err := api.Store.CreateProject(r.Context(), userID, req.Name)
	if err != nil {
		api.ErrorReporter("Failed to create project", err,
			"component", "handler",
			"user_id", userID,
			"project_name", req.Name,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("Project created successfully",
		"component", "handler",
		"user_id", userID,
		"project_id", project.ID,
		"project_name", project.Name,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(protocol.CreateProjectResponse{Project: *project})
}

func (api *API) HandleListProjects(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger.Info("HandleListProjects started",
		"component", "handler",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	userID, ok := GetUserID(r.Context())
	if !ok || userID == "" {
		logger.Warn("List projects request missing user context",
			"component", "handler",
			"remote_addr", r.RemoteAddr,
		)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	logger.Info("Fetching projects for user",
		"component", "handler",
		"user_id", userID,
	)

	projects, err := api.Store.GetProjectsByUser(r.Context(), userID)
	if err != nil {
		api.ErrorReporter("Failed to list projects", err,
			"component", "handler",
			"user_id", userID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if projects == nil {
		projects = []protocol.Project{}
		logger.Info("No projects found, returning empty array",
			"component", "handler",
			"user_id", userID,
		)
	}

	logger.Info("Projects retrieved successfully",
		"component", "handler",
		"user_id", userID,
		"project_count", len(projects),
		"duration_ms", time.Since(start).Milliseconds(),
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (api *API) HandleUpdateProject(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger.Info("HandleUpdateProject started",
		"component", "handler",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	userID, ok := GetUserID(r.Context())
	if !ok || userID == "" {
		logger.Warn("Update project request missing user context",
			"component", "handler",
			"remote_addr", r.RemoteAddr,
		)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req protocol.UpdateProjectNameRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode update project request", err,
			"component", "handler",
			"user_id", userID,
		)
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	logger.Info("Updating project name",
		"component", "handler",
		"user_id", userID,
		"project_id", req.ID,
		"new_name", req.Name,
	)

	err := api.Store.UpdateProjectName(r.Context(), userID, req.ID, req.Name)
	if err != nil {
		api.ErrorReporter("Update failed", err,
			"component", "handler",
			"user_id", userID,
			"project_id", req.ID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	logger.Info("Project updated successfully",
		"component", "handler",
		"user_id", userID,
		"project_id", req.ID,
		"new_name", req.Name,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	w.WriteHeader(http.StatusOK)
}

func (api *API) HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger.Info("HandleDeleteProject started",
		"component", "handler",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	userID, ok := GetUserID(r.Context())
	if !ok || userID == "" {
		logger.Warn("Delete project request missing user context",
			"component", "handler",
			"remote_addr", r.RemoteAddr,
		)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req protocol.DeleteProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode delete project request", err,
			"component", "handler",
			"user_id", userID,
		)
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	logger.Info("Deleting project",
		"component", "handler",
		"user_id", userID,
		"project_id", req.ID,
	)

	if err := api.Store.DeleteProject(r.Context(), userID, req.ID); err != nil {
		api.ErrorReporter("Delete failed", err,
			"component", "handler",
			"user_id", userID,
			"project_id", req.ID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	logger.Info("Project deleted successfully",
		"component", "handler",
		"user_id", userID,
		"project_id", req.ID,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	w.WriteHeader(http.StatusOK)
}

func (api *API) HandleRotateKey(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger.Info("HandleRotateKey started",
		"component", "handler",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	userID, ok := GetUserID(r.Context())
	if !ok || userID == "" {
		logger.Warn("Rotate key request missing user context",
			"component", "handler",
			"remote_addr", r.RemoteAddr,
		)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req protocol.RotateProjectAPIKeyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode rotate key request", err,
			"component", "handler",
			"user_id", userID,
		)
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	logger.Info("Rotating API key",
		"component", "handler",
		"user_id", userID,
		"project_id", req.ID,
	)

	newKey, err := api.Store.RotateAPIKey(r.Context(), userID, req.ID)
	if err != nil {
		api.ErrorReporter("Rotation failed", err,
			"component", "handler",
			"user_id", userID,
			"project_id", req.ID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		http.Error(w, "Rotation failed", http.StatusInternalServerError)
		return
	}

	logger.Info("API key rotated successfully",
		"component", "handler",
		"user_id", userID,
		"project_id", req.ID,
		"new_key_length", len(newKey),
		"duration_ms", time.Since(start).Milliseconds(),
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(protocol.RotateProjectAPIKeyResponse{APIKey: newKey})
}

func (api *API) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger.Info("HandleDeleteAccount started",
		"component", "handler",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	userID, ok := GetUserID(r.Context())
	if !ok || userID == "" {
		logger.Warn("Delete account request missing user context",
			"component", "handler",
			"remote_addr", r.RemoteAddr,
		)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	logger.Info("Deleting user account",
		"component", "handler",
		"user_id", userID,
	)

	if err := api.Store.DeleteUser(r.Context(), userID); err != nil {
		api.ErrorReporter("Failed to delete account", err,
			"component", "handler",
			"user_id", userID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	logger.Info("User account deleted successfully",
		"component", "handler",
		"user_id", userID,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	w.WriteHeader(http.StatusOK)
}
