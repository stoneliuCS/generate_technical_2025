package services

import (
	"generate_technical_challenge_2025/internal/database/models"
	"generate_technical_challenge_2025/internal/transactions"
	"generate_technical_challenge_2025/internal/utils"
	"log/slog"

	"github.com/google/uuid"
)

// Pointers could be nil or have the actual value, always check if the error is nil
// before dereferencing a pointer otherwise you may get a null pointer dereference.
type MemberService interface {
	CreateMember(*models.Member) (*uuid.UUID, error)
	CreateScore(*models.Score) (int, error)
	LogFrontendUsageAsync(uuid.UUID)
	GetMember(string, string) (*uuid.UUID, error)
	CheckMemberExistsByEmailAndNuid(string, string) (bool, error)
	CheckMemberExistsById(uuid.UUID) (bool, error)
}

type MemberServiceImpl struct {
	logger       *slog.Logger
	transactions transactions.MemberTransactions
	usageLogger  *utils.UsageLogger
}

// LogFrontendUsageAsync implements MemberService.
// No need to call this within a goroutine, it already handles asynchronicity.
func (u *MemberServiceImpl) LogFrontendUsageAsync(userID uuid.UUID) {
	u.usageLogger.LogUsage(userID)
}

// CreateScore implements MemberService.
func (u *MemberServiceImpl) CreateScore(score *models.Score) (int, error) {
	return u.transactions.InsertScore(score)
}

// CheckMemberExistsById implements MemberService.
func (u *MemberServiceImpl) CheckMemberExistsById(id uuid.UUID) (bool, error) {
	return u.transactions.MemberExistsById(id)
}

func CreateMemberService(logger *slog.Logger, transactions transactions.MemberTransactions) MemberService {
	usageLogger := utils.NewUsageLogger(transactions)
	return &MemberServiceImpl{
		logger:       logger,
		transactions: transactions,
		usageLogger:  usageLogger,
	}
}

func (u *MemberServiceImpl) CheckMemberExistsByEmailAndNuid(email string, nuid string) (bool, error) {
	return u.transactions.MemberExistsByEmailAndNuid(email, nuid)
}

// CreateUser implements UserService.
func (u *MemberServiceImpl) CreateMember(member *models.Member) (*uuid.UUID, error) {
	return u.transactions.InsertMember(member)
}

func (u *MemberServiceImpl) GetMember(email string, nuid string) (*uuid.UUID, error) {
	return u.transactions.GetMember(email, nuid)
}
