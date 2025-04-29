package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"

	"github.com/olenka-91/shorturl/internal/models"
	"github.com/olenka-91/shorturl/internal/storage"
)

type urlStoreTranslationService struct {
	urlStore storage.Storage
	baseURL  string
}

func NewUrlStoreTranslationService(sbaseURL string, st storage.Storage) *urlStoreTranslationService {
	serv := urlStoreTranslationService{urlStore: st, baseURL: sbaseURL}
	return &serv
}

func (r *urlStoreTranslationService) generateShortUrl() string {
	bytes := make([]byte, 6) // 6 байт = 8 символов в base64
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	// Кодируем байты в строку base64
	return base64.URLEncoding.EncodeToString(bytes)
}

func (r *urlStoreTranslationService) ShortURL(ctx context.Context, longURL []byte) (string, error) {
	shortURL64 := r.generateShortUrl()
	return r.baseURL + shortURL64, r.urlStore.SaveShortURL(ctx, shortURL64, string(longURL))

}

func (r *urlStoreTranslationService) LongURL(ctx context.Context, shortURL string) (string, error) {
	return r.urlStore.GetOriginalURL(ctx, shortURL)
}

func (r *urlStoreTranslationService) PingDB() error {
	if st, ok := r.urlStore.(storage.DBStorage); ok {
		return st.PingDB()
	}
	return errors.New("Database operations not supported by this storage")
}

func (r *urlStoreTranslationService) CloseDB() error {
	if st, ok := r.urlStore.(storage.DBStorage); ok {
		return st.CloseDB()
	}
	return nil
}

func (r *urlStoreTranslationService) PostURLBatch(ctx context.Context, batch models.ArrBatchInput) ([]models.BatchOutput, error) {
	b := batch.Validate()
	if b == false {
		return nil, errors.New("Не найдено входных значений")
	}

	fullBatch := make([]models.BatchForPost, len(batch))

	for i, input := range batch {
		fullBatch[i].CorrelationID = input.CorrelationID
		fullBatch[i].OriginalURL = input.OriginalURL
		fullBatch[i].ShortURL = r.generateShortUrl()
	}
	outputBatch, err := r.urlStore.PostURLBatch(ctx, fullBatch)
	if err != nil {
		return nil, err
	}

	return outputBatch, nil

}
