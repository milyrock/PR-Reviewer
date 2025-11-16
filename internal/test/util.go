package test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// SetupTestDB создает тестовую БД используя testcontainers
func SetupTestDB(ctx context.Context) (*sqlx.DB, func(), error) {
	// Используем контекст с увеличенным таймаутом
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	postgresContainer, err := postgres.Run(ctxWithTimeout,
		"postgres:15",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForSQL("5432/tcp", "pgx", func(host string, port nat.Port) string {
				return fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable", host, port.Port())
			}).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Применяем миграции
	if err := applyMigrations(db); err != nil {
		return nil, nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	cleanup := func() {
		db.Close()
		postgresContainer.Terminate(ctx)
	}

	return db, cleanup, nil
}

// applyMigrations применяет SQL миграции из db/pr.sql
func applyMigrations(db *sqlx.DB) error {
	// Получаем путь к файлу миграции относительно корня проекта
	migrationPath := "db/pr.sql"

	// Если файл не найден, пробуем найти его относительно текущей директории
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		// Пробуем найти файл, поднимаясь вверх по директориям
		wd, _ := os.Getwd()
		for i := 0; i < 5; i++ {
			testPath := filepath.Join(wd, migrationPath)
			if _, err := os.Stat(testPath); err == nil {
				migrationPath = testPath
				break
			}
			wd = filepath.Dir(wd)
		}
	}

	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}
