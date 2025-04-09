package storage

import "github.com/jmoiron/sqlx"

type UrlStore interface {
	PingDB() error
	SaveShortURL(shortURL, originalURL string) error
	GetOriginalURL(shortURL string) (string, error)
}

type Repository struct {
	UrlStore
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{UrlStore: NewUrlStorePostgres(db)}
}
