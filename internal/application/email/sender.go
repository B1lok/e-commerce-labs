package email

type Sender interface {
	ConfirmationEmail(email *ConfirmationEmail) error
	WeatherDailyEmail(email *WeatherDailyEmail) error
	WeatherHourlyEmail(email *WeatherHourlyEmail) error
}
