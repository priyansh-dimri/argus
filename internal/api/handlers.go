package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type Analyzer interface {
	Analyze(ctx context.Context, req protocol.AnalysisRequest) (protocol.AnalysisResponse, error)
}

type Store interface{}

type API struct {
	Analyzer Analyzer
	Store    Store
}

func (api *API) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
}
