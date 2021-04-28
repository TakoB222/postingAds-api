package repository

import (
	"github.com/TakoB222/postingAds-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user domain.User)(int, error)
	GetUser(email, password_hash string)(domain.User, error)
	GetSessionByRefreshToken(refreshToken string)(domain.Session, error)
	DeleteSessionByUserId(userId string) error
	SetSession(session domain.Session) error
}

type Admin interface {
}

type Ad interface {

}

type Repository struct {
	Authorization
	Admin
	Ad
}

func NewRepositories(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
	}
}
