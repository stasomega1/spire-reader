package store

import (
	"github.com/jmoiron/sqlx"
	pgdb "spire-reader/app/store/pgdb"
)

type Store struct {
	db           *sqlx.DB
	dbRepository *pgdb.Repository
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) PgdbRepository() *pgdb.Repository {
	if store.dbRepository != nil {
		return store.dbRepository
	}
	store.dbRepository = pgdb.NewRepository(store.db)
	return store.dbRepository
}
