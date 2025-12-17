package api

import "net/http"

func NewRouter(api *API) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /analyze", api.HandleAnalyze)
	mux.HandleFunc("POST /projects", api.HandleCreateProject)
	mux.HandleFunc("GET /projects", api.HandleListProjects)
	return mux
}
