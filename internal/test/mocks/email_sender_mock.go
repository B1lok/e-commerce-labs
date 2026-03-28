package mocks

import (
	"github.com/stretchr/testify/mock"
	"weather-api/internal/application/email"
)

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) ConfirmationEmail(mail *email.ConfirmationEmail) error {
	args := m.Called(mail)
	return args.Error(0)
}

func (m *MockEmailSender) WeatherDailyEmail(mail *email.WeatherDailyEmail) error {
	args := m.Called(mail)
	return args.Error(0)
}

func (m *MockEmailSender) WeatherHourlyEmail(mail *email.WeatherHourlyEmail) error {
	args := m.Called(mail)
	return args.Error(0)
}
