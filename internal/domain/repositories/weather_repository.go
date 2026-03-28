package repositories

import (
	"context"
	"weather-api/internal/domain/models"
)

type WeatherRepository interface {
	GetWeather(ctx context.Context, city string) (*models.Weather, error)
	GetDailyForecast(ctx context.Context, city string) (*models.WeatherDaily, error)
	GetHourlyForecast(ctx context.Context, city string) (*models.WeatherHourly, error)
}
