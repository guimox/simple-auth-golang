package config

import (
	"net/http"

	"github.com/guimox/simple-auth-golang/internal/handlers"
	"github.com/guimox/simple-auth-golang/internal/middleware"
)

func SetupRoutes(authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/public", handlers.Public)

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/logout", authHandler.Logout)

	mux.Handle("/protected", authMiddleware.Authorize(http.HandlerFunc(handlers.Protected)))

	return middleware.LogMiddleware(mux)
}
