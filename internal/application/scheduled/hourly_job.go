package scheduled

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"weather-api/internal/application/email"
	"weather-api/internal/domain/models"
	"weather-api/internal/domain/repositories"
)

type HourlyWeatherUpdateJob struct {
	WeatherRepository      repositories.WeatherRepository
	SubscriptionRepository repositories.SubscriptionRepository
	sender                 email.Sender
	host                   string
	workerCount            int
}

func NewHourlyWeatherUpdateJob(
	weatherRepo repositories.WeatherRepository,
	subscriptionRepo repositories.SubscriptionRepository,
	sender email.Sender,
	host string,
) *HourlyWeatherUpdateJob {
	return &HourlyWeatherUpdateJob{
		WeatherRepository:      weatherRepo,
		SubscriptionRepository: subscriptionRepo,
		sender:                 sender,
		host:                   host,
		workerCount:            10,
	}
}

type EmailTask struct {
	subscription  *models.Subscription
	weatherHourly *models.WeatherHourly
}

func (h *HourlyWeatherUpdateJob) Name() string {
	return "HourlyWeatherUpdateJob"
}

func (h *HourlyWeatherUpdateJob) Schedule() string {
	return "0 0 * * * *"
}

func (h *HourlyWeatherUpdateJob) Run(ctx context.Context) error {
	log.Println("HourlyWeatherUpdateJob started...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	freq := models.Frequency("hourly")

	groupedSubscriptions, err := h.SubscriptionRepository.FindGroupedSubscriptions(ctx, &freq)
	if err != nil {
		log.Printf("Failed to fetch subscriptions: %v", err)
		return err
	}

	totalSubscriptions := 0
	for _, group := range groupedSubscriptions {
		totalSubscriptions += len(group.Subscriptions)
	}

	taskChan := make(chan EmailTask, min(totalSubscriptions, 1000))
	errChan := make(chan error, h.workerCount)

	var wg sync.WaitGroup

	for i := 0; i < h.workerCount; i++ {
		wg.Add(1)
		go h.emailWorker(ctx, taskChan, errChan, &wg)
	}

	for _, subscriptionGroup := range groupedSubscriptions {

		select {
		case <-ctx.Done():
			close(taskChan)
			return ctx.Err()
		default:
		}

		weatherHourly, err := h.WeatherRepository.GetHourlyForecast(ctx, subscriptionGroup.City)
		if err != nil {
			log.Printf("Failed to fetch weather for city %s: %v", subscriptionGroup.City, err)
			continue
		}

		for _, subscription := range subscriptionGroup.Subscriptions {
			select {
			case <-ctx.Done():
				close(taskChan)
				return ctx.Err()
			case taskChan <- EmailTask{
				subscription:  subscription,
				weatherHourly: weatherHourly,
			}:
			}
		}
	}

	close(taskChan)

	errorCollector := make(chan []error, 1)
	go func() {
		var errors []error
		for {
			select {
			case err, ok := <-errChan:
				if !ok {
					errorCollector <- errors
					return
				}
				if err != nil {
					errors = append(errors, err)
				}
			case <-ctx.Done():
				errorCollector <- errors
				return
			}
		}
	}()

	wg.Wait()
	close(errChan)

	errors := <-errorCollector

	if len(errors) > 0 {
		log.Printf("Job completed with %d errors", len(errors))
		return errors[0]
	}

	log.Println("HourlyWeatherUpdateJob completed successfully")
	return nil
}

func (h *HourlyWeatherUpdateJob) emailWorker(ctx context.Context, taskChan <-chan EmailTask, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskChan {
		select {
		case <-ctx.Done():
			return
		default:
			err := h.sender.WeatherHourlyEmail(h.toWeatherHourlyEmail(task.subscription, task.weatherHourly))
			if err != nil {
				log.Printf("Error sending email to %s: %v", task.subscription.Email, err)
				select {
				case errChan <- err:
				default:
					log.Printf("Error channel full, additional error: %v", err)
				}
			}
		}
	}
}

func (h *HourlyWeatherUpdateJob) toWeatherHourlyEmail(subscription *models.Subscription, weatherDaily *models.WeatherHourly) *email.WeatherHourlyEmail {
	return &email.WeatherHourlyEmail{
		To:             subscription.Email,
		Frequency:      string(subscription.Frequency),
		WeatherHourly:  weatherDaily,
		UnsubscribeUrl: fmt.Sprintf("%sapi/unsubscribe/%s", h.host, subscription.Token),
	}
}
