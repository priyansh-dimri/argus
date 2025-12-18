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
	return &SupabaseStore{
		db:       db,
		randRead: rand.Read,
	}
}

func (s *SupabaseStore) SaveThreat(ctx context.Context, projectID string, req protocol.AnalysisRequest, res protocol.AnalysisResponse) error {
	headersJSON, _ := json.Marshal(req.Headers)
	metadataJSON, _ := json.Marshal(req.MetaData)

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

	_, err := s.db.Exec(ctx, query,
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
		return fmt.Errorf("error while inserting threat log: %w", err)
	}

	return nil
}

func (s *SupabaseStore) CreateProject(ctx context.Context, userID string, name string) (*protocol.Project, error) {
	apiKey, err := s.generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate api key: %w", err)
	}

	const query = `
		INSERT INTO projects (user_id, name, api_key)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	var projectID string
	var createdAt time.Time

	err = s.db.QueryRow(ctx, query, userID, name, apiKey).Scan(&projectID, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert project: %w", err)
	}

	return &protocol.Project{
		ID:        projectID,
		UserID:    userID,
		Name:      name,
		APIKey:    apiKey,
		CreatedAt: createdAt,
	}, nil
}

func (s *SupabaseStore) GetProjectIDByKey(ctx context.Context, apiKey string) (string, error) {
	const query = `SELECT id FROM projects WHERE api_key = $1`
	var id string
	err := s.db.QueryRow(ctx, query, apiKey).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("invalid api key")
		}
		return "", fmt.Errorf("database error: %w", err)
	}
	return id, nil
}

func (s *SupabaseStore) GetProjectsByUser(ctx context.Context, userID string) ([]protocol.Project, error) {
	const query = `SELECT id, user_id, name, api_key, created_at FROM projects WHERE user_id = $1`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []protocol.Project
	for rows.Next() {
		var p protocol.Project
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.APIKey, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (s *SupabaseStore) generateAPIKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := s.randRead(bytes); err != nil {
		return "", err
	}
	return "argus_" + hex.EncodeToString(bytes), nil
}
