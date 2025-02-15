package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/guimox/simple-auth-golang/internal/models"
	"github.com/guimox/simple-auth-golang/internal/repository"
	"github.com/guimox/simple-auth-golang/internal/utils"
)

type AuthHandler struct {
	UserRepo  *repository.UserRepository
	TokenRepo *repository.TokenRepository
}

func NewAuthHandler(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo, TokenRepo: tokenRepo}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) < 4 || len(password) < 8 {
		http.Error(w, "Invalid username/password", http.StatusNotAcceptable)
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Username:       username,
		HashedPassword: hashedPassword,
	}

	err = h.UserRepo.CreateUser(user)
	if err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.UserRepo.GetUserByUsername(username)
	if err != nil || !utils.CheckPasswordHash(password, user.HashedPassword) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionToken, err := utils.GenerateToken(32)
	if err != nil {
		http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
		return
	}

	csrfToken, err := utils.GenerateToken(32)
	if err != nil {
		http.Error(w, "Failed to generate CSRF token", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	token := &models.Token{
		Token:     sessionToken,
		CSRFToken: csrfToken,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
	}

	err = h.TokenRepo.CreateToken(token)
	if err != nil {
		http.Error(w, "Failed to store tokens", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  expiresAt,
		HttpOnly: false,
	})

	fmt.Fprintln(w, "Login successful!")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "No session token found", http.StatusBadRequest)
		return
	}

	err = h.TokenRepo.DeleteToken(sessionCookie.Value)
	if err != nil {
		http.Error(w, "Failed to delete session token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	fmt.Fprintln(w, "Logout successful!")
}
