package postgres

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"net/url"
	"weather-api/internal/config"
)

func RunMigrations(cfg config.Config) {
	escapedPassword := url.QueryEscape(cfg.DBPassword)
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, escapedPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	m, err := migrate.New("file://migrations", connectionString)
	if err != nil {
		log.Println("Migration initialization failed")
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Migration failed: %v", err)
	}
}

func RunMigrationsWithPath(cfg config.Config, migrationPath string) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	m, err := migrate.New(migrationPath, connectionString)
	if err != nil {
		log.Printf("Migration initialization failed: %v", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Migration failed: %v", err)
	}
}
