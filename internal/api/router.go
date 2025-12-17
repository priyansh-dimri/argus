package api

import "net/http"

func NewRouter(api *API, mw *Middleware) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /analyze", mw.AuthSDK(api.HandleAnalyze))
	mux.HandleFunc("POST /projects", mw.AuthDashboard(api.HandleCreateProject))
	mux.HandleFunc("GET /projects", mw.AuthDashboard(api.HandleListProjects))

	return mux
}
