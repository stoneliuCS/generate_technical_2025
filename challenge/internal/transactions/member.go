package transactions

import (
	"generate_technical_challenge_2025/internal/database/models"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberTransactions interface {
	InsertMember(*models.Member) (*uuid.UUID, error)
}

type MemberTransactionsImpl struct {
	logger *slog.Logger
	db     *gorm.DB
}

func CreateMemberTransactions(logger *slog.Logger, db *gorm.DB) MemberTransactions {
	return &MemberTransactionsImpl{logger: logger, db: db}
}

func (u MemberTransactionsImpl) InsertMember(member *models.Member) (*uuid.UUID, error) {
	res := u.db.Create(&member)
	if res.Error != nil {
		return nil, res.Error
	}
	return &member.ID, nil
}
