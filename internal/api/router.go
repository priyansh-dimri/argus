package api

import "net/http"

func NewRouter(api *API, mw *Middleware) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /analyze", mw.AuthSDK(api.HandleAnalyze))
	mux.HandleFunc("POST /projects", mw.AuthDashboard(api.HandleCreateProject))
	mux.HandleFunc("GET /projects", mw.AuthDashboard(api.HandleListProjects))
	mux.HandleFunc("PATCH /projects", mw.AuthDashboard(api.HandleUpdateProject))
	mux.HandleFunc("DELETE /projects", mw.AuthDashboard(api.HandleDeleteProject))
	mux.HandleFunc("POST /rotate-key", mw.AuthDashboard(api.HandleRotateKey))
	mux.HandleFunc("DELETE /account", mw.AuthDashboard(api.HandleDeleteAccount))

	return mux
}
