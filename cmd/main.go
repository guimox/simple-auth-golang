package main

import (
	"github.com/guimox/simple-auth-golang/db"
	"github.com/guimox/simple-auth-golang/internal/config"
	"github.com/guimox/simple-auth-golang/internal/handlers"
	"github.com/guimox/simple-auth-golang/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db.InitDB(cfg.DatabaseURL)
	defer db.DB.Close() // Ensure the database connection is closed when the program exits

	// Initialize repositories
	repos := config.InitializeRepositories()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(repos.UserRepo, repos.TokenRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(repos.TokenRepo)

	// Setup routes
	router := config.SetupRoutes(authHandler, authMiddleware)

	// Start server
	config.StartServer(cfg.ServerPort, router)
}
