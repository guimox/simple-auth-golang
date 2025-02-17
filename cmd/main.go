package main

import (
	"github.com/guimox/simple-auth-golang/db"
	"github.com/guimox/simple-auth-golang/internal/config"
	"github.com/guimox/simple-auth-golang/internal/handlers"
	"github.com/guimox/simple-auth-golang/internal/middleware"
)

func main() {
	cfg := config.LoadConfig()

	db.InitDB(cfg.DatabaseURL)
	defer db.DB.Close()

	repos := config.InitializeRepositories()

	authHandler := handlers.NewAuthHandler(repos.UserRepo, repos.TokenRepo)

	authMiddleware := middleware.NewAuthMiddleware(repos.TokenRepo)

	router := config.SetupRoutes(authHandler, authMiddleware)

	config.StartServer(cfg.ServerPort, router)
}
