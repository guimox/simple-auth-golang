package middleware

import (
	"net/http"
	"time"

	"github.com/guimox/simple-auth-golang/internal/repository"
)

type AuthMiddleware struct {
	TokenRepo *repository.TokenRepository
}

func NewAuthMiddleware(tokenRepo *repository.TokenRepository) *AuthMiddleware {
	return &AuthMiddleware{TokenRepo: tokenRepo}
}

func (m *AuthMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized: No session token found", http.StatusUnauthorized)
			return
		}

		token, err := m.TokenRepo.GetToken(sessionCookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid session token", http.StatusUnauthorized)
			return
		}

		if token.ExpiresAt.Before(time.Now()) {
			http.Error(w, "Unauthorized: Session token expired", http.StatusUnauthorized)
			return
		}

		csrf := r.Header.Get("csrf_token")
		if csrf == "" {
			http.Error(w, "Unauthorized: Missing CSRF token", http.StatusUnauthorized)
			return
		}

		if csrf != token.CSRFToken {
			http.Error(w, "Unauthorized: CSRF token mismatch", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
