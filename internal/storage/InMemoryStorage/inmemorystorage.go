package inmemorystorage

import (
	"fmt"
	"sync"
)

type InMemoryStorage struct {
	urlStorage map[string]string
	mu         sync.Mutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{urlStorage: make(map[string]string)}
}

func (s *InMemoryStorage) SaveShortURL(shortURL64, longURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urlStorage[shortURL64] = string(longURL)
	return nil
}

func (s *InMemoryStorage) GetOriginalURL(shortURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	longURL, ok := s.urlStorage[shortURL]
	if !ok {
		return "", fmt.Errorf("Not found")
	}
	return longURL, nil
}
