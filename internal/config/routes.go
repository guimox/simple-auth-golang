package config

import (
	"net/http"

	"github.com/guimox/simple-auth-golang/internal/handlers"
	"github.com/guimox/simple-auth-golang/internal/middleware"
)

// SetupRoutes configures and returns the HTTP router
func SetupRoutes(authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/public", handlers.Public)

	// Auth routes
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/logout", authHandler.Logout)

	// Protected route
	mux.Handle("/protected", authMiddleware.Authorize(http.HandlerFunc(handlers.Protected)))

	// Add logging middleware
	return middleware.LogMiddleware(mux)
}
