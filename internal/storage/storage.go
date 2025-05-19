package storage

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/olenka-91/shorturl/config"
	"github.com/olenka-91/shorturl/internal/models"
	filestorage "github.com/olenka-91/shorturl/internal/storage/FileStorage"
	inmemorystorage "github.com/olenka-91/shorturl/internal/storage/InMemoryStorage"
	"github.com/olenka-91/shorturl/internal/storage/repository"
)

type Storage interface {
	SaveShortURL(ctx context.Context, shortURL, originalURL string, userID int) error
	GetOriginalURL(ctx context.Context, shortURL string, userID int) (string, error)
	PostURLBatch(ctx context.Context, batch []models.BatchForPost, userID int) ([]models.BatchOutput, error)
}

type DBStorage interface {
	Storage
	PingDB() error
	CloseDB() error
	ListURLsByUser(ctx context.Context, userID int) ([]models.URLsForUser, error)
}

func NewStorage(cfg config.MyConfigs) (Storage, error) {
	// Псевдокод дальше
	if cfg.DBDSN != "" {
		// Создаём хранилище на базе БД
		db, err := repository.NewPostgresDB(cfg.DBDSN)
		if err != nil {
			log.Errorf("error occured while connecting to DB: %s", err.Error())
			return nil, err
		}
		log.Debug("Created repository storage")
		return repository.NewRepository(db), nil
	}

	// Проверяем наличие переменной/флага для файла
	if cfg.FilePath != "" {
		log.Debug("Created file storage")
		return filestorage.NewFileStorage(cfg.FilePath), nil
	}

	// Всё остальное — используем хранилище в памяти
	log.Debug("Created in memory storage")
	return inmemorystorage.NewInMemoryStorage(), nil
}
