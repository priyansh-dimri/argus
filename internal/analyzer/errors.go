package analyzer

type AnalyzerErr string

func (e AnalyzerErr) Error() string {
	return string(e)
}

const (
	ErrMalformedAIResponse AnalyzerErr = "malformed AI response"
	ErrAIGenerateFailed    AnalyzerErr = "AI generate failed"
)
