package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrAuth = errors.New("Unauthorized")

func Authorize(r *http.Request) error {
	// Get session cookie
	st, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			fmt.Println("Authorization failed: No session cookie found")
		} else {
			fmt.Println("Authorization failed: Error retrieving session cookie -", err)
		}
		return ErrAuth
	}

	if st.Value == "" {
		fmt.Println("Authorization failed: Session token is empty")
		return ErrAuth
	}

	// Split token into parts (expected format: "randomToken:username")
	parts := strings.Split(st.Value, ":")
	if len(parts) != 2 {
		fmt.Println("Authorization failed: Invalid session token format")
		return ErrAuth
	}

	token, username := parts[0], parts[1]

	// Check if the user exists
	user, ok := users[username]
	if !ok {
		fmt.Printf("Authorization failed: User '%s' not found\n", username)
		return ErrAuth
	}

	// Validate session token
	storedParts := strings.Split(user.SessionToken, ":")
	if len(storedParts) != 2 || storedParts[0] != token {
		fmt.Println("Authorization failed: Session token mismatch")
		return ErrAuth
	}

	// Get CSRF token from the request headers
	csrf := r.Header.Get("csrf_token")
	if csrf == "" {
		fmt.Println("Authorization failed: Missing CSRF token in request")
		return ErrAuth
	}

	// Validate CSRF token
	if csrf != user.CSRFToken {
		fmt.Println("Authorization failed: CSRF token mismatch")
		return ErrAuth
	}

	fmt.Println("Authorization successful for user:", username)
	return nil
}
