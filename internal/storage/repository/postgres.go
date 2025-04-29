package repository

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func runMigrations(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	// Создаем экземпляр миграции
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../schema", // Путь к миграциям на диске
		"postgres",            // Тип базы данных (например, postgres)
		driver,                // Драйвер базы данных
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Выполняем миграции
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Migrations don't needed: no changes.")
		} else {
			log.Fatalf("Maigrations failed: %v", err)
		}
	} else {
		log.Println("Migrations applied successfully")
	}

	return nil
}

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

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}
