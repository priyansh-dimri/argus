package storage

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
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

	t.Run("create project successfully", func(t *testing.T) {
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

	t.Run("detect api key gen failure on create project", func(t *testing.T) {
		store := NewSupabaseStore(mock)
		store.randRead = func(b []byte) (n int, err error) {
			return 0, errors.New("entropy error")
		}

		_, err := store.CreateProject(ctx, "user_123", "Fail Project")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to generate api key") {
			t.Errorf("expected api key generation error, got: %v", err)
		}
	})

	t.Run("detect database insertion failure on create project", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO projects").
			WithArgs("user_123", "Fail Project", pgxmock.AnyArg()).
			WillReturnError(errors.New("db insert failed"))

		_, err := store.CreateProject(ctx, "user_123", "Fail Project")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to insert project") {
			t.Errorf("expected db insert error, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})

	t.Run("get project id by key", func(t *testing.T) {
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

	t.Run("detect invalid key on get project", func(t *testing.T) {
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

	t.Run("detect database error on get project by key", func(t *testing.T) {
		apiKey := "argus_db_error_key"

		mock.ExpectQuery("SELECT id FROM projects").
			WithArgs(apiKey).
			WillReturnError(errors.New("db connection lost"))

		_, err := store.GetProjectIDByKey(ctx, apiKey)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "database error") {
			t.Errorf("expected 'database error' wrapper, got %v", err)
		}
	})

	t.Run("get project by user successfully", func(t *testing.T) {
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

	t.Run("detect query failure on get projects by user", func(t *testing.T) {
		userID := "user_query_fail"

		mock.ExpectQuery("SELECT id, user_id, name, api_key, created_at FROM projects").
			WithArgs(userID).
			WillReturnError(errors.New("deadlock detected"))

		_, err := store.GetProjectsByUser(ctx, userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to query projects") {
			t.Errorf("expected 'failed to query projects' wrapper, got %v", err)
		}
	})

	t.Run("detect scan failure on get projects by user", func(t *testing.T) {
		userID := "user_scan_fail"

		rows := pgxmock.NewRows([]string{"id", "user_id", "name", "api_key", "created_at"}).
			AddRow("p1", userID, "Project A", "key_a", "NOT_A_TIMESTAMP")

		mock.ExpectQuery("SELECT id, user_id, name, api_key, created_at FROM projects").
			WithArgs(userID).
			WillReturnRows(rows)

		_, err := store.GetProjectsByUser(ctx, userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to scan project row") {
			t.Errorf("expected 'failed to scan project row' wrapper, got %v", err)
		}
	})

	t.Run("update project name successfully", func(t *testing.T) {
		projectID := "project_123"
		newName := "Renamed Project"

		mock.ExpectExec("UPDATE projects SET name").
			WithArgs(newName, projectID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := store.UpdateProjectName(ctx, projectID, newName)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})

	t.Run("detect database error on update project name", func(t *testing.T) {
		projectID := "project_123"
		newName := "Renamed Project"

		mock.ExpectExec("UPDATE projects SET name").
			WithArgs(newName, projectID).
			WillReturnError(errors.New("db update failed"))

		err := store.UpdateProjectName(ctx, projectID, newName)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to update project name") {
			t.Errorf("expected 'failed to update project name' wrapper, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})

	t.Run("detect not found on update project name", func(t *testing.T) {
		projectID := "proj_missing"
		newName := "Renamed Project"

		mock.ExpectExec("UPDATE projects SET name").
			WithArgs(newName, projectID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err := store.UpdateProjectName(ctx, projectID, newName)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "project not found" {
			t.Errorf("expected 'project not found', got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})

	t.Run("rotate api key successfully", func(t *testing.T) {
		projectID := "project_123"

		fixed := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
		expectedKey := "argus_" + hex.EncodeToString(fixed)

		store := NewSupabaseStore(mock)
		store.randRead = func(b []byte) (n int, err error) {
			copy(b, fixed)
			return len(b), nil
		}

		mock.ExpectExec("UPDATE projects SET api_key").
			WithArgs(expectedKey, projectID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		newKey, err := store.RotateAPIKey(ctx, projectID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if newKey != expectedKey {
			t.Errorf("expected new key %s, got %s", expectedKey, newKey)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})

	t.Run("detect api key generation failure on rotate api key", func(t *testing.T) {
		projectID := "project_123"

		store := NewSupabaseStore(mock)
		store.randRead = func(b []byte) (n int, err error) {
			return 0, errors.New("entropy error")
		}

		_, err := store.RotateAPIKey(ctx, projectID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to generate new key") {
			t.Errorf("expected 'failed to generate new key' wrapper, got %v", err)
		}
	})

	t.Run("detect database error on rotate api key", func(t *testing.T) {
		projectID := "project_123"

		fixed := []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99}
		expectedKey := "argus_" + hex.EncodeToString(fixed)

		store := NewSupabaseStore(mock)
		store.randRead = func(b []byte) (n int, err error) {
			copy(b, fixed)
			return len(b), nil
		}

		mock.ExpectExec("UPDATE projects SET api_key").
			WithArgs(expectedKey, projectID).
			WillReturnError(errors.New("db update failed"))

		_, err := store.RotateAPIKey(ctx, projectID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to update api key") {
			t.Errorf("expected 'failed to update api key' wrapper, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})

	t.Run("detect not found on rotate api key", func(t *testing.T) {
		projectID := "proj_missing"

		fixed := []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}
		expectedKey := "argus_" + hex.EncodeToString(fixed)

		store := NewSupabaseStore(mock)
		store.randRead = func(b []byte) (n int, err error) {
			copy(b, fixed)
			return len(b), nil
		}

		mock.ExpectExec("UPDATE projects SET api_key").
			WithArgs(expectedKey, projectID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		_, err := store.RotateAPIKey(ctx, projectID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "project not found" {
			t.Errorf("expected 'project not found', got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})

	t.Run("delete project successfully", func(t *testing.T) {
		projectID := "project_123"

		mock.ExpectExec("DELETE FROM projects").
			WithArgs(projectID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := store.DeleteProject(ctx, projectID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})

	t.Run("detect database error on delete project", func(t *testing.T) {
		projectID := "project_123"

		mock.ExpectExec("DELETE FROM projects").
			WithArgs(projectID).
			WillReturnError(errors.New("db delete failed"))

		err := store.DeleteProject(ctx, projectID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to delete project") {
			t.Errorf("expected 'failed to delete project', got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})

	t.Run("detect not found on delete project", func(t *testing.T) {
		projectID := "proj_missing"

		mock.ExpectExec("DELETE FROM projects").
			WithArgs(projectID).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err := store.DeleteProject(ctx, projectID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "project not found" {
			t.Errorf("expected 'project not found', got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation are not met: %s", err)
		}
	})
}
