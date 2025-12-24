package main

import (
	"log"
	"os"
	"path/filepath"

	api "github.com/kvloginov/cup-of-team/backend/internal/api/handlers"
	"github.com/kvloginov/cup-of-team/backend/internal/infra/db"
	"github.com/kvloginov/cup-of-team/backend/internal/infra/http"
	"github.com/kvloginov/cup-of-team/backend/internal/infra/repository"
	"github.com/kvloginov/cup-of-team/backend/internal/usecase/team"
)

func main() {
	// Get configuration from environment variables
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "db/cup-of-team.db")

	// Ensure db directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create db directory: %v", err)
	}

	// Initialize database
	database, err := db.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()
	log.Printf("Database initialized at %s", dbPath)

	// Create repository
	repo := repository.New(database.DB)

	// Create usecases
	teamUsecase := team.NewUsecase(repo)

	// Create handlers
	handlers := api.NewHandlers(teamUsecase)

	// Create server
	server := http.NewServer(http.Config{
		Port: ":" + port,
	})

	// Register routes (API routes without /api prefix, it will be added automatically)
	handlers.RegisterRoutes(server)

	// Start server
	log.Printf("🚀 API server starting on http://localhost:%s", port)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
