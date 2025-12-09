package protocol

// Analysis Request is sent from client to Argus backend
type AnalysisRequest struct {
	Log      string
	IP       string
	Headers  map[string]string
	Route    string
	MetaData map[string]string
}

// Analysis Response is result sent by Argus backend
type AnalysisResponse struct {
	IsThreat   bool
	Reason     string
	Confidence float64
}
