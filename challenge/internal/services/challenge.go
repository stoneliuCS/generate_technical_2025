package services

import (
	"generate_technical_challenge_2025/internal/transactions"
	"log/slog"

	"github.com/google/uuid"
)

// Represents the state of the invasion
type InvasionState struct {
	Budget         uint
	WallDurability uint
}

// ALIEN DATA DEFINITIONS
type AlienType int

const (
	Regular AlienType = iota
	Swift
	Boss
)

// Represents an invading alien
type Alien struct {
	Id   uuid.UUID
	Hp   uint
	Atk  uint
	Type AlienType
}

// WEAPON DATA DEFINITIONS
type WeaponType int

const (
	Turret WeaponType = iota
	MachineGun
	RayGun
)

type Weapon struct {
	Atk  uint
	Cost uint
	Type WeaponType
}

type ChallengeService interface {
	GenerateUniqueChallenge(id uuid.UUID)
}

type ChallengeServiceImpl struct {
	logger       *slog.Logger
	transactions transactions.ChallengeTransactions
}

// GenerateUniqueChallenge implements ChallengeService.
func (c ChallengeServiceImpl) GenerateUniqueChallenge(id uuid.UUID) {
	panic("unimplemented")
}

func CreateChallengeService(logger *slog.Logger, transactions transactions.ChallengeTransactions) ChallengeService {
	return ChallengeServiceImpl{
		logger: logger, transactions: transactions,
	}
}
