package email

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"net/http"
	"weather-api/internal/application/email"
	"weather-api/pkg/errors"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Sender struct {
	config EmailConfig
	dialer *gomail.Dialer
}

func NewEmailSender(config EmailConfig) *Sender {
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	return &Sender{
		config: config,
		dialer: dialer,
	}
}

func (s *Sender) ConfirmationEmail(email *email.ConfirmationEmail) error {
	return s.sendEmail("templates/confirm.html", email.To, "Confirm your subscription", email)
}

func (s *Sender) WeatherDailyEmail(email *email.WeatherDailyEmail) error {
	return s.sendEmail("templates/daily.html", email.To, "Your weather daily forecast", email)
}

func (s *Sender) WeatherHourlyEmail(email *email.WeatherHourlyEmail) error {
	return s.sendEmail("templates/hourly.html", email.To, "Your weather hourly forecast", email)
}

func (s *Sender) sendEmail(templatePath, to, subject string, data any) error {
	htmlBody, err := renderTemplate(templatePath, data)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", s.config.Username)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", htmlBody)

	if err := s.dialer.DialAndSend(message); err != nil {
		fmt.Println(err.Error())
		return errors.New("Failed to send email", http.StatusInternalServerError)
	}

	return nil
}

func renderTemplate(templatePath string, data any) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Printf("Template parse error: %v\n", err)
		return "", errors.New("Failed to parse template", http.StatusInternalServerError)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Printf("Template execute error: %v\n", err)
		return "", errors.New("Failed to render template", http.StatusInternalServerError)
	}

	return buf.String(), nil
}
