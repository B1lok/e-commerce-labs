package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"weather-api/internal/application/command"
	"weather-api/internal/domain/models"
	"weather-api/internal/test/mocks"
	"weather-api/pkg/errors"
)

func TestSubscriptionService_Subscribe_Success(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	cmd := &command.SubscribeCommand{
		Email:     "test@example.com",
		City:      "berlin",
		Frequency: "daily",
	}

	validatedCity := "Berlin"
	mockValidator.On("Validate", "berlin").Return(validatedCity, nil)

	lookup := &models.SubscriptionLookup{
		Email:     "test@example.com",
		City:      validatedCity,
		Frequency: models.Frequency("daily"),
	}

	subscription := &models.Subscription{
		Email:     "test@example.com",
		City:      validatedCity,
		Frequency: models.Frequency("daily"),
	}
	mockRepo.On("ExistByLookup", ctx, lookup).Return(false, nil)

	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Subscription")).Return(subscription, nil)

	mockSender.On("ConfirmationEmail", mock.AnythingOfType("*email.ConfirmationEmail")).Return(nil)

	err := service.Subscribe(ctx, cmd)

	assert.NoError(t, err)
	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

func TestSubscriptionService_Subscribe_InvalidCity(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	cmd := &command.SubscribeCommand{
		Email:     "test@example.com",
		City:      "invalidcity",
		Frequency: "daily",
	}

	validationErr := errors.New("City not found", http.StatusNotFound)
	mockValidator.On("Validate", "invalidcity").Return(nil, validationErr)

	err := service.Subscribe(ctx, cmd)

	assert.Error(t, err)
	assert.Equal(t, validationErr, err)
	mockValidator.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "ExistByLookup")
	mockRepo.AssertNotCalled(t, "Create")
	mockSender.AssertNotCalled(t, "ConfirmationEmail")
}

func TestSubscriptionService_Subscribe_AlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	cmd := &command.SubscribeCommand{
		Email:     "test@example.com",
		City:      "berlin",
		Frequency: "daily",
	}

	validatedCity := "Berlin"
	mockValidator.On("Validate", "berlin").Return(validatedCity, nil)

	lookup := &models.SubscriptionLookup{
		Email:     "test@example.com",
		City:      validatedCity,
		Frequency: "daily",
	}
	mockRepo.On("ExistByLookup", ctx, lookup).Return(true, nil)

	err := service.Subscribe(ctx, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Email already subscribed")
	apiErr, ok := errors.IsAPIError(err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusConflict, apiErr.Code)

	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
	mockSender.AssertNotCalled(t, "ConfirmationEmail")
}

func TestSubscriptionService_Subscribe_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	cmd := &command.SubscribeCommand{
		Email:     "test@example.com",
		City:      "berlin",
		Frequency: "daily",
	}

	validatedCity := "Berlin"
	mockValidator.On("Validate", "berlin").Return(validatedCity, nil)

	lookup := &models.SubscriptionLookup{
		Email:     "test@example.com",
		City:      validatedCity,
		Frequency: "daily",
	}
	repoErr := errors.New("Database connection failed", http.StatusInternalServerError)
	mockRepo.On("ExistByLookup", ctx, lookup).Return(false, repoErr)

	err := service.Subscribe(ctx, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check if email exists")
	apiErr, ok := errors.IsAPIError(err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Code)

	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
	mockSender.AssertNotCalled(t, "ConfirmationEmail")
}

func TestSubscriptionService_Subscribe_InvalidFrequency(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	cmd := &command.SubscribeCommand{
		Email:     "test@example.com",
		City:      "berlin",
		Frequency: "invalid",
	}

	validatedCity := "Berlin"
	mockValidator.On("Validate", "berlin").Return(validatedCity, nil)

	validatedCmd := &command.SubscribeCommand{
		Email:     "test@example.com",
		City:      "Berlin",
		Frequency: "invalid",
	}

	mockRepo.On("ExistByLookup", ctx, validatedCmd.ToSubscriptionLookup()).Return(false, nil)

	err := service.Subscribe(ctx, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid input")
	apiErr, ok := errors.IsAPIError(err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)

	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
	mockSender.AssertNotCalled(t, "ConfirmationEmail")
}

func TestSubscriptionService_Subscribe_CreateError(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	cmd := &command.SubscribeCommand{
		Email:     "test@example.com",
		City:      "berlin",
		Frequency: "daily",
	}

	validatedCity := "Berlin"
	mockValidator.On("Validate", "berlin").Return(validatedCity, nil)

	lookup := &models.SubscriptionLookup{
		Email:     "test@example.com",
		City:      validatedCity,
		Frequency: "daily",
	}
	mockRepo.On("ExistByLookup", ctx, lookup).Return(false, nil)

	createErr := errors.New("Database error", http.StatusInternalServerError)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Subscription")).Return(nil, createErr)

	err := service.Subscribe(ctx, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create subscription")
	apiErr, ok := errors.IsAPIError(err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Code)

	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockSender.AssertNotCalled(t, "ConfirmationEmail")
}

func TestSubscriptionService_Subscribe_EmailError(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	cmd := &command.SubscribeCommand{
		Email:     "test@example.com",
		City:      "berlin",
		Frequency: "daily",
	}

	validatedCity := "Berlin"
	mockValidator.On("Validate", "berlin").Return(validatedCity, nil)

	lookup := &models.SubscriptionLookup{
		Email:     "test@example.com",
		City:      validatedCity,
		Frequency: "daily",
	}

	mockRepo.On("ExistByLookup", ctx, lookup).Return(false, nil)
	subscription := &models.Subscription{
		Email:     "test@example.com",
		City:      validatedCity,
		Frequency: models.Frequency("daily"),
	}
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Subscription")).Return(subscription, nil)

	emailErr := errors.New("Failed to send email", http.StatusInternalServerError)
	mockSender.On("ConfirmationEmail", mock.AnythingOfType("*email.ConfirmationEmail")).Return(emailErr)

	err := service.Subscribe(ctx, cmd)

	assert.Error(t, err)
	assert.Equal(t, emailErr, err)

	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

func TestSubscriptionService_Confirm_Success(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	token := "valid-token-123"

	subscription := &models.Subscription{
		ID:        1,
		Email:     "test@example.com",
		City:      "Berlin",
		Frequency: models.Frequency("daily"),
		Token:     token,
		Confirmed: false,
	}

	mockRepo.On("FindByToken", ctx, token).Return(subscription, nil)

	updatedSubscription := &models.Subscription{
		ID:        1,
		Email:     "test@example.com",
		City:      "Berlin",
		Frequency: models.Frequency("daily"),
		Token:     token,
		Confirmed: true,
	}
	mockRepo.On("Update", ctx, mock.AnythingOfType("*models.Subscription")).Return(updatedSubscription, nil)

	err := service.Confirm(ctx, token)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_Confirm_TokenNotFound(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	token := "invalid-token"

	mockRepo.On("FindByToken", ctx, token).Return(nil, nil)

	err := service.Confirm(ctx, token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Token not found")
	apiErr, ok := errors.IsAPIError(err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestSubscriptionService_Confirm_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	token := "valid-token-123"

	repoErr := errors.New("Database error", http.StatusInternalServerError)
	mockRepo.On("FindByToken", ctx, token).Return(nil, repoErr)

	err := service.Confirm(ctx, token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find subscription")
	apiErr, ok := errors.IsAPIError(err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Code)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestSubscriptionService_Confirm_UpdateError(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	token := "valid-token-123"

	subscription := &models.Subscription{
		ID:        1,
		Email:     "test@example.com",
		City:      "Berlin",
		Frequency: models.Frequency("daily"),
		Token:     token,
		Confirmed: false,
	}

	mockRepo.On("FindByToken", ctx, token).Return(subscription, nil)

	updateErr := errors.New("Update failed", http.StatusInternalServerError)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*models.Subscription")).Return(nil, updateErr)

	err := service.Confirm(ctx, token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update subscription")
	apiErr, ok := errors.IsAPIError(err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_Unsubscribe_Success(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	token := "valid-token-123"

	subscription := &models.Subscription{
		ID:        1,
		Email:     "test@example.com",
		City:      "Berlin",
		Frequency: models.Frequency("daily"),
		Token:     token,
		Confirmed: true,
	}

	mockRepo.On("FindByToken", ctx, token).Return(subscription, nil)
	mockRepo.On("Delete", ctx, uint(1)).Return(nil)

	err := service.Unsubscribe(ctx, token)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_Unsubscribe_TokenNotFound(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	token := "invalid-token"

	mockRepo.On("FindByToken", ctx, token).Return(nil, nil)

	err := service.Unsubscribe(ctx, token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Token not found")
	apiErr, ok := errors.IsAPIError(err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestSubscriptionService_Unsubscribe_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	token := "valid-token-123"

	repoErr := errors.New("Database error", http.StatusInternalServerError)
	mockRepo.On("FindByToken", ctx, token).Return(nil, repoErr)

	err := service.Unsubscribe(ctx, token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find subscription")
	apiErr, ok := errors.IsAPIError(err)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Code)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestSubscriptionService_Unsubscribe_DeleteError(t *testing.T) {
	mockRepo := new(mocks.MockSubscriptionRepository)
	mockValidator := new(mocks.MockCityValidator)
	mockSender := new(mocks.MockEmailSender)
	host := "http://example.com/"

	service := NewSubscriptionService(mockRepo, mockValidator, mockSender, host)

	ctx := context.Background()
	token := "valid-token-123"

	subscription := &models.Subscription{
		ID:        1,
		Email:     "test@example.com",
		City:      "Berlin",
		Frequency: models.Frequency("daily"),
		Token:     token,
		Confirmed: true,
	}

	mockRepo.On("FindByToken", ctx, token).Return(subscription, nil)

	deleteErr := errors.New("Delete failed", http.StatusInternalServerError)
	mockRepo.On("Delete", ctx, uint(1)).Return(deleteErr)

	err := service.Unsubscribe(ctx, token)

	assert.Error(t, err)
	assert.Equal(t, deleteErr, err)

	mockRepo.AssertExpectations(t)
}
