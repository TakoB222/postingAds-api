package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
}

func NewRepositories(db *sqlx.DB) *Repository {
	return &Repository{}
}
