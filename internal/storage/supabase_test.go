package storage

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

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

	t.Run("insert threat log into db", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO threat_logs").
			WithArgs(
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

		err := store.SaveThreat(context.Background(), request, response)

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

		err := store.SaveThreat(context.Background(), request, response)

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

		err := store.SaveThreat(context.Background(), reqMissingMethod, response)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectation not met: %s", err)
		}
	})
}
