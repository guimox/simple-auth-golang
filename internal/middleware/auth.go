package middleware

import (
	"net/http"
	"time"

	"github.com/guimox/simple-auth-golang/internal/repository"
)

// AuthMiddleware is a middleware that checks if the request is authorized
type AuthMiddleware struct {
	TokenRepo *repository.TokenRepository
}

// NewAuthMiddleware creates a new instance of AuthMiddleware
func NewAuthMiddleware(tokenRepo *repository.TokenRepository) *AuthMiddleware {
	return &AuthMiddleware{TokenRepo: tokenRepo}
}

// Authorize checks for a valid session token and CSRF token
func (m *AuthMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized: No session token found", http.StatusUnauthorized)
			return
		}

		// Validate the session token
		token, err := m.TokenRepo.GetToken(sessionCookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid session token", http.StatusUnauthorized)
			return
		}

		// Check if the token has expired
		if token.ExpiresAt.Before(time.Now()) {
			http.Error(w, "Unauthorized: Session token expired", http.StatusUnauthorized)
			return
		}

		// Get CSRF token from the request headers
		csrf := r.Header.Get("csrf_token")
		if csrf == "" {
			http.Error(w, "Unauthorized: Missing CSRF token", http.StatusUnauthorized)
			return
		}

		// Validate CSRF token
		if csrf != token.CSRFToken {
			http.Error(w, "Unauthorized: CSRF token mismatch", http.StatusUnauthorized)
			return
		}

		// If all checks pass, call the next handler
		next.ServeHTTP(w, r)
	})
}
