package services

import (
	"generate_technical_challenge_2025/internal/transactions"
	"generate_technical_challenge_2025/internal/utils"
	"log/slog"

	"github.com/google/uuid"
)

type ChallengeService interface {
	GenerateUniqueAlienChallenge(id uuid.UUID) []InvasionState
	SolveAlienChallenge(state InvasionState) InvasionState
	GenerateUniqueFrontendChallenge(id uuid.UUID) []DetailedAlien
}

type ChallengeServiceImpl struct {
	logger       *slog.Logger
	transactions transactions.ChallengeTransactions
}

const (
	LOWER_HP_BOUND              = 50
	UPPER_HP_BOUND              = 100
	NUM_WAVES_LOWER_BOUND       = 5
	NUM_WAVES_UPPER_BOUND       = 10
	LOWER_DETAILED_ALIEN_AMOUNT = 10
	UPPER_DETAILED_ALIEN_AMOUNT = 100
)

var alienTypes = []AlienType{
	AlienTypeRegular,
	AlienTypeElite,
	AlienTypeBoss,
}

// GenerateUniqueFrontendChallenge implements ChallengeService.
func (c ChallengeServiceImpl) GenerateUniqueFrontendChallenge(id uuid.UUID) []DetailedAlien {
	rng := utils.CreateRNGFromHash(id)
	numAliens := utils.GenerateRandomNumWithinRange(rng, LOWER_DETAILED_ALIEN_AMOUNT, UPPER_DETAILED_ALIEN_AMOUNT)

	aliens := []DetailedAlien{}
	for idx := range numAliens {
		alien := GenerateDetailedAlien(rng, id, idx)
		aliens = append(aliens, alien)
	}

	return aliens
}

// GenerateUniqueAlienChallenge implements ChallengeService.
func (c ChallengeServiceImpl) GenerateUniqueAlienChallenge(id uuid.UUID) []InvasionState {
	rng := utils.CreateRNGFromHash(id)
	numWaves := utils.GenerateRandomNumWithinRange(rng, NUM_WAVES_LOWER_BOUND, NUM_WAVES_UPPER_BOUND)
	waves := []InvasionState{}
	for range numWaves {
		aliens := GenerateAlienInvasion(rng)
		hp := utils.GenerateRandomNumWithinRange(rng, LOWER_HP_BOUND, UPPER_HP_BOUND)
		invasionState := CreateInvasionState(aliens, hp)
		waves = append(waves, invasionState)
	}
	return waves
}

// SolveChallenge implements ChallengeService.
func (c ChallengeServiceImpl) SolveAlienChallenge(state InvasionState) InvasionState {
	panic("Not implemented.")
}

func CreateChallengeService(logger *slog.Logger, transactions transactions.ChallengeTransactions) ChallengeService {
	return ChallengeServiceImpl{
		logger: logger, transactions: transactions,
	}
}
