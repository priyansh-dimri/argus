package argus

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCoraza(t *testing.T) {
	t.Run("successfully block SQLi from request URL", func(t *testing.T) {
		waf, err := NewWAF()
		if err != nil {
			t.Fatalf("Failed to init WAF: %v", err)
		}

		payload := "' OR 1=1"
		path := "/search?q=" + url.QueryEscape(payload)
		req := httptest.NewRequest("GET", path, nil)

		isThreat, err := waf.Check(req)
		if err != nil {
			t.Fatalf("WAF check failed: %v", err)
		}

		if !isThreat {
			t.Error("Expected SQLi to be blocked, but it passed")
		}
	})

	t.Run("successfully block XSS attack from request URL", func(t *testing.T) {
		waf, err := NewWAF()
		if err != nil {
			t.Fatalf("Failed to init WAF: %v", err)
		}

		payload := "<script>alert(1)</script>"
		path := "/search?q=" + url.QueryEscape(payload)
		req := httptest.NewRequest("GET", path, nil)

		isThreat, err := waf.Check(req)
		if err != nil {
			t.Fatalf("WAF check failed: %v", err)
		}

		if !isThreat {
			t.Error("Expected XSS to be blocked, but it passed")
		}
	})

	t.Run("allow clean requests", func(t *testing.T) {
		waf, err := NewWAF()
		if err != nil {
			t.Fatalf("Failed to init WAF: %v", err)
		}

		clean_req := httptest.NewRequest("GET", "/search?q=hello_world", nil)

		isThreat, err := waf.Check(clean_req)
		if err != nil {
			t.Fatalf("WAF check failed: %v", err)
		}

		if isThreat {
			t.Error("Expected clean request to pass, but it was blocked")
		}
	})

	t.Run("Successfully block SQLi from request body", func(t *testing.T) {
		waf, err := NewWAF()
		if err != nil {
			t.Fatalf("Failed to init WAF: %v", err)
		}

		body := strings.NewReader(`{"username": "admin", "password": "' OR 1=1"}`)
		req := httptest.NewRequest("POST", "/login", body)
		req.Header.Set("Content-Type", "application/json")

		isThreat, err := waf.Check(req)
		if err != nil {
			t.Fatalf("WAF check failed: %v", err)
		}

		if !isThreat {
			t.Error("Expected Body SQLi to be blocked, but it passed")
		}
	})

	t.Run("Successfully block XSS attack from request body", func(t *testing.T) {
		waf, err := NewWAF()
		if err != nil {
			t.Fatalf("Failed to init WAF: %v", err)
		}

		safe_body := strings.NewReader(`{"username": "johndoe", "bio": "Just a normal user."}`)
		req := httptest.NewRequest("POST", "/profile", safe_body)
		req.Header.Set("Content-Type", "application/json")

		isThreat, err := waf.Check(req)
		if err != nil {
			t.Fatalf("WAF check failed: %v", err)
		}

		if isThreat {
			t.Error("Expected clean body to pass, but it was blocked")
		}
	})
}
