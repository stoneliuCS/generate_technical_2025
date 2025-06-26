package services

import (
	"generate_technical_challenge_2025/internal/database/models"
	"generate_technical_challenge_2025/internal/transactions"
	"log/slog"

	"github.com/google/uuid"
)

type MemberService interface {
	// Inserts a user into the database, returning its unique identifier upon success.
	CreateMember(*models.Member) (*uuid.UUID, error)
	GetMember(string, string) (*uuid.UUID, error)
}

type MemberServiceImpl struct {
	logger       *slog.Logger
	transactions transactions.MemberTransactions
}

func CreateMemberService(logger *slog.Logger, transactions transactions.MemberTransactions) MemberService {
	return &MemberServiceImpl{
		logger:       logger,
		transactions: transactions,
	}
}

// CreateUser implements UserService.
func (u *MemberServiceImpl) CreateMember(member *models.Member) (*uuid.UUID, error) {
	return u.transactions.InsertMember(member)
}

func (u *MemberServiceImpl) GetMember(email string, nuid string) (*uuid.UUID, error) {
	return u.transactions.GetMember(email, nuid)
}
