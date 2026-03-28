package stubs

import "weather-api/internal/application/email"

type SenderStub struct {
	ConfirmationEmailFn  func(email *email.ConfirmationEmail) error
	WeatherDailyEmailFn  func(email *email.WeatherDailyEmail) error
	WeatherHourlyEmailFn func(email *email.WeatherHourlyEmail) error
}

func NewSenderStub() *SenderStub {
	return &SenderStub{
		ConfirmationEmailFn:  nil,
		WeatherDailyEmailFn:  nil,
		WeatherHourlyEmailFn: nil,
	}
}

func (s *SenderStub) ConfirmationEmail(email *email.ConfirmationEmail) error {
	if s.ConfirmationEmailFn != nil {
		return s.ConfirmationEmailFn(email)
	}
	return nil
}

func (s *SenderStub) WeatherDailyEmail(email *email.WeatherDailyEmail) error {
	if s.WeatherDailyEmailFn != nil {
		return s.WeatherDailyEmailFn(email)
	}
	return nil
}

func (s *SenderStub) WeatherHourlyEmail(email *email.WeatherHourlyEmail) error {
	if s.WeatherHourlyEmailFn != nil {
		return s.WeatherHourlyEmailFn(email)
	}
	return nil
}
