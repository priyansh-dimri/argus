export interface ADR {
  id: string;
  title: string;
  problem: string;
  alternatives: { name: string; tradeoff: string }[];
  decision: string;
  consequences: {
    positive: string[];
    negative: string[];
  };
}

export const adrs: ADR[] = [
  {
    id: "adr-001",
    title: "WAF Singleton with sync.Once",
    problem:
      "Coraza WAF initialization loads 6 CRS rule files (~2MB) and compiles regex patterns. Creating per-request instances causes 20ms+ startup latency and excessive memory allocation.",
    alternatives: [
      {
        name: "Per-request WAF",
        tradeoff:
          "Simplest code but 20ms initialization overhead per request, making <10µs latency impossible.",
      },
      {
        name: "sync.Pool of WAF instances",
        tradeoff:
          "Reduces init cost but requires pool management. CRS rules are immutable—pooling doesn't help parallelism.",
      },
      {
        name: "Global singleton with sync.Once",
        tradeoff:
          "Thread-safe initialization, zero per-request cost. Transactions are already goroutine-safe.",
      },
    ],
    decision:
      "Implemented singleton WAF using sync.Once in rules.go. Single waf.NewTransaction() creates isolated transaction contexts per request, avoiding shared state issues.",
    consequences: {
      positive: [
        "262µs average latency (vs 20ms+ with per-request init)",
        "4K RPS sustained throughput on single core",
        "56% parallel efficiency at 4 cores despite shared instance",
      ],
      negative: [
        "Global state complicates testing (requires mocking)",
        "Cannot swap CRS rules without restart (acceptable for production)",
      ],
    },
  },
  {
    id: "adr-002",
    title: "Protection Modes with Mode Router",
    problem:
      "Different use cases require different latency/security trade-offs. High-traffic APIs need async analysis; financial apps need synchronous AI validation.",
    alternatives: [
      {
        name: "Single hardcoded mode",
        tradeoff:
          "Simplest implementation but forces one-size-fits-all approach, losing flexibility.",
      },
      {
        name: "Runtime mode switching via API",
        tradeoff:
          "Maximum flexibility but adds API surface, race conditions, and unclear failure semantics.",
      },
      {
        name: "Config-driven mode with compile-time validation",
        tradeoff:
          "Validates modes at startup, prevents runtime errors. Requires restart to change (acceptable).",
      },
    ],
    decision:
      "Implemented three modes in config.go (LATENCY_FIRST, SMART_SHIELD, PARANOID) with switch-case router in middleware.go. Mode selected via Config struct at initialization.",
    consequences: {
      positive: [
        "LATENCY_FIRST achieves 87µs middleware overhead (pure async)",
        "PARANOID provides 100% AI coverage for zero-day threats",
        "Clear separation of concerns; testable in isolation",
      ],
      negative: [
        "Mode change requires app restart (mitigated by container orchestration)",
        "Three code paths increase maintenance surface",
      ],
    },
  },
  {
    id: "adr-003",
    title: "Circuit Breaker for AI API Resilience",
    problem:
      "Gemini AI API can fail due to rate limits, network issues, or service outages. Cascading failures to AI would block all requests in PARANOID mode.",
    alternatives: [
      {
        name: "No resilience pattern",
        tradeoff:
          "Simplest but causes total service failure when AI is down. Unacceptable for production.",
      },
      {
        name: "Retry with exponential backoff",
        tradeoff:
          "Handles transient errors but amplifies load during outages, worsening cascades.",
      },
      {
        name: "Circuit breaker (gobreaker/v2)",
        tradeoff:
          "Fails fast after threshold (3 consecutive failures), prevents cascade. 40ns per-call overhead.",
      },
    ],
    decision:
      "Integrated gobreaker with 3-failure threshold, 30s timeout, 60s reset interval in breaker.go. Wraps Client.SendAnalysis() calls in all modes.",
    consequences: {
      positive: [
        "40ns per-call overhead (measured in benchmarks)",
        "HALF_OPEN state allows graceful recovery",
        "Prevents thundering herd to AI service during outages",
      ],
      negative: [
        "False positives possible if AI has brief (<30s) hiccup",
        "Requires monitoring circuit state (not yet implemented)",
      ],
    },
  },
  {
    id: "adr-004",
    title: "Body Buffering with Reset for WAF + AI",
    problem:
      "http.Request.Body is an io.ReadCloser (stream) consumed by first read. WAF needs body for inspection, AI analysis needs same body, and upstream handler needs original body.",
    alternatives: [
      {
        name: "Read body twice via TeeReader",
        tradeoff:
          "Requires network re-read or caching. Complex error handling and doubled memory usage.",
      },
      {
        name: "Custom io.ReadCloser wrapper",
        tradeoff:
          "Allows re-reading but needs careful state management. High implementation complexity.",
      },
      {
        name: "io.ReadAll + io.NopCloser reset",
        tradeoff:
          "Simple buffer-once, reset-many approach. 67KB avg memory allocation per request.",
      },
    ],
    decision:
      "Buffer body with io.ReadAll in middleware.Protect(), store []byte, reset via io.NopCloser(bytes.NewBuffer(bodyBytes)) after each read.",
    consequences: {
      positive: [
        "Single memory allocation per request (predictable)",
        "WAF, AI, and handler all see identical body",
        "Stdlib-only solution (no external deps)",
      ],
      negative: [
        "67KB allocation per POST request (acceptable for threat API)",
        "Not suitable for streaming uploads",
      ],
    },
  },
  {
    id: "adr-005",
    title: "Async Threat Storage with Goroutines",
    problem:
      "Saving threat logs to Supabase (INSERT + JSON marshaling) adds 50-100ms latency. Blocking on DB writes destroys latency targets.",
    alternatives: [
      {
        name: "Synchronous DB write in request path",
        tradeoff:
          "Simplest code but adds 50-100ms to every request. Unacceptable for <10µs goals.",
      },
      {
        name: "Message queue (Kafka, RabbitMQ)",
        tradeoff:
          "Production-grade durability but requires separate service, adds operational complexity.",
      },
      {
        name: "Background goroutine with fire-and-forget",
        tradeoff:
          "Zero request latency impact. Risk: goroutine leaks if uncontrolled.",
      },
    ],
    decision:
      "Fire goroutine with go func() in handlers.HandleAnalyze(). Pass background context (not request context) to avoid cancellation. Log errors via ErrorReporter.",
    consequences: {
      positive: [
        "0µs latency impact on request path",
        "Natural backpressure via DB connection pool",
        "Simple implementation (15 lines of code)",
      ],
      negative: [
        "No delivery guarantee if app crashes before DB write",
      ],
    },
  },
  {
    id: "adr-006",
    title: "Token-Based Authentication Strategy",
    problem:
      "Two distinct auth needs: SDK clients (stateless API keys) and dashboard users (session-based JWT). Shared middleware would couple concerns.",
    alternatives: [
      {
        name: "Single OAuth2 flow for both",
        tradeoff:
          "Unified auth but forces SDK users into complex OAuth leading to poor DX.",
      },
      {
        name: "Separate middleware per route group",
        tradeoff:
          "Clean separation but duplicates CORS, context injection, error handling logic.",
      },
      {
        name: "Dual middleware with shared context injection",
        tradeoff:
          "Reuses context.WithValue pattern. AuthSDK for /analyze, AuthDashboard for /projects.",
      },
    ],
    decision:
      "Implemented AuthSDK (Bearer API key → projectID) and AuthDashboard (JWT → userID) in api/middleware.go. Both inject IDs into context for downstream handlers.",
    consequences: {
      positive: [
        "SDK auth: single database lookup via GetProjectIDByKey",
        "Dashboard auth: JWT validation with Supabase HMAC secret",
      ],
      negative: [
        "Two auth code paths to maintain",
      ],
    },
  },
  {
    id: "adr-007",
    title: "Metadata-Driven Context Awareness",
    problem:
      "AI must distinguish between malicious payloads and legitimate user content (e.g., blog posts explaining SQL injection). False positives destroy UX.",
    alternatives: [
      {
        name: "Pure payload analysis",
        tradeoff:
          "Simplest but causes false positives on tutorial content, code snippets, etc.",
      },
      {
        name: "Allowlist by route",
        tradeoff:
          "Works for known safe routes but brittle, requires constant updates.",
      },
      {
        name: "Client-provided metadata with AI prompt rules",
        tradeoff:
          "Clients tag context (e.g., 'blog_editor'). AI prompt explicitly trusts metadata.",
      },
    ],
    decision:
      "Added MetaData map[string]string to AnalysisRequest in protocol. Prompt in prompt.go states: 'Metadata is provided by authentic users so you MUST trust them always.'",
    consequences: {
      positive: [
        "Clients control context without code changes",
      ],
      negative: [
        "Metadata can be spoofed by malicious clients",
      ],
    },
  },
];
