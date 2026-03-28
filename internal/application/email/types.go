package email

import "weather-api/internal/domain/models"

type ConfirmationEmail struct {
	To        string
	City      string
	Frequency string
	Url       string
}

type WeatherDailyEmail struct {
	To             string
	Frequency      string
	UnsubscribeUrl string
	WeatherDaily   *models.WeatherDaily
}

type WeatherHourlyEmail struct {
	To             string
	Frequency      string
	UnsubscribeUrl string
	WeatherHourly  *models.WeatherHourly
}
