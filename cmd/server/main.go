package main

import (
	"context"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/priyansh-dimri/argus/internal/analyzer"
	"github.com/priyansh-dimri/argus/internal/api"
	"github.com/priyansh-dimri/argus/internal/storage"
	"github.com/priyansh-dimri/argus/pkg/logger"
)

func main() {
	logger.InitLogger()

	geminiKey := os.Getenv("GEMINI_API_KEY")
	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if geminiKey == "" || dbURL == "" {
		logger.Error("Missing environment variables", nil, "GEMINI_API_KEY", geminiKey != "", "DATABASE_URL", dbURL != "")
		os.Exit(1)
	}

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Error("Failed to connect to database", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	store := storage.NewSupabaseStore(dbPool)

	// AI Client
	aiClient, err := analyzer.NewGeminiClient(ctx, geminiKey, "gemini-2.5-flash")
	if err != nil {
		logger.Error("Failed to initialize Gemini client", err)
		os.Exit(1)
	}

	core := analyzer.NewAnalyzer(aiClient)
	handler := api.NewAPI(core, store)

	router := api.NewRouter(handler)

	logger.Info("Starting Argus API", "port", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Error("Server failed", err)
		os.Exit(1)
	}
}
