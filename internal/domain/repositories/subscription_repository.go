package repositories

import (
	"context"
	"weather-api/internal/domain/models"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error)
	ExistByLookup(ctx context.Context, lookup *models.SubscriptionLookup) (bool, error)
	Update(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error)
	Delete(ctx context.Context, id uint) error
	FindByToken(ctx context.Context, token string) (*models.Subscription, error)
	FindGroupedSubscriptions(ctx context.Context, frequency *models.Frequency) ([]*models.GroupedSubscription, error)
}
