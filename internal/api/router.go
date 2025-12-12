package api

import "net/http"

func NewRouter(api *API) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /analyze", api.HandleAnalyze)

	return mux
}
