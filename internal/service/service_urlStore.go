package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"sync"
)

type urlStoreTranslationService struct {
	/*repo repository.Song*/
	urlStore map[string]string
	mu       sync.Mutex
	baseURL  string
}

func NewUrlStoreTranslationService(sbaseURL string /*r repository.Song*/) *urlStoreTranslationService {
	return &urlStoreTranslationService{urlStore: make(map[string]string), baseURL: sbaseURL /*repo: r*/}
}

func (r *urlStoreTranslationService) ShortURL(longURL []byte) (string, error) {
	bytes := make([]byte, 6) // 6 байт = 8 символов в base64
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	// Кодируем байты в строку base64

	shortURL64 := base64.URLEncoding.EncodeToString(bytes)

	r.mu.Lock()
	defer r.mu.Unlock()
	r.urlStore[shortURL64] = string(longURL)

	fmt.Println("longURL=", string(longURL))
	fmt.Println("shortURL=", shortURL64)
	return r.baseURL + shortURL64, nil

}

func (r *urlStoreTranslationService) LongURL(shortURL string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	longURL, ok := r.urlStore[shortURL]
	if !ok {
		return "", fmt.Errorf("Not found")
	}
	fmt.Println("shortURL=", shortURL)
	fmt.Println("longURL=", longURL)
	return longURL, nil
}
