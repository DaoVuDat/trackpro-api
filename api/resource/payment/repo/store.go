package paymentrepo

import "database/sql"

type postgresStore struct {
	db *sql.DB
}

func NewPostgresStore(store *sql.DB) *postgresStore {
	return &postgresStore{
		db: store,
	}
}
