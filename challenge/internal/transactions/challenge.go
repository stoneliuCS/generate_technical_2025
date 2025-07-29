package transactions

import (
	"log/slog"

	"gorm.io/gorm"
)

type ChallengeTransactions interface{}

type ChallengeTransactionsImpl struct {
	logger *slog.Logger
	db     *gorm.DB
}

// SaveAlienChallengeSolutionsForMember implements ChallengeTransactions.

func CreateChallengeTransactions(logger *slog.Logger, db *gorm.DB) ChallengeTransactions {
	return ChallengeTransactionsImpl{logger: logger, db: db}
}
