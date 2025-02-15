package repository

import (
	"database/sql"

	"github.com/guimox/simple-auth-golang/internal/models"
)

type TokenRepository struct {
	DB *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{DB: db}
}

func (r *TokenRepository) CreateToken(token *models.Token) error {
	query := `
		INSERT INTO tokens (token, csrf_token, expires_at, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	err := r.DB.QueryRow(query, token.Token, token.CSRFToken, token.ExpiresAt, token.UserID).Scan(&token.ID)
	return err
}

func (r *TokenRepository) GetToken(token string) (*models.Token, error) {
	t := &models.Token{}
	query := `
		SELECT id, token, csrf_token, expires_at, user_id
		FROM tokens
		WHERE token = $1`
	err := r.DB.QueryRow(query, token).Scan(&t.ID, &t.Token, &t.CSRFToken, &t.ExpiresAt, &t.UserID)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TokenRepository) DeleteToken(token string) error {
	query := `DELETE FROM tokens WHERE token = $1`
	_, err := r.DB.Exec(query, token)
	return err
}
