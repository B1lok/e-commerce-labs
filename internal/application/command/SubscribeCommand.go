package command

import "weather-api/internal/domain/models"

type SubscribeCommand struct {
	Email     string
	City      string
	Frequency string
}

func (c *SubscribeCommand) ToSubscriptionLookup() *models.SubscriptionLookup {
	return &models.SubscriptionLookup{
		Email:     c.Email,
		City:      c.City,
		Frequency: models.Frequency(c.Frequency),
	}
}
