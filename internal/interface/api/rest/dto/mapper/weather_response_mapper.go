package mapper

import (
	"weather-api/internal/application/common"
	"weather-api/internal/interface/api/rest/dto/response"
)

func ToWeatherResponse(weatherResult *common.WeatherResult) *response.WeatherResponse {
	return &response.WeatherResponse{
		Temperature: weatherResult.Temperature,
		Humidity:    weatherResult.Humidity,
		Description: weatherResult.Description,
	}
}
