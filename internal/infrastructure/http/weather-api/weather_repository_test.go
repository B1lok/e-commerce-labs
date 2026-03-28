package weather_api

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func mockHTTPClient(response string, statusCode int) *http.Client {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	}))

	return &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(server.URL)
			},
		},
	}
}

type MockClock struct{}

func (MockClock) Now() time.Time {
	t, _ := time.Parse("2006-01-02 15:04", "2025-05-18 00:00")
	return t
}

func loadJSONFile(t *testing.T, filename string) string {
	path := filepath.Join("testdata", filename)
	content, err := os.ReadFile(path)
	require.NoError(t, err, "Failed to read test data file: %s", path)
	return string(content)
}

func TestGetWeather(t *testing.T) {
	mockResponse := loadJSONFile(t, "current_weather_response.json")

	repo := NewWeatherRepository("dummy-api-key")
	repo.client = mockHTTPClient(mockResponse, http.StatusOK)

	weather, err := repo.GetWeather(context.Background(), "London")

	require.NoError(t, err)
	require.NotNil(t, weather)
	assert.Equal(t, 11.0, weather.Temperature)
	assert.Equal(t, "Partly cloudy", weather.Description)
}

func TestGetDailyForecast(t *testing.T) {
	mockResponse := loadJSONFile(t, "forecast_response.json")

	repo := NewWeatherRepository("dummy-api-key")
	repo.client = mockHTTPClient(mockResponse, http.StatusOK)

	forecast, err := repo.GetDailyForecast(context.Background(), "London")

	require.NoError(t, err)
	require.NotNil(t, forecast)

	assert.Equal(t, "London", forecast.Location)
	assert.Equal(t, "Partly Cloudy", forecast.Condition)
	assert.Equal(t, 12.3, forecast.AvgTempC)
}

func TestGetHourlyForecast(t *testing.T) {
	mockResponse := loadJSONFile(t, "forecast_response.json")

	repo := NewWeatherRepository("dummy-api-key")
	repo.client = mockHTTPClient(mockResponse, http.StatusOK)
	repo.SetClock(MockClock{})

	forecast, err := repo.GetHourlyForecast(context.Background(), "London")

	require.NoError(t, err)
	require.NotNil(t, forecast)

	assert.Equal(t, "London", forecast.Location)
	assert.Equal(t, "Clear", forecast.Condition)
	assert.Equal(t, 10.2, forecast.TempC)
}

func TestHandleAPIErrorResponse(t *testing.T) {
	mockErrorResponse := loadJSONFile(t, "not_found_response.json")

	repo := NewWeatherRepository("dummy-api-key")
	repo.client = mockHTTPClient(mockErrorResponse, http.StatusBadRequest)

	_, err := repo.GetWeather(context.Background(), "NonExistentCity")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "City not found")
}

func TestDirectMapping(t *testing.T) {
	currentJSON := loadJSONFile(t, "current_weather_response.json")
	dailyJSON := loadJSONFile(t, "forecast_response.json")
	hourlyJSON := loadJSONFile(t, "forecast_response.json")

	var currentResponse WeatherRepositoryResponse
	var dailyResponse WeatherDailyResponse
	var hourlyResponse WeatherHourlyResponse

	err := json.Unmarshal([]byte(currentJSON), &currentResponse)
	require.NoError(t, err)

	err = json.Unmarshal([]byte(dailyJSON), &dailyResponse)
	require.NoError(t, err)

	err = json.Unmarshal([]byte(hourlyJSON), &hourlyResponse)
	require.NoError(t, err)

	weather := ToWeather(&currentResponse)
	dailyForecast := ToWeatherDaily(&dailyResponse)
	hourlyForecast := ToWeatherHourly(&hourlyResponse, MockClock{}.Now())

	assert.NotNil(t, weather)
	assert.NotNil(t, dailyForecast)
	assert.NotNil(t, hourlyForecast)
}
