package main

import (
	"errors"
	"net/http"
)

var ErrUnauthorized = errors.New("Unauthorized")

func Authorize(r *http.Request) error {
	username := r.FormValue("username")
	user, ok := users[username]
	if !ok {
		return ErrUnauthorized
	}

	st, err := r.Cookie("session-token")
	if err != nil || st.Value == "" || st.Value != user.SessionToken {
		return ErrUnauthorized
	}

	csrf := r.Header.Get("X-CSRF-Token")
	if csrf != user.CRSFToken || csrf == "" {
		return ErrUnauthorized
	}

	return nil
}
