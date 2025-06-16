package transactions

import (
	"log/slog"

	"gorm.io/gorm"
)

type Transactions struct{}

func CreateTransactions(logger *slog.Logger, db *gorm.DB) *Transactions {
	return &Transactions{}
}
