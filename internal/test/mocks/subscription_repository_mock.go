package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"weather-api/internal/domain/models"
)

type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) ExistByLookup(ctx context.Context, lookup *models.SubscriptionLookup) (bool, error) {
	args := m.Called(ctx, lookup)
	return args.Bool(0), args.Error(1)
}

func (m *MockSubscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
	args := m.Called(ctx, subscription)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) FindByToken(ctx context.Context, token string) (*models.Subscription, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Update(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
	args := m.Called(ctx, subscription)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) FindGroupedSubscriptions(ctx context.Context, frequency *models.Frequency) ([]*models.GroupedSubscription, error) {
	args := m.Called(ctx, frequency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.GroupedSubscription), args.Error(1)
}
