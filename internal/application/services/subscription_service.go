package services

import (
	"context"
	"net/http"
	"weather-api/internal/application/command"
	"weather-api/internal/application/email"
	"weather-api/internal/domain/models"
	"weather-api/internal/domain/repositories"
	"weather-api/internal/domain/validator"
	"weather-api/pkg/errors"
)

type SubscriptionService struct {
	repository repositories.SubscriptionRepository
	validator  validator.CityValidator
	sender     email.Sender
	host       string
}

func NewSubscriptionService(repository repositories.SubscriptionRepository, validator validator.CityValidator, sender email.Sender, host string) *SubscriptionService {
	return &SubscriptionService{repository: repository, validator: validator, sender: sender, host: host}
}

func (s *SubscriptionService) Subscribe(ctx context.Context, subscribeCommand *command.SubscribeCommand) error {
	validatedCity, err := s.validator.Validate(subscribeCommand.City)
	if err != nil {
		return err
	}
	subscribeCommand.City = *validatedCity
	exists, err := s.repository.ExistByLookup(ctx, subscribeCommand.ToSubscriptionLookup())
	if err != nil {
		return errors.Wrap(err, "failed to check if email exists", http.StatusInternalServerError)
	}
	if exists {
		return errors.New("Email already subscribed", http.StatusConflict)
	}

	newSubscription, err := models.NewSubscription(
		subscribeCommand.Email,
		subscribeCommand.City,
		models.Frequency(subscribeCommand.Frequency),
	)
	if err != nil {
		return errors.Wrap(err, "Invalid input", http.StatusBadRequest)
	}

	savedSubscription, err := s.repository.Create(ctx, newSubscription)
	if err != nil {
		return errors.Wrap(err, "failed to create subscription", http.StatusInternalServerError)
	}
	confirmationEmail := &email.ConfirmationEmail{
		To:        savedSubscription.Email,
		City:      savedSubscription.City,
		Frequency: string(savedSubscription.Frequency),
		Url:       s.host + "api/confirm/" + newSubscription.Token,
	}
	if err := s.sender.ConfirmationEmail(confirmationEmail); err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) Confirm(ctx context.Context, token string) error {
	subscription, err := s.repository.FindByToken(ctx, token)
	if err != nil {
		return errors.Wrap(err, "failed to find subscription", http.StatusInternalServerError)
	}

	if subscription == nil {
		return errors.New("Token not found", http.StatusNotFound)
	}

	subscription.SetConfirmed(true)

	_, err = s.repository.Update(ctx, subscription)
	if err != nil {
		return errors.Wrap(err, "failed to update subscription", http.StatusInternalServerError)
	}
	return err
}

func (s *SubscriptionService) Unsubscribe(ctx context.Context, token string) error {
	subscription, err := s.repository.FindByToken(ctx, token)
	if err != nil {
		return errors.Wrap(err, "failed to find subscription", http.StatusInternalServerError)
	}

	if subscription == nil {
		return errors.New("Token not found", http.StatusNotFound)
	}

	return s.repository.Delete(ctx, subscription.ID)
}
