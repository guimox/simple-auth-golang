package config

import (
	"github.com/guimox/simple-auth-golang/db"
	"github.com/guimox/simple-auth-golang/internal/repository"
)

type Repositories struct {
	UserRepo  *repository.UserRepository
	TokenRepo *repository.TokenRepository
}

func InitializeRepositories() Repositories {
	return Repositories{
		UserRepo:  repository.NewUserRepository(db.DB),
		TokenRepo: repository.NewTokenRepository(db.DB),
	}
}
