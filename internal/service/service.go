package service

import "github.com/olenka-91/shorturl/internal/storage/repository"

type urlStoreTranslation interface {
	ShortURL(longURL []byte) (string, error)
	LongURL(shortURL string) (string, error)
	PingDB() error
}

type Service struct {
	urlStoreTranslation
}

func NewService(sbaseURL string, r *repository.Repository) *Service {
	return &Service{urlStoreTranslation: NewUrlStoreTranslationService(sbaseURL, r.UrlStore)}
}
