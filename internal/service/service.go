package service

import (
	"context"

	"github.com/olenka-91/shorturl/internal/models"
	"github.com/olenka-91/shorturl/internal/storage"
)

type urlStoreTranslation interface {
	ShortURL(ctx context.Context, longURL []byte, userID int) (string, error)
	LongURL(ctx context.Context, shortURL string, userID int) (string, error)
	PostURLBatch(ctx context.Context, batch models.ArrBatchInput, userID int) ([]models.BatchOutput, error)
	PingDB() error
	CloseDB() error
	ListURLsByUser(ctx context.Context, userID int) ([]models.URLsForUser, error)
}

type Authorization interface {
	GenerateToken(userID int) (string, error)
	ParseToken(accessToken string) (int, error)
}

type Service struct {
	urlStoreTranslation
	Authorization
}

func NewService(sbaseURL string, st storage.Storage) *Service {
	return &Service{
		urlStoreTranslation: NewUrlStoreTranslationService(sbaseURL, st),
		Authorization:       NewAuthorizationService(),
	}
}
