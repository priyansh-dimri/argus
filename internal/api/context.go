package api

import "context"

type contextKey string

const (
	projectIDKey contextKey = "project_id"
	userIDKey    contextKey = "user_id"
)

func WithProjectID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, projectIDKey, id)
}

func GetProjectID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(projectIDKey).(string)
	return id, ok
}

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
