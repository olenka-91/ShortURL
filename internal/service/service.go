package service

import (
	"context"

	"github.com/olenka-91/shorturl/internal/models"
	"github.com/olenka-91/shorturl/internal/storage"
)

type urlStoreTranslation interface {
	ShortURL(ctx context.Context, longURL []byte) (string, error)
	LongURL(ctx context.Context, shortURL string) (string, error)
	PostURLBatch(ctx context.Context, batch models.ArrBatchInput) ([]models.BatchOutput, error)
	PingDB() error
	CloseDB() error
}

type Service struct {
	urlStoreTranslation
}

func NewService(sbaseURL string, st storage.Storage) *Service {
	return &Service{urlStoreTranslation: NewUrlStoreTranslationService(sbaseURL, st)}
}
