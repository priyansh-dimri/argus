package argus

import (
	"bytes"
	"io"
	"net/http"

	"github.com/priyansh-dimri/argus/pkg/protocol"
)

type Middleware struct {
	Client AnalysisSender
	WAF    RuleEngine
	Config Config
}

func NewMiddleware(client AnalysisSender, waf RuleEngine, config Config) *Middleware {
	return &Middleware{
		Client: client,
		WAF:    waf,
		Config: config,
	}
}

func (m *Middleware) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		resetBody := func() {
			if bodyBytes != nil {
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		wafResult, _ := m.WAF.Check(r)

		resetBody()

		switch m.Config.Mode {
		case LatencyFirst:
			m.handleLatencyFirst(w, r, next, wafResult, bodyBytes)
		case Paranoid:
			m.handleParanoid(w, r, next, wafResult, bodyBytes)
		case SmartShield:
			fallthrough
		default:
			m.handleSmartShield(w, r, next, wafResult, bodyBytes)
		}
	})
}

func (m *Middleware) handleLatencyFirst(w http.ResponseWriter, r *http.Request, next http.Handler, wafBlocked bool, body []byte) {
	if wafBlocked {
		go m.sendAsyncLog(r, body, wafBlocked)
		http.Error(w, "Blocked by Argus Shield", http.StatusForbidden)
		return
	}
	go m.sendAsyncLog(r, body, wafBlocked)
	next.ServeHTTP(w, r)
}

func (m *Middleware) handleSmartShield(w http.ResponseWriter, r *http.Request, next http.Handler, wafBlocked bool, body []byte) {
	if !wafBlocked {
		go m.sendAsyncLog(r, body, wafBlocked)
		next.ServeHTTP(w, r)
		return
	}

	resp, err := m.sendSyncAnalysis(r, body, wafBlocked)

	if err == nil && !*resp.IsThreat {
		next.ServeHTTP(w, r)
		return
	}

	http.Error(w, "Blocked by Argus Smart Shield", http.StatusForbidden)
}

func (m *Middleware) handleParanoid(w http.ResponseWriter, r *http.Request, next http.Handler, wafBlocked bool, body []byte) {
	resp, err := m.sendSyncAnalysis(r, body, wafBlocked)
	if err == nil && *resp.IsThreat {
		http.Error(w, "Blocked by Argus Paranoid Shield", http.StatusForbidden)
		return
	}

	next.ServeHTTP(w, r)
}

func (m *Middleware) buildPayload(r *http.Request, body []byte, wafBlocked bool) protocol.AnalysisRequest {
	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	headers["Method"] = r.Method

	meta := map[string]string{
		"waf_result": "PASS",
	}
	if wafBlocked {
		meta["waf_result"] = "BLOCK"
	}

	return protocol.AnalysisRequest{
		Log:      string(body),
		IP:       r.RemoteAddr,
		Route:    r.URL.Path,
		Headers:  headers,
		MetaData: meta,
	}
}

func (m *Middleware) sendAsyncLog(r *http.Request, body []byte, wafBlocked bool) {
	req := m.buildPayload(r, body, wafBlocked)
	_, _ = m.Client.SendAnalysis(req)
}

func (m *Middleware) sendSyncAnalysis(r *http.Request, body []byte, wafBlocked bool) (protocol.AnalysisResponse, error) {
	req := m.buildPayload(r, body, wafBlocked)
	return m.Client.SendAnalysis(req)
}
