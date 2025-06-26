package transactions

import (
	"generate_technical_challenge_2025/internal/database/models"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberTransactions interface {
	InsertMember(*models.Member) (*uuid.UUID, error)
	GetMember(string, string) (*uuid.UUID, error)
	MemberExistsByEmailAndNuid(string, string) (bool, error)
	MemberExistsById(uuid.UUID) (bool, error)
}

type MemberTransactionsImpl struct {
	logger *slog.Logger
	db     *gorm.DB
}

// MemberExistsById implements MemberTransactions.
func (u *MemberTransactionsImpl) MemberExistsById(id uuid.UUID) (bool, error) {
	var member models.Member
	res := u.db.Where("id = ?", id).Limit(1).Find(&member)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func CreateMemberTransactions(logger *slog.Logger, db *gorm.DB) MemberTransactions {
	return &MemberTransactionsImpl{logger: logger, db: db}
}

// MemberExists implements MemberTransactions.
func (u *MemberTransactionsImpl) MemberExistsByEmailAndNuid(email string, nuid string) (bool, error) {
	var member models.Member
	res := u.db.Where("email = ?", email).Where("nuid = ?", nuid).Limit(1).Find(&member)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (u MemberTransactionsImpl) InsertMember(member *models.Member) (*uuid.UUID, error) {
	res := u.db.Create(&member)
	if res.Error != nil {
		return nil, res.Error
	}
	return &member.ID, nil
}

// GetMember implements MemberTransactions.
func (u *MemberTransactionsImpl) GetMember(email string, nuid string) (*uuid.UUID, error) {
	var member models.Member
	res := u.db.Where("email = ?", email).Where("nuid = ?", nuid).First(&member)
	if res.Error != nil {
		return nil, res.Error
	}
	return &member.ID, nil
}
