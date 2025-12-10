package protocol

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
