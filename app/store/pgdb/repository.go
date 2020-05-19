package pgdb

import (
	"github.com/jmoiron/sqlx"
	"spire-reader/app/model"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (repository *Repository) Version() (string, error) {
	return model.Version, nil
}
