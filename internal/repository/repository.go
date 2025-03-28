package repository

import "github.com/jmoiron/sqlx"

type UrlStore interface {
	PingDB() error
}

type Repository struct {
	UrlStore
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{UrlStore: NewUrlStorePostgres(db)}
}
