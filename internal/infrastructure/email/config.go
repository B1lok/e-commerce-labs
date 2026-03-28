package email

import (
	"weather-api/internal/config"
)

func CreateConfig(cfg config.Config) EmailConfig {
	return EmailConfig{
		Host:     cfg.EmailHost,
		Port:     cfg.EmailPort,
		Username: cfg.EmailUser,
		Password: cfg.EmailPassword,
	}
}
