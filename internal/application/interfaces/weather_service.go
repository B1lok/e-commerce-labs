package interfaces

import (
	"context"
	"weather-api/internal/application/query"
)

type WeatherService interface {
	GetWeather(ctx context.Context, city string) (*query.WeatherQueryResult, error)
}
