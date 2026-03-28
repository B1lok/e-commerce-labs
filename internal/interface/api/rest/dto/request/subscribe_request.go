package request

import "weather-api/internal/application/command"

type SubscribeRequest struct {
	Email     string `form:"email" binding:"required,email"`
	City      string `form:"city" binding:"required"`
	Frequency string `form:"frequency" binding:"required,oneof=hourly daily"`
}

func (req *SubscribeRequest) ToSubscribeCommand() *command.SubscribeCommand {
	return &command.SubscribeCommand{
		Email:     req.Email,
		City:      req.City,
		Frequency: req.Frequency,
	}
}
