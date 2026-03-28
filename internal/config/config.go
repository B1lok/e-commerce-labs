package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	ServerHost    string `mapstructure:"SERVER_HOST"`
	ServerPort    string `mapstructure:"SERVER_PORT"`
	WeatherApiKey string `mapstructure:"WEATHER_API_KEY"`
	EmailHost     string `mapstructure:"EMAIL_HOST"`
	EmailPort     int    `mapstructure:"EMAIL_PORT"`
	EmailUser     string `mapstructure:"EMAIL_USERNAME"`
	EmailPassword string `mapstructure:"EMAIL_PASSWORD"`
}

func LoadConfig() (Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Print("No .env file found, binding individual environment variables.")
		bindAllEnvVars()
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

func bindAllEnvVars() {
	_ = viper.BindEnv("DB_HOST")
	_ = viper.BindEnv("DB_PORT")
	_ = viper.BindEnv("DB_USER")
	_ = viper.BindEnv("DB_PASSWORD")
	_ = viper.BindEnv("DB_NAME")
	_ = viper.BindEnv("SERVER_HOST")
	_ = viper.BindEnv("SERVER_PORT")
	_ = viper.BindEnv("WEATHER_API_KEY")
	_ = viper.BindEnv("EMAIL_HOST")
	_ = viper.BindEnv("EMAIL_PORT")
	_ = viper.BindEnv("EMAIL_USERNAME")
	_ = viper.BindEnv("EMAIL_PASSWORD")
}
