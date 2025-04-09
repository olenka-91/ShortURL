package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/olenka-91/shorturl/internal/repository"
	"github.com/olenka-91/shorturl/internal/storage"
)

type urlStoreTranslationService struct {
	urlStore *storage.Storage
	repo     repository.UrlStore

	baseURL string
}

func NewUrlStoreTranslationService(sbaseURL string, r repository.UrlStore) *urlStoreTranslationService {
	st := urlStoreTranslationService{urlStore: storage.NewStorage(), baseURL: sbaseURL, repo: r}
	//	st.urlStore.LoadFromFile(config.MyConfigs.FileName)
	return &st
}

func (r *urlStoreTranslationService) ShortURL(longURL []byte) (string, error) {
	bytes := make([]byte, 6) // 6 байт = 8 символов в base64
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	// Кодируем байты в строку base64

	shortURL64 := base64.URLEncoding.EncodeToString(bytes)
	r.urlStore.Add(shortURL64, string(longURL))

	fmt.Println("longURL=", string(longURL))
	fmt.Println("shortURL=", shortURL64)
	return r.baseURL + shortURL64, nil

}

func (r *urlStoreTranslationService) LongURL(shortURL string) (string, error) {
	return r.urlStore.Get(shortURL)
}

func (r *urlStoreTranslationService) PingDB() error {
	return r.repo.PingDB()
}
