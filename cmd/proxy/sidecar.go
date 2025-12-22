package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/priyansh-dimri/argus/pkg/argus"
)

func main() {
	targetAddr := getEnv("TARGET_URL", "http://localhost:5000")
	listenPort := getEnv("SIDECAR_PORT", "8000")
	apiKey := getEnv("ARGUS_API_KEY", "")
	argusAPIURL := getEnv("ARGUS_API_URL", "http://localhost:8080")

	if apiKey == "" {
		log.Fatal("ARGUS_API_KEY is required")
	}

	ascii, err := os.ReadFile("ascii.txt")
	if err == nil {
		fmt.Println(string(ascii))
		fmt.Println()
	}

	fmt.Printf("🛡️ Argus Sidecar v1.0 | Port: %s\n", listenPort)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	target, err := url.Parse(targetAddr)
	if err != nil {
		log.Fatalf("Invalid TARGET_URL: %v", err)
	}

	waf, err := argus.NewWAF()
	if err != nil {
		log.Fatalf("Error initializing WAF: %v", err)
	}
	client := argus.NewClient(argusAPIURL, apiKey, 20*time.Second)

	mwLatency := argus.NewMiddleware(client, waf, argus.Config{Mode: argus.LatencyFirst})
	mwSmart := argus.NewMiddleware(client, waf, argus.Config{Mode: argus.SmartShield})
	mwParanoid := argus.NewMiddleware(client, waf, argus.Config{Mode: argus.Paranoid})

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error: %v", err)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/latency-first/", stripAndProtect(mwLatency, "/latency-first", proxy))
	mux.Handle("/smart-shield/", stripAndProtect(mwSmart, "/smart-shield", proxy))
	mux.Handle("/paranoid/", stripAndProtect(mwParanoid, "/paranoid", proxy))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Argus Multi-Mode Sidecar Active.\nUse /latency-first/, /smart-shield/, or /paranoid/ as your entry point.")
	})

	fmt.Printf("Argus Sidecar Dynamic Routing Active\n")
	fmt.Printf("Proxying to: %s\n", targetAddr)
	fmt.Printf("Listening on :%s\n", listenPort)

	log.Fatal(http.ListenAndServe(":"+listenPort, mux))
}

func stripAndProtect(mw *argus.Middleware, prefix string, next http.Handler) http.Handler {
	protected := mw.Protect(next)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		protected.ServeHTTP(w, r)
	})
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
