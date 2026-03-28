package stubs

import (
	"context"
	"net/http"
	"weather-api/internal/domain/models"
	"weather-api/pkg/errors"
)

type WeatherRepositoryStub struct {
	GetWeatherFn        func(ctx context.Context, city string) (*models.Weather, error)
	GetDailyForecastFn  func(ctx context.Context, city string) (*models.WeatherDaily, error)
	GetHourlyForecastFn func(ctx context.Context, city string) (*models.WeatherHourly, error)
}

func NewWeatherRepositoryStub() *WeatherRepositoryStub {
	return &WeatherRepositoryStub{
		GetWeatherFn:        nil,
		GetDailyForecastFn:  nil,
		GetHourlyForecastFn: nil,
	}
}

func (s *WeatherRepositoryStub) GetWeather(ctx context.Context, city string) (*models.Weather, error) {
	if s.GetWeatherFn != nil {
		return s.GetWeatherFn(ctx, city)
	}
	if city == "InvalidCity" {
		return nil, errors.New("City not found", http.StatusNotFound)
	}
	return &models.Weather{
		Temperature: 20.5,
		Humidity:    60.0,
		Description: "Clear sky",
	}, nil
}

func (s *WeatherRepositoryStub) GetDailyForecast(ctx context.Context, city string) (*models.WeatherDaily, error) {
	if s.GetDailyForecastFn != nil {
		return s.GetDailyForecastFn(ctx, city)
	}
	return &models.WeatherDaily{
		Location:   city,
		Date:       "2025-05-18",
		MaxTempC:   22.0,
		MinTempC:   15.0,
		AvgTempC:   18.5,
		WillItRain: false,
		ChanceRain: 10,
		WillItSnow: false,
		ChanceSnow: 0,
		Condition:  "Sunny",
		Icon:       "sunny.png",
	}, nil
}

func (s *WeatherRepositoryStub) GetHourlyForecast(ctx context.Context, city string) (*models.WeatherHourly, error) {
	if s.GetHourlyForecastFn != nil {
		return s.GetHourlyForecastFn(ctx, city)
	}
	return &models.WeatherHourly{
		Location:   city,
		Time:       "12:00",
		TempC:      20.0,
		WillItRain: false,
		ChanceRain: 5,
		WillItSnow: false,
		ChanceSnow: 0,
		Condition:  "Partly cloudy",
		Icon:       "cloudy.png",
	}, nil
}
