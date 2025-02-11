package auth

import (
	"errors"
	"net/http"

	"github.com/guimox/simple-auth-golang/internal/models"
	"github.com/guimox/simple-auth-golang/internal/store"
)

var ErrUnauthorized = errors.New("unauthorized")

func Authorize(r *http.Request) (*models.User, error) {
	// Get session token from cookie
	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" {
		return nil, ErrUnauthorized
	}

	// Retrieve session from storage
	session, err := store.GetSession(st.Value)
	if err != nil {
		return nil, ErrUnauthorized
	}

	// Retrieve user from storage
	user, err := store.GetUser(session.Username)
	if err != nil {
		return nil, ErrUnauthorized
	}

	// Validate CSRF token
	csrf := r.Header.Get("X-CSRF-Token")
	if csrf != user.CRSFToken || csrf == "" {
		return nil, ErrUnauthorized
	}

	return &user, nil
}
