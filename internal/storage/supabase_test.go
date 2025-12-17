package storage

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/priyansh-dimri/argus/pkg/protocol"
)

func TestSupabase(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error when not expected: %s", err)
	}
	defer mock.Close()

	store := NewSupabaseStore(mock)

	request := protocol.AnalysisRequest{
		Log: "admin' --",
		IP:  "11.1.2.3",
		Headers: map[string]string{
			"Method": "POST",
		},
		Route: "/api/login",
		MetaData: map[string]string{
			"app": "authService",
		},
	}

	isThreat := true
	reason := "SQLi attack"
	confidence := 0.99

	response := protocol.AnalysisResponse{
		IsThreat:   &isThreat,
		Reason:     &reason,
		Confidence: &confidence,
	}

	headersBytes, _ := json.Marshal(request.Headers)
	metaDataBytes, _ := json.Marshal(request.MetaData)

	projectID := "project_123"

	t.Run("insert threat log into db", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO threat_logs").
			WithArgs(
				projectID,
				request.IP,
				request.Route,
				"POST",
				headersBytes,
				metaDataBytes,
				request.Log,
				*response.IsThreat,
				*response.Reason,
				*response.Confidence,
			).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := store.SaveThreat(context.Background(), projectID, request, response)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})

	t.Run("detect error when insertion fails", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO threat_logs").
			WithArgs(
				projectID,
				request.IP,
				request.Route,
				"POST",
				headersBytes,
				metaDataBytes,
				request.Log,
				*response.IsThreat,
				*response.Reason,
				*response.Confidence,
			).
			WillReturnError(errors.New("db connection lost"))

		err := store.SaveThreat(context.Background(), projectID, request, response)

		if err == nil {
			t.Error("expected error, got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})

	t.Run("use default UNKNOWN when method header is missing", func(t *testing.T) {
		reqMissingMethod := request
		reqMissingMethod.Headers = map[string]string{}
		missingMethodHeadersBytes, _ := json.Marshal(reqMissingMethod.Headers)
		missingMethodMetadataBytes, _ := json.Marshal(reqMissingMethod.MetaData)

		mock.ExpectExec("INSERT INTO threat_logs").
			WithArgs(
				projectID,
				reqMissingMethod.IP,
				reqMissingMethod.Route,
				"UNKNOWN",
				missingMethodHeadersBytes,
				missingMethodMetadataBytes,
				reqMissingMethod.Log,
				*response.IsThreat,
				*response.Reason,
				*response.Confidence,
			).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := store.SaveThreat(context.Background(), projectID, reqMissingMethod, response)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})
}

func TestSupabaseStore_ProjectManagement(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	defer mock.Close()

	store := NewSupabaseStore(mock)
	ctx := context.Background()

	t.Run("CreateProject", func(t *testing.T) {
		userID := "user_123"
		name := "My First Project"
		expectedID := "proj_uuid"
		expectedTime := time.Now()

		mock.ExpectQuery("INSERT INTO projects").
			WithArgs(userID, name, pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{"id", "created_at"}).AddRow(expectedID, expectedTime))

		proj, err := store.CreateProject(ctx, userID, name)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if proj.ID != expectedID {
			t.Errorf("expected id %s, got %s", expectedID, proj.ID)
		}
		if proj.UserID != userID {
			t.Errorf("expected user_id %s, got %s", userID, proj.UserID)
		}
		if len(proj.APIKey) < 6 || proj.APIKey[:6] != "argus_" {
			t.Errorf("malformed api key: %s", proj.APIKey)
		}
	})

	t.Run("GetProjectIDByKey - Success", func(t *testing.T) {
		apiKey := "argus_valid_key"
		expectedID := "proj_uuid"

		mock.ExpectQuery("SELECT id FROM projects").
			WithArgs(apiKey).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(expectedID))

		id, err := store.GetProjectIDByKey(ctx, apiKey)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id != expectedID {
			t.Errorf("expected id %s, got %s", expectedID, id)
		}
	})

	t.Run("GetProjectIDByKey - Invalid Key", func(t *testing.T) {
		apiKey := "argus_invalid_key"

		mock.ExpectQuery("SELECT id FROM projects").
			WithArgs(apiKey).
			WillReturnError(pgx.ErrNoRows)

		_, err := store.GetProjectIDByKey(ctx, apiKey)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "invalid api key" {
			t.Errorf("expected 'invalid api key', got %v", err)
		}
	})

	t.Run("GetProjectsByUser", func(t *testing.T) {
		userID := "user_123"

		rows := pgxmock.NewRows([]string{"id", "user_id", "name", "api_key", "created_at"}).
			AddRow("p1", userID, "Project A", "key_a", time.Now()).
			AddRow("p2", userID, "Project B", "key_b", time.Now())

		mock.ExpectQuery("SELECT id, user_id, name, api_key, created_at FROM projects").
			WithArgs(userID).
			WillReturnRows(rows)

		projects, err := store.GetProjectsByUser(ctx, userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(projects) != 2 {
			t.Errorf("expected 2 projects, got %d", len(projects))
		}
		if projects[0].Name != "Project A" {
			t.Errorf("expected first project name 'Project A', got %s", projects[0].Name)
		}
	})
}
