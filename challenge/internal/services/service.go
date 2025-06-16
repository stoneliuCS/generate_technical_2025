package services

import (
	"generate_technical_challenge_2025/internal/transactions"
	"log/slog"
)

// Services available in Hangouts
type Services struct{}

// Creates all the services available.
func CreateServices(logger *slog.Logger, transactions *transactions.Transactions) *Services {
	return &Services{}
}
