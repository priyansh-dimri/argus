package protocol

import "time"

// Analysis Request is sent from client to Argus backend
type AnalysisRequest struct {
	Log      string            `json:"log"`
	IP       string            `json:"ip"`
	Headers  map[string]string `json:"headers"`
	Route    string            `json:"route"`
	MetaData map[string]string `json:"metadata"`
}

// Analysis Response is result sent by Argus backend
type AnalysisResponse struct {
	IsThreat   *bool    `json:"is_threat"`
	Reason     *string  `json:"reason"`
	Confidence *float64 `json:"confidence"`
}

type Project struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateProjectRequest struct {
	Name string `json:"name"`
}

type CreateProjectResponse struct {
	Project Project `json:"project"`
}

type DashboardStats struct {
	TotalRequests   int            `json:"total_requests"`
	BlockedRequests int            `json:"blocked_requests"`
	ThreatsByType   map[string]int `json:"threats_by_type"`
	RecentThreats   []ThreatLog    `json:"recent_threats"`
}

type ThreatLog struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	Route     string    `json:"route"`
	Reason    string    `json:"reason"`
}
