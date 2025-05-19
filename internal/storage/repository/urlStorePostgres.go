package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/olenka-91/shorturl/internal/models"
)

type UrlStorePostgres struct {
	db *sqlx.DB
}

func NewUrlStorePostgres(db *sqlx.DB) *UrlStorePostgres {
	return &UrlStorePostgres{db: db}
}

func (u *UrlStorePostgres) PingDB() error {
	return u.db.Ping()
}

func (u *UrlStorePostgres) CloseDB() error {
	return u.db.Close()
}

func (u *UrlStorePostgres) SaveShortURL(ctx context.Context, shortURL, originalURL string, userID int) error {
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	var id int
	row := tx.QueryRowContext(ctx, "INSERT into urls (short, long, user_id) values ($1, $2, $3) ON CONFLICT (long) DO NOTHING RETURNING id", shortURL, originalURL, strconv.Itoa(userID))
	err = row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.NewDBError(originalURL, err)
		} else {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

func (u *UrlStorePostgres) GetOriginalURL(ctx context.Context, shortURL string, userID int) (string, error) {
	row := u.db.QueryRowContext(ctx, "SELECT long FROM urls WHERE short=$1", shortURL)
	var originalUrl string
	err := row.Scan(&originalUrl)
	if err != nil {
		return "", err
	}

	return originalUrl, nil
}

func (u *UrlStorePostgres) PostURLBatch(ctx context.Context, batch []models.BatchForPost, userID int) ([]models.BatchOutput, error) {
	var outputBatch []models.BatchOutput

	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO urls (short, long, user_id) 
		VALUES ($1, $2, $3)	
		ON CONFLICT (long) DO NOTHING 
		RETURNING id
    	`)
	//
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, input := range batch {
		var id int
		err := stmt.QueryRowContext(ctx, input.ShortURL, input.OriginalURL, strconv.Itoa(userID)).Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			} else {
				tx.Rollback()
				return nil, err
			}
		}

		outputBatch = append(outputBatch, models.BatchOutput{
			CorrelationID: input.CorrelationID,
			ShortURL:      input.ShortURL,
		})

	}

	return outputBatch, tx.Commit()
}

func (u *UrlStorePostgres) ListURLsByUser(ctx context.Context, userID int) ([]models.URLsForUser, error) {

	var URLs []models.URLsForUser
	err := u.db.SelectContext(ctx, &URLs, "SELECT long, short FROM urls WHERE user_ID=$1", strconv.Itoa(userID))
	if err != nil {
		return nil, err
	}

	return URLs, nil
}
