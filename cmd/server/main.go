package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/priyansh-dimri/argus/internal/analyzer"
	"github.com/priyansh-dimri/argus/internal/api"
	"github.com/priyansh-dimri/argus/internal/storage"
	"github.com/priyansh-dimri/argus/pkg/logger"
)

func main() {
	startTime := time.Now()
	logger.InitLogger()
	logger.Info("Argus API Starting", "component", "main")
	logger.Info("Loading environment variables", "component", "main")

	geminiKey := os.Getenv("GEMINI_API_KEY")
	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Info("PORT not set, using default", "component", "main", "port", port)
	} else {
		logger.Info("PORT loaded from environment", "component", "main", "port", port)
	}

	if geminiKey == "" || dbURL == "" {
		logger.Error("Missing environment variables", nil, "GEMINI_API_KEY", geminiKey != "", "DATABASE_URL", dbURL != "")
		os.Exit(1)
	}

	logger.Info("All required environment variables loaded",
		"component", "main",
		"gemini_key_length", len(geminiKey),
		"db_url_length", len(dbURL),
	)

	logger.Info("Initializing database connection", "component", "main")
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Error("Failed to parse DATABASE_URL", err,
			"component", "main",
			"db_url_length", len(dbURL),
		)
		os.Exit(1)
	}

	logger.Info("Database config parsed successfully", "component", "main")

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	logger.Info("Set query execution mode",
		"component", "main",
		"mode", "SimpleProtocol",
	)

	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error("Failed to connect to database", err,
			"component", "main",
		)
		os.Exit(1)
	}

	defer dbPool.Close()

	logger.Info("Database connection pool created successfully",
		"component", "main",
		"max_conns", config.MaxConns,
		"min_conns", config.MinConns,
	)

	logger.Info("Testing database connection", "component", "main")
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := dbPool.Ping(pingCtx); err != nil {
		logger.Error("Database ping failed", err, "component", "main")
		os.Exit(1)
	}

	logger.Info("Database connection test successful", "component", "main")

	logger.Info("Initializing storage layer", "component", "main")
	store := storage.NewSupabaseStore(dbPool)

	logger.Info("Initializing Gemini AI client",
		"component", "main",
		"model", "gemini-2.5-flash",
	)

	aiClient, err := analyzer.NewGeminiClient(ctx, geminiKey, "gemini-2.5-flash")
	if err != nil {
		logger.Error("Failed to initialize Gemini client", err,
			"component", "main",
			"model", "gemini-2.5-flash",
		)
		os.Exit(1)
	}

	logger.Info("Gemini AI client initialized successfully", "component", "main")

	logger.Info("Initializing analyzer", "component", "main")
	core := analyzer.NewAnalyzer(aiClient)

	logger.Info("Initializing API handlers", "component", "main")
	handler := api.NewAPI(core, store)

	logger.Info("Initializing authentication middleware", "component", "main")
	authMiddleware := api.NewMiddleware(store)

	logger.Info("Setting up HTTP router", "component", "main")
	router := api.NewRouter(handler, authMiddleware)

	logger.Info("Wrapping router with CORS middleware", "component", "main")
	corsHandler := authMiddleware.CORS(router)

	initDuration := time.Since(startTime)
	logger.Info("Initialization complete, starting HTTP server",
		"component", "main",
		"port", port,
		"init_duration_ms", initDuration.Milliseconds(),
	)

	serverAddr := ":" + port
	logger.Info("=== Argus API Ready ===",
		"component", "main",
		"address", serverAddr,
		"startup_time_ms", initDuration.Milliseconds(),
	)

	if err := http.ListenAndServe(serverAddr, corsHandler); err != nil {
		logger.Error("HTTP server failed", err,
			"component", "main",
			"address", serverAddr,
		)
		os.Exit(1)
	}
}
