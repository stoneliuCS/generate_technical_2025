package services

import (
	"generate_technical_challenge_2025/internal/database/models"
	"generate_technical_challenge_2025/internal/transactions"
	"log/slog"

	"github.com/google/uuid"
)

type UserService interface {
	// Inserts a user into the database, returning its unique identifier upon success.
	CreateUser(*models.User) (*uuid.UUID, error)
}

type UserServiceImpl struct {
	logger       *slog.Logger
	transactions transactions.UserTransactions
}

func CreateUserService(logger *slog.Logger, transactions transactions.UserTransactions) UserService {
	return &UserServiceImpl{
		logger:       logger,
		transactions: transactions,
	}
}

// CreateUser implements UserService.
func (u *UserServiceImpl) CreateUser(user *models.User) (*uuid.UUID, error) {
	return u.transactions.InsertUser(user)
}
