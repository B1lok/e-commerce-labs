package mapper

import (
	"weather-api/internal/application/common"
	"weather-api/internal/domain/models"
)

func NewWeatherResult(weather *models.Weather) *common.WeatherResult {
	return &common.WeatherResult{
		Temperature: weather.Temperature,
		Humidity:    weather.Humidity,
		Description: weather.Description,
	}
}
