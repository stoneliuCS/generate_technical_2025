package services

import (
	"generate_technical_challenge_2025/internal/database/models"
	"generate_technical_challenge_2025/internal/transactions"
	"log/slog"

	"github.com/google/uuid"
)

// Pointers could be nil or have the actual value, always check if the error is nil 
// before dereferencing a pointer otherwise you may get a null pointer dereference.
type MemberService interface {
	// Inserts a user into the database, returning its unique identifier upon success.
	CreateMember(*models.Member) (*uuid.UUID, error)
	// Gets the id for a member from its email and nuid, or an error if not found.
	GetMember(string, string) (*uuid.UUID, error)
	// Checks if the member exists in the database.
	CheckMemberExists(string, string) (*bool, error)
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

func (u *MemberServiceImpl) CheckMemberExists(email string, nuid string) (*bool, error) {
	return u.transactions.MemberExists(email, nuid)
}

// CreateUser implements UserService.
func (u *MemberServiceImpl) CreateMember(member *models.Member) (*uuid.UUID, error) {
	return u.transactions.InsertMember(member)
}

func (u *MemberServiceImpl) GetMember(email string, nuid string) (*uuid.UUID, error) {
	return u.transactions.GetMember(email, nuid)
}
