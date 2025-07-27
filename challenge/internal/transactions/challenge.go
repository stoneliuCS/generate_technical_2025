package transactions

import (
	"generate_technical_challenge_2025/internal/database/models"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const BATCH_SIZE = 20

type ChallengeTransactions interface {
	SaveAlienChallengeSolutionsForMember(sols []models.AlienChallengeSolution) error
	CheckIfMemberHasSolution(memberID uuid.UUID, challengeID uuid.UUID) (bool, error)
}

type ChallengeTransactionsImpl struct {
	logger *slog.Logger
	db     *gorm.DB
}

// CheckIfMemberHasSolutions implements ChallengeTransactions.
func (c ChallengeTransactionsImpl) CheckIfMemberHasSolution(memberID uuid.UUID, challengeID uuid.UUID) (bool, error) {
	panic("unimplemented")
}

// SaveAlienChallengeSolutionsForMember implements ChallengeTransactions.
func (c ChallengeTransactionsImpl) SaveAlienChallengeSolutionsForMember(sols []models.AlienChallengeSolution) error {
	res := c.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(sols, BATCH_SIZE)
	if res.Error != nil {
		return res.Error
	} else {
		return nil
	}
}

func CreateChallengeTransactions(logger *slog.Logger, db *gorm.DB) ChallengeTransactions {
	return ChallengeTransactionsImpl{logger: logger, db: db}
}
