package argus

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/corazawaf/coraza/v3"
)

type RuleEngine interface {
	Check(r *http.Request) (bool, error)
}

type WAFWrapper struct {
	waf coraza.WAF
}

var _ RuleEngine = (*WAFWrapper)(nil)

//go:embed rules/*.conf rules/*.data
var rulesFS embed.FS

func NewWAF() (*WAFWrapper, error) {
	cfg := coraza.NewWAFConfig()

	root, _ := fs.Sub(rulesFS, "rules")
	cfg = cfg.WithRootFS(root)

	files := []string{
		"coraza.conf",
		"crs-setup.conf",
		"REQUEST-901-INITIALIZATION.conf",
		"REQUEST-941-APPLICATION-ATTACK-XSS.conf",
		"REQUEST-942-APPLICATION-ATTACK-SQLI.conf",
		"REQUEST-949-BLOCKING-EVALUATION.conf",
	}

	var err error
	for _, file := range files {
		cfg, err = parseRuleFile(cfg, file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rule file %s: %w", file, err)
		}
	}

	waf, err := coraza.NewWAF(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize coraza waf: %w", err)
	}

	return &WAFWrapper{waf: waf}, nil
}

func parseRuleFile(cfg coraza.WAFConfig, filename string) (coraza.WAFConfig, error) {
	f, err := rulesFS.Open("rules/" + filename)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return cfg, err
	}

	return cfg.WithDirectives(string(data)), nil
}

func (w *WAFWrapper) Check(r *http.Request) (bool, error) {
	tx := w.waf.NewTransaction()
	defer tx.Close()

	tx.ProcessConnection(r.RemoteAddr, 0, "", 0)
	tx.ProcessURI(r.URL.String(), r.Method, r.Proto)

	for k, vv := range r.Header {
		for _, v := range vv {
			tx.AddRequestHeader(k, v)
		}
	}
	if it := tx.ProcessRequestHeaders(); it != nil {
		return true, nil
	}

	if r.Body != nil && r.Body != http.NoBody {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return false, fmt.Errorf("failed to read request body: %w", err)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if _, _, err := tx.WriteRequestBody(bodyBytes); err != nil {
			return false, fmt.Errorf("failed to write body to waf: %w", err)
		}
	}

	if it, err := tx.ProcessRequestBody(); err != nil {
		return false, fmt.Errorf("failed to process request body: %w", err)
	} else if it != nil {
		return true, nil
	}

	return tx.IsInterrupted(), nil
}
