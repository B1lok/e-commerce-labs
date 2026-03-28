package services

import (
	"context"
	"weather-api/internal/application/mapper"
	"weather-api/internal/application/query"
	"weather-api/internal/domain/repositories"
)

type WeatherService struct {
	weatherRepository repositories.WeatherRepository
}

func NewWeatherService(weatherRepository repositories.WeatherRepository) *WeatherService {
	return &WeatherService{weatherRepository: weatherRepository}
}

func (s *WeatherService) GetWeather(ctx context.Context, city string) (*query.WeatherQueryResult, error) {
	weather, err := s.weatherRepository.GetWeather(ctx, city)
	if err != nil {
		return nil, err
	}
	var queryResult query.WeatherQueryResult
	queryResult.Result = mapper.NewWeatherResult(weather)
	return &queryResult, nil
}
