package postgres

import (
	"weather-api/internal/domain/models"
)

func ToEntity(subscription *models.Subscription) *SubscriptionEntity {
	if subscription == nil {
		return nil
	}

	return &SubscriptionEntity{
		ID:        subscription.ID,
		Email:     subscription.Email,
		City:      subscription.City,
		Frequency: subscription.Frequency,
		Token:     subscription.Token,
		Confirmed: subscription.Confirmed,
		CreatedAt: subscription.CreatedAt,
		UpdatedAt: subscription.UpdatedAt,
	}
}

func ToDomain(entity *SubscriptionEntity) (*models.Subscription, error) {
	if entity == nil {
		return nil, nil
	}

	subscription := &models.Subscription{
		ID:        entity.ID,
		Email:     entity.Email,
		City:      entity.City,
		Frequency: entity.Frequency,
		Token:     entity.Token,
		Confirmed: entity.Confirmed,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	return subscription, nil
}

func ToDomainList(entities []SubscriptionEntity) ([]*models.Subscription, error) {
	if entities == nil {
		return nil, nil
	}

	result := make([]*models.Subscription, 0, len(entities))
	for _, entity := range entities {
		sub, err := ToDomain(&entity)
		if err != nil {
			return nil, err
		}
		result = append(result, sub)
	}

	return result, nil
}
