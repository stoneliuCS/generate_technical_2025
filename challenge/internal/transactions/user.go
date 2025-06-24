package transactions

import (
	"generate_technical_challenge_2025/internal/database/models"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTransactions interface {
	InsertUser(*models.User) (*uuid.UUID, error)
}

type UserTransactionsImpl struct {
	logger *slog.Logger
	db     *gorm.DB
}

func CreateUserTransactions(logger *slog.Logger, db *gorm.DB) UserTransactions {
	return &UserTransactionsImpl{logger: logger, db: db}
}

func (u UserTransactionsImpl) InsertUser(user *models.User) (*uuid.UUID, error) {
	res := u.db.Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user.ID, nil
}
