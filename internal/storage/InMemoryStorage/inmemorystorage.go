package inmemorystorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/olenka-91/shorturl/internal/models"
)

type InMemoryStorage struct {
	urlStorage map[string]string
	mu         sync.Mutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{urlStorage: make(map[string]string)}
}

func (s *InMemoryStorage) SaveShortURL(ctx context.Context, shortURL64, longURL string, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urlStorage[shortURL64] = string(longURL)
	return nil
}

func (s *InMemoryStorage) GetOriginalURL(ctx context.Context, shortURL string, userID int) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	longURL, ok := s.urlStorage[shortURL]
	if !ok {
		return "", fmt.Errorf("Not found")
	}
	return longURL, nil
}

func (s *InMemoryStorage) PostURLBatch(ctx context.Context, batch []models.BatchForPost, userID int) ([]models.BatchOutput, error) {
	return nil, nil
}
