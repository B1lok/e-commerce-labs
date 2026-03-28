package postgres

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	"weather-api/internal/config"
)

func RunMigrations(cfg config.Config) {
	escapedPassword := url.QueryEscape(cfg.DBPassword)
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, escapedPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	m, err := migrate.New("file://migrations", connectionString)
	if err != nil {
		slog.Error("Migration initialization failed", "error", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Migration failed", "error", err)
	}
}

func RunMigrationsWithPath(cfg config.Config, migrationPath string) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	m, err := migrate.New(migrationPath, connectionString)
	if err != nil {
		slog.Error("Migration initialization failed", "error", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("Migration failed", "error", err)
	}
}
