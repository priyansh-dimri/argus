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
	SaveThreat(ctx context.Context, req protocol.AnalysisRequest, res protocol.AnalysisResponse) error
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
		if err := api.Store.SaveThreat(bgContext, req, res); err != nil {
			if api.ErrorReporter != nil {
				api.ErrorReporter("Failed to save threat log", "error", err)
			}
		}
	}()
}
