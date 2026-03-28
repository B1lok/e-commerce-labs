package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"weather-api/internal/domain/models"
)

type MockWeatherRepository struct {
	mock.Mock
}

func (m *MockWeatherRepository) GetWeather(ctx context.Context, city string) (*models.Weather, error) {
	args := m.Called(ctx, city)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Weather), args.Error(1)
}

func (m *MockWeatherRepository) GetDailyForecast(ctx context.Context, city string) (*models.WeatherDaily, error) {
	args := m.Called(ctx, city)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WeatherDaily), args.Error(1)
}

func (m *MockWeatherRepository) GetHourlyForecast(ctx context.Context, city string) (*models.WeatherHourly, error) {
	args := m.Called(ctx, city)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WeatherHourly), args.Error(1)
}
