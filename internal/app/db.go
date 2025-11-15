package app

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/milyrock/PR-Reviewer/internal/config"
)

func InitDB(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf(
		"host=%s user=%s password=%s database=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.Username,
		cfg.Password,
		cfg.Database,
		cfg.Port,
	)

	db, err := sqlx.Open("pgx", dataSource)
	if err != nil {
		return nil, fmt.Errorf("create pool of connections to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	migrationFile := "./db/pr.sql"
	if _, err := os.Stat(migrationFile); err == nil {
		if err := initSchema(db, migrationFile); err != nil {
			return nil, fmt.Errorf("failed to initialize schema: %w", err)
		}
	}

	return db, nil
}

func initSchema(db *sqlx.DB, migrationFile string) error {
	sqlBytes, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}
