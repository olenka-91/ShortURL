package storage

import (
	"github.com/olenka-91/shorturl/config"
	filestorage "github.com/olenka-91/shorturl/internal/storage/FileStorage"
	inmemorystorage "github.com/olenka-91/shorturl/internal/storage/InMemoryStorage"
)

type Storage interface {
	SaveShortURL(shortURL, originalURL string) error
	GetOriginalURL(shortURL string) (string, error)
}

func NewStorage(cfg config.MyConfigs) (Storage, error) {
	// Псевдокод дальше
	if cfg.DatabaseDSN != "" {
		// Создаём хранилище на базе БД
		return NewDBStorage(cfg.DatabaseDSN), nil
	}

	// Проверяем наличие переменной/флага для файла
	if cfg.FilePath != "" {
		return filestorage.NewFileStorage(cfg.FilePath), nil
	}

	// Всё остальное — используем хранилище в памяти
	return inmemorystorage.NewInMemoryStorage(), nil
}
