package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"weather-api/internal/domain/models"
	"weather-api/internal/test/mocks"
	"weather-api/pkg/errors"
)

func TestWeatherService_GetWeather_Success(t *testing.T) {
	mockRepo := new(mocks.MockWeatherRepository)
	service := NewWeatherService(mockRepo)
	ctx := context.Background()
	city := "Berlin"

	mockWeather := &models.Weather{
		Temperature: 20.5,
		Humidity:    70,
		Description: "Partly Cloudy",
	}

	mockRepo.On("GetWeather", ctx, city).Return(mockWeather, nil)

	result, err := service.GetWeather(ctx, city)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockWeather.Temperature, result.Result.Temperature)
	assert.Equal(t, mockWeather.Humidity, result.Result.Humidity)
	assert.Equal(t, mockWeather.Description, result.Result.Description)
	mockRepo.AssertExpectations(t)
}

func TestWeatherService_GetWeather_EmptyCity(t *testing.T) {
	mockRepo := new(mocks.MockWeatherRepository)
	service := NewWeatherService(mockRepo)
	ctx := context.Background()
	city := "invalid city123"
	expectedErr := errors.New("invalid city123", http.StatusNotFound)

	mockRepo.On("GetWeather", ctx, city).Return(nil, expectedErr)

	result, err := service.GetWeather(ctx, city)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
