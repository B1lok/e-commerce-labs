package weather_api

import (
	"time"
	"weather-api/internal/domain/models"
)

func ToWeather(weatherRepositoryResponse *WeatherRepositoryResponse) *models.Weather {
	return &models.Weather{
		Temperature: weatherRepositoryResponse.Current.TempC,
		Humidity:    weatherRepositoryResponse.Current.Humidity,
		Description: weatherRepositoryResponse.Current.Condition.Text,
	}
}

func ToWeatherDaily(response *WeatherDailyResponse) *models.WeatherDaily {
	first := response.Forecast.Forecastday[0]

	return &models.WeatherDaily{
		Location:   response.Location.Name,
		Date:       first.Date,
		MaxTempC:   first.Day.MaxtempC,
		MinTempC:   first.Day.MintempC,
		AvgTempC:   first.Day.AvgtempC,
		WillItRain: first.Day.DailyWillItRain == 1,
		ChanceRain: first.Day.DailyChanceOfRain,
		WillItSnow: first.Day.DailyWillItSnow == 1,
		ChanceSnow: first.Day.DailyChanceOfSnow,
		Condition:  first.Day.Condition.Text,
		Icon:       first.Day.Condition.Icon,
	}
}

func ToWeatherHourly(response *WeatherHourlyResponse, targetTime time.Time) *models.WeatherHourly {
	currentTime := targetTime.Truncate(time.Hour).Format("2006-01-02 15:04")
	for _, hour := range response.Forecast.Forecastday[0].Hour {
		if hour.Time == currentTime {
			return &models.WeatherHourly{
				Location:   response.Location.Name,
				Time:       hour.Time,
				TempC:      hour.TempC,
				WillItRain: hour.WillItRain == 1,
				ChanceRain: hour.ChanceOfRain,
				WillItSnow: hour.WillItSnow == 1,
				ChanceSnow: hour.ChanceOfSnow,
				Condition:  hour.Condition.Text,
				Icon:       hour.Condition.Icon,
			}
		}
	}

	return nil
}
