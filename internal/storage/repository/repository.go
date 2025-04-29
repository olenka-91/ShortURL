package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/olenka-91/shorturl/internal/models"
)

type UrlStore interface {
	PingDB() error
	CloseDB() error
	SaveShortURL(ctx context.Context, shortURL, originalURL string) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	PostURLBatch(ctx context.Context, batch []models.BatchForPost) ([]models.BatchOutput, error)
}

type Repository struct {
	UrlStore
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{UrlStore: NewUrlStorePostgres(db)}
}
