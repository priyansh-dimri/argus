package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()
}

type SupabaseStore struct {
	db DB
}

func NewSupabaseStore(db DB) *SupabaseStore {
	return &SupabaseStore{
		db: db,
	}
}

func (s *SupabaseStore) SaveThreat(ctx context.Context, req protocol.AnalysisRequest, res protocol.AnalysisResponse) error {
	headersJSON, _ := json.Marshal(req.Headers)
	metadataJSON, _ := json.Marshal(req.MetaData)

	const query = `
		INSERT INTO threat_logs (
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
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
	`

	method := "UNKNOWN"

	if m, ok := req.Headers["Method"]; ok {
		method = m
	}

	_, err := s.db.Exec(ctx, query,
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
