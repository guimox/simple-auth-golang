package repository

import (
	"database/sql"
	"errors"

	"github.com/guimox/simple-auth-golang/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, hashed_password) VALUES ($1, $2) RETURNING id`
	err := r.DB.QueryRow(query, user.Username, user.HashedPassword).Scan(&user.ID)
	if err != nil {
		return errors.New("user already exists")
	}
	return nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, hashed_password FROM users WHERE username = $1`
	err := r.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.HashedPassword)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
