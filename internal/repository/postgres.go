package repository

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func NewPostgresDB(cfg string) (*sqlx.DB, error) {

	db, err := sqlx.Open("pgx", cfg)
	if err != nil {
		return nil, err
	} else {
		log.Info("Database connected: ", cfg)
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
