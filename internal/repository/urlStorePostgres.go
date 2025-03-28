package repository

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
