package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/priyansh-dimri/argus/pkg/logger"
	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Close()
}

type SupabaseStore struct {
	db       DB
	randRead func(b []byte) (n int, err error)
}

func NewSupabaseStore(db DB) *SupabaseStore {
	logger.Info("Initializing Supabase store", "component", "storage")
	store := &SupabaseStore{
		db:       db,
		randRead: rand.Read,
	}
	logger.Info("Supabase store initialized successfully", "component", "storage")
	return store
}

func (s *SupabaseStore) SaveThreat(ctx context.Context, projectID string, req protocol.AnalysisRequest, res protocol.AnalysisResponse) error {
	start := time.Now()
	logger.Info("Saving threat log to database",
		"component", "storage",
		"operation", "SaveThreat",
		"project_id", projectID,
		"is_threat", res.IsThreat,
	)

	headersJSON, err := json.Marshal(req.Headers)
	if err != nil {
		logger.Error("Failed to marshal headers", err,
			"component", "storage",
			"operation", "SaveThreat",
			"project_id", projectID,
		)
		headersJSON = []byte("{}")
	}

	metadataJSON, err := json.Marshal(req.MetaData)
	if err != nil {
		logger.Error("Failed to marshal metadata", err,
			"component", "storage",
			"operation", "SaveThreat",
			"project_id", projectID,
		)
		metadataJSON = []byte("{}")
	}

	logger.Info("Marshaled request data",
		"component", "storage",
		"headers_size", len(headersJSON),
		"metadata_size", len(metadataJSON),
	)

	const query = `
		INSERT INTO threat_logs (
			project_id,
			ip, 
			route, 
			method,
			headers,
			metadata, 
			payload, 
			is_threat, 
			reason, 
			confidence, 
			timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
	`

	method := "UNKNOWN"

	if m, ok := req.Headers["Method"]; ok {
		method = m
	}

	_, err = s.db.Exec(ctx, query,
		projectID,
		req.IP,
		req.Route,
		method,
		headersJSON,
		metadataJSON,
		req.Log,
		*res.IsThreat,
		*res.Reason,
		*res.Confidence,
	)

	if err != nil {
		logger.Error("Failed to insert threat log", err,
			"component", "storage",
			"operation", "SaveThreat",
			"project_id", projectID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return fmt.Errorf("error while inserting threat log: %w", err)
	}

	logger.Info("Threat log saved successfully",
		"component", "storage",
		"operation", "SaveThreat",
		"project_id", projectID,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return nil
}

func (s *SupabaseStore) CreateProject(ctx context.Context, userID string, name string) (*protocol.Project, error) {
	start := time.Now()
	logger.Info("Creating new project",
		"component", "storage",
		"operation", "CreateProject",
		"user_id", userID,
		"project_name", name,
	)

	apiKey, err := s.generateAPIKey()
	if err != nil {
		logger.Error("Failed to generate API key", err,
			"component", "storage",
			"operation", "CreateProject",
			"user_id", userID,
		)
		return nil, fmt.Errorf("failed to generate api key: %w", err)
	}

	logger.Info("API key generated",
		"component", "storage",
		"key_length", len(apiKey),
	)

	const query = `
		INSERT INTO projects (user_id, name, api_key)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	var projectID string
	var createdAt time.Time

	err = s.db.QueryRow(ctx, query, userID, name, apiKey).Scan(&projectID, &createdAt)
	if err != nil {
		logger.Error("Failed to insert project", err,
			"component", "storage",
			"operation", "CreateProject",
			"user_id", userID,
			"project_name", name,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return nil, fmt.Errorf("failed to insert project: %w", err)
	}

	logger.Info("Project created successfully",
		"component", "storage",
		"operation", "CreateProject",
		"user_id", userID,
		"project_id", projectID,
		"project_name", name,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return &protocol.Project{
		ID:        projectID,
		UserID:    userID,
		Name:      name,
		APIKey:    apiKey,
		CreatedAt: createdAt,
	}, nil
}

func (s *SupabaseStore) GetProjectIDByKey(ctx context.Context, apiKey string) (string, error) {
	start := time.Now()
	logger.Info("Looking up project by API key",
		"component", "storage",
		"operation", "GetProjectIDByKey",
		"key_length", len(apiKey),
	)

	const query = `SELECT id FROM projects WHERE api_key = $1`
	var id string
	err := s.db.QueryRow(ctx, query, apiKey).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Warn("API key not found",
				"component", "storage",
				"operation", "GetProjectIDByKey",
				"duration_ms", time.Since(start).Milliseconds(),
			)
			return "", fmt.Errorf("invalid api key")
		}
		logger.Error("Database error during project lookup", err,
			"component", "storage",
			"operation", "GetProjectIDByKey",
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return "", fmt.Errorf("database error: %w", err)
	}

	logger.Info("Project found by API key",
		"component", "storage",
		"operation", "GetProjectIDByKey",
		"project_id", id,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return id, nil
}

func (s *SupabaseStore) GetProjectsByUser(ctx context.Context, userID string) ([]protocol.Project, error) {
	start := time.Now()
	logger.Info("Fetching projects for user",
		"component", "storage",
		"operation", "GetProjectsByUser",
		"user_id", userID,
	)

	const query = `SELECT id, user_id, name, api_key, created_at FROM projects WHERE user_id = $1`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		logger.Error("Failed to query projects", err,
			"component", "storage",
			"operation", "GetProjectsByUser",
			"user_id", userID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []protocol.Project
	rowCount := 0
	for rows.Next() {
		var p protocol.Project
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.APIKey, &p.CreatedAt); err != nil {
			logger.Error("Failed to scan project row", err,
				"component", "storage",
				"operation", "GetProjectsByUser",
				"user_id", userID,
				"row_number", rowCount,
			)
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}
		projects = append(projects, p)
		rowCount++
	}

	logger.Info("Projects fetched successfully",
		"component", "storage",
		"operation", "GetProjectsByUser",
		"user_id", userID,
		"project_count", len(projects),
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return projects, nil
}

func (s *SupabaseStore) UpdateProjectName(ctx context.Context, userID string, projectID string, newName string) error {
	start := time.Now()
	logger.Info("Updating project name",
		"component", "storage",
		"operation", "UpdateProjectName",
		"user_id", userID,
		"project_id", projectID,
		"new_name", newName,
	)

	const query = `UPDATE projects SET name = $1 WHERE id = $2 AND user_id = $3`
	tag, err := s.db.Exec(ctx, query, newName, projectID, userID)
	if err != nil {
		logger.Error("Failed to update project name", err,
			"component", "storage",
			"operation", "UpdateProjectName",
			"user_id", userID,
			"project_id", projectID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return fmt.Errorf("failed to update project name: %w", err)
	}

	rowsAffected := tag.RowsAffected()
	if rowsAffected == 0 {
		logger.Warn("Project not found for update",
			"component", "storage",
			"operation", "UpdateProjectName",
			"user_id", userID,
			"project_id", projectID,
		)
		return fmt.Errorf("project not found")
	}

	logger.Info("Project name updated successfully",
		"component", "storage",
		"operation", "UpdateProjectName",
		"user_id", userID,
		"project_id", projectID,
		"new_name", newName,
		"rows_affected", rowsAffected,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return nil
}

func (s *SupabaseStore) RotateAPIKey(ctx context.Context, userID string, projectID string) (string, error) {
	start := time.Now()
	logger.Info("Rotating API key",
		"component", "storage",
		"operation", "RotateAPIKey",
		"user_id", userID,
		"project_id", projectID,
	)

	newKey, err := s.generateAPIKey()
	if err != nil {
		logger.Error("Failed to generate new API key", err,
			"component", "storage",
			"operation", "RotateAPIKey",
			"user_id", userID,
			"project_id", projectID,
		)
		return "", fmt.Errorf("failed to generate new key: %w", err)
	}

	logger.Info("New API key generated",
		"component", "storage",
		"key_length", len(newKey),
	)

	const query = `UPDATE projects SET api_key = $1 WHERE id = $2 AND user_id = $3`
	tag, err := s.db.Exec(ctx, query, newKey, projectID, userID)
	if err != nil {
		logger.Error("Failed to update API key", err,
			"component", "storage",
			"operation", "RotateAPIKey",
			"user_id", userID,
			"project_id", projectID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return "", fmt.Errorf("failed to update api key: %w", err)
	}

	rowsAffected := tag.RowsAffected()
	if rowsAffected == 0 {
		logger.Warn("Project not found for key rotation",
			"component", "storage",
			"operation", "RotateAPIKey",
			"user_id", userID,
			"project_id", projectID,
		)
		return "", fmt.Errorf("project not found")
	}

	logger.Info("API key rotated successfully",
		"component", "storage",
		"operation", "RotateAPIKey",
		"user_id", userID,
		"project_id", projectID,
		"rows_affected", rowsAffected,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return newKey, nil
}

func (s *SupabaseStore) DeleteProject(ctx context.Context, userID string, projectID string) error {
	start := time.Now()
	logger.Info("Deleting project",
		"component", "storage",
		"operation", "DeleteProject",
		"user_id", userID,
		"project_id", projectID,
	)

	const query = `DELETE FROM projects WHERE id = $1 AND user_id = $2`
	tag, err := s.db.Exec(ctx, query, projectID, userID)
	if err != nil {
		logger.Error("Failed to delete project", err,
			"component", "storage",
			"operation", "DeleteProject",
			"user_id", userID,
			"project_id", projectID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected := tag.RowsAffected()
	if rowsAffected == 0 {
		logger.Warn("Project not found for deletion",
			"component", "storage",
			"operation", "DeleteProject",
			"user_id", userID,
			"project_id", projectID,
		)
		return fmt.Errorf("project not found")
	}

	logger.Info("Project deleted successfully",
		"component", "storage",
		"operation", "DeleteProject",
		"user_id", userID,
		"project_id", projectID,
		"rows_affected", rowsAffected,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return nil
}

func (s *SupabaseStore) DeleteUser(ctx context.Context, userID string) error {
	start := time.Now()
	logger.Info("Deleting user account",
		"component", "storage",
		"operation", "DeleteUser",
		"user_id", userID,
	)

	const query = `DELETE FROM auth.users WHERE id = $1`
	tag, err := s.db.Exec(ctx, query, userID)
	if err != nil {
		logger.Error("Failed to delete user", err,
			"component", "storage",
			"operation", "DeleteUser",
			"user_id", userID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected := tag.RowsAffected()
	if rowsAffected == 0 {
		logger.Warn("User not found for deletion",
			"component", "storage",
			"operation", "DeleteUser",
			"user_id", userID,
		)
		return fmt.Errorf("user not found")
	}

	logger.Info("User deleted successfully",
		"component", "storage",
		"operation", "DeleteUser",
		"user_id", userID,
		"rows_affected", rowsAffected,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return nil
}

func (s *SupabaseStore) generateAPIKey() (string, error) {
	logger.Info("Generating new API key",
		"component", "storage",
		"operation", "generateAPIKey",
	)

	bytes := make([]byte, 16)
	if _, err := s.randRead(bytes); err != nil {
		logger.Error("Failed to generate random bytes for API key", err,
			"component", "storage",
			"operation", "generateAPIKey",
		)
		return "", err
	}

	apiKey := "argus_" + hex.EncodeToString(bytes)
	logger.Info("API key generated successfully",
		"component", "storage",
		"operation", "generateAPIKey",
		"key_length", len(apiKey),
	)

	return apiKey, nil
}
