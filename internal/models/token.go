package models

import "time"

type Token struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	CSRFToken string    `json:"csrf_token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    int       `json:"user_id"`
}
