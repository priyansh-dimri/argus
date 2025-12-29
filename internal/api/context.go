package api

import (
	"context"

	"github.com/priyansh-dimri/argus/pkg/logger"
)

type contextKey string

const (
	projectIDKey contextKey = "project_id"
	userIDKey    contextKey = "user_id"
)

func WithProjectID(ctx context.Context, id string) context.Context {
	logger.Info("Setting project ID in context",
		"component", "context",
		"project_id", id,
	)
	return context.WithValue(ctx, projectIDKey, id)
}

func GetProjectID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(projectIDKey).(string)
	if !ok || id == "" {
		logger.Warn("Failed to retrieve project ID from context",
			"component", "context",
			"found", ok,
		)
	} else {
		logger.Info("Retrieved project ID from context",
			"component", "context",
			"project_id", id,
		)
	}
	return id, ok
}

func WithUserID(ctx context.Context, id string) context.Context {
	logger.Info("Setting user ID in context",
		"component", "context",
		"user_id", id,
	)
	return context.WithValue(ctx, userIDKey, id)
}

func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	if !ok || id == "" {
		logger.Warn("Failed to retrieve user ID from context",
			"component", "context",
			"found", ok,
		)
	} else {
		logger.Info("Retrieved user ID from context",
			"component", "context",
			"user_id", id,
		)
	}
	return id, ok
}
