package postgres

import (
	"MentorApiProject/internal/config"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
	"log"
)

func Connect(cfg config.PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.Dsn())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(cfg config.PostgresConfig) error {
	m, err := migrate.New("file://migrations", cfg.Dsn())
	if err != nil {
		return err
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("POSTGRES MIGRATIONS: No changes detected.")
		return nil
	}
	if err != nil {
		return err
	}

	log.Println("POSTGRES MIGRATIONS: Migrations completed successfully.")
	return nil
}
