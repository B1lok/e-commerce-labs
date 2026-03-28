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

type DailyWeatherUpdateJob struct {
	WeatherRepository      repositories.WeatherRepository
	SubscriptionRepository repositories.SubscriptionRepository
	sender                 email.Sender
	host                   string
	workerCount            int
}

func NewDailyWeatherUpdateJob(
	weatherRepo repositories.WeatherRepository,
	subscriptionRepo repositories.SubscriptionRepository,
	sender email.Sender,
	host string,
) *DailyWeatherUpdateJob {
	return &DailyWeatherUpdateJob{
		WeatherRepository:      weatherRepo,
		SubscriptionRepository: subscriptionRepo,
		sender:                 sender,
		host:                   host,
		workerCount:            10,
	}
}

type DailyEmailTask struct {
	subscription *models.Subscription
	weatherDaily *models.WeatherDaily
}

func (d *DailyWeatherUpdateJob) Name() string {
	return "DailyWeatherUpdateJob"
}

func (d *DailyWeatherUpdateJob) Schedule() string {
	return "0 0 8 * * *"
}

func (d *DailyWeatherUpdateJob) Run(ctx context.Context) error {
	log.Println("DailyWeatherUpdateJob started...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	freq := models.Frequency("daily")

	groupedSubscriptions, err := d.SubscriptionRepository.FindGroupedSubscriptions(ctx, &freq)
	if err != nil {
		log.Printf("Failed to fetch subscriptions: %v", err)
		return err
	}

	totalSubscriptions := 0
	for _, group := range groupedSubscriptions {
		totalSubscriptions += len(group.Subscriptions)
	}

	taskChan := make(chan DailyEmailTask, min(totalSubscriptions, 1000))
	errChan := make(chan error, d.workerCount)

	var wg sync.WaitGroup

	for i := 0; i < d.workerCount; i++ {
		wg.Add(1)
		go d.emailWorker(ctx, taskChan, errChan, &wg)
	}

	for _, subscriptionGroup := range groupedSubscriptions {

		select {
		case <-ctx.Done():
			close(taskChan)
			return ctx.Err()
		default:
		}

		weatherDaily, err := d.WeatherRepository.GetDailyForecast(ctx, subscriptionGroup.City)
		if err != nil {
			log.Printf("Failed to fetch weather for city %s: %v", subscriptionGroup.City, err)
			continue
		}

		for _, subscription := range subscriptionGroup.Subscriptions {
			select {
			case <-ctx.Done():
				close(taskChan)
				return ctx.Err()
			case taskChan <- DailyEmailTask{
				subscription: subscription,
				weatherDaily: weatherDaily,
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

	log.Println("DailyWeatherUpdateJob completed successfully")
	return nil
}

func (d *DailyWeatherUpdateJob) emailWorker(ctx context.Context, taskChan <-chan DailyEmailTask, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskChan {
		select {
		case <-ctx.Done():
			return
		default:
			err := d.sender.WeatherDailyEmail(d.toWeatherDailyEmail(task.subscription, task.weatherDaily))
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

func (d *DailyWeatherUpdateJob) toWeatherDailyEmail(subscription *models.Subscription, weatherDaily *models.WeatherDaily) *email.WeatherDailyEmail {
	return &email.WeatherDailyEmail{
		To:             subscription.Email,
		Frequency:      string(subscription.Frequency),
		WeatherDaily:   weatherDaily,
		UnsubscribeUrl: fmt.Sprintf("%sapi/unsubscribe/%s", d.host, subscription.Token),
	}
}
