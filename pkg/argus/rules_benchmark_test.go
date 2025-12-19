package argus

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

var benchmarkResult bool

func BenchmarkWAF_Check_Clean(b *testing.B) {
	waf, err := NewWAF()
	if err != nil {
		b.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/api/v1/users?id=12345&sort=asc", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		isThreat, err := waf.Check(req)
		if err != nil {
			b.Fatalf("Check failed: %v", err)
		}
		benchmarkResult = isThreat
	}
}

func BenchmarkWAF_Check_SQLi(b *testing.B) {
	waf, err := NewWAF()
	if err != nil {
		b.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/api/v1/search?q='%20OR%201=1%20--", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		isThreat, err := waf.Check(req)
		if err != nil {
			b.Fatalf("Check failed: %v", err)
		}
		benchmarkResult = isThreat
	}
}

func BenchmarkWAF_Check_POST_JSON(b *testing.B) {
	waf, err := NewWAF()
	if err != nil {
		b.Fatal(err)
	}

	payload := `{"username": "admin", "bio": "Standard user bio with enough length to trigger body inspection rules.", "settings": {"theme": "dark", "notifications": true}}`

	req := httptest.NewRequest("POST", "/api/v1/profile", nil)
	req.Header.Set("Content-Type", "application/json")

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		req.Body = io.NopCloser(strings.NewReader(payload))
		isThreat, err := waf.Check(req)
		if err != nil {
			b.Fatalf("Check failed: %v", err)
		}
		benchmarkResult = isThreat
	}
}

func BenchmarkWAF_Check_XSS(b *testing.B) {
	waf, err := NewWAF()
	if err != nil {
		b.Fatal(err)
	}

	payload := `{"bio": "<script>alert(1)</script>"}`

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		req := httptest.NewRequest("POST", "/profile", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		isThreat, err := waf.Check(req)
		if err != nil {
			b.Fatalf("Check failed: %v", err)
		}
		benchmarkResult = isThreat
	}
}

func BenchmarkWAF_Check_Parallel(b *testing.B) {
	waf, err := NewWAF()
	if err != nil {
		b.Fatal(err)
	}

	payload := `{"key": "value"}`

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("POST", "/api/data", strings.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")

			isThreat, err := waf.Check(req)
			if err != nil {
				b.Fatalf("Check failed: %v", err)
			}
			benchmarkResult = isThreat
		}
	})
}

func BenchmarkWAF_Check_Parallel_SQLi(b *testing.B) {
	waf, err := NewWAF()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/api/search?q='%20OR%201=1%20--", nil)

			isThreat, err := waf.Check(req)
			if err != nil {
				b.Fatalf("Check failed: %v", err)
			}
			benchmarkResult = isThreat
		}
	})
}

func BenchmarkWAF_Check_Parallel_XSS(b *testing.B) {
	waf, err := NewWAF()
	if err != nil {
		b.Fatal(err)
	}

	payload := `{"bio": "<script>alert(1)</script>"}`

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("POST", "/profile", strings.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")

			isThreat, err := waf.Check(req)
			if err != nil {
				b.Fatalf("Check failed: %v", err)
			}
			benchmarkResult = isThreat
		}
	})
}
