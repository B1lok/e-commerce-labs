package postgres

import (
	"time"
	"weather-api/internal/domain/models"
)

type SubscriptionEntity struct {
	ID        uint `gorm:"primaryKey"`
	Email     string
	City      string
	Frequency models.Frequency `gorm:"type:frequency_enum"`
	Token     string
	Confirmed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SubscriptionEntity) TableName() string {
	return "subscriptions"
}
