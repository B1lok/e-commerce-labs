package interfaces

import (
	"context"
	"weather-api/internal/application/command"
)

type SubscriptionService interface {
	Subscribe(ctx context.Context, subscribeCommand *command.SubscribeCommand) error
	Confirm(ctx context.Context, token string) error
	Unsubscribe(ctx context.Context, token string) error
}
