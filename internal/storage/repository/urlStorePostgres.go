package storage

import (
	"github.com/jmoiron/sqlx"
)

type UrlStorePostgres struct {
	db *sqlx.DB
}

func NewUrlStorePostgres(db *sqlx.DB) *UrlStorePostgres {
	return &UrlStorePostgres{db: db}
}

func (u *UrlStorePostgres) PingDB() error {
	return u.db.Ping()
}

func (u *UrlStorePostgres) SaveShortURL(shortURL, originalURL string) error {
	return nil
}

func (u *UrlStorePostgres) GetOriginalURL(shortURL string) (string, error) {
	return "", nil
}
