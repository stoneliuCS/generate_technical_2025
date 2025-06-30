package services

import (
	"generate_technical_challenge_2025/internal/transactions"
	"log/slog"

	"github.com/google/uuid"
)

type ChallengeService interface {
	GenerateUniqueAlienChallenge(id uuid.UUID) InvasionState
	SolveAlienChallenge(state InvasionState) InvasionState
}

type ChallengeServiceImpl struct {
	logger       *slog.Logger
	transactions transactions.ChallengeTransactions
}

// GenerateUniqueAlienChallenge implements ChallengeService.
func (c ChallengeServiceImpl) GenerateUniqueAlienChallenge(id uuid.UUID) InvasionState {
	return GenerateInvasionState(id)
}

// SolveChallenge implements ChallengeService.
func (c ChallengeServiceImpl) SolveAlienChallenge(state InvasionState) InvasionState {
	// Define all the states we want to run against
	allInvasionStates := []InvasionState{}
	// First generate all possible combinations of weapon purchases given the current budget.
	allPossibleGunPurchases := GenerateAllPossibleWeaponPurchasesFromBudget(state.Budget)
	for _, weaponPurchases := range allPossibleGunPurchases {
		// Begin with the same state across all weapon purchases
		currentState := state
		for _, weapon := range weaponPurchases {
			currentState = currentState.PurchaseWeapon(weapon)
		}
		allInvasionStates = append(allInvasionStates, currentState)
	}
	panic("Not implemented.")
}

func CreateChallengeService(logger *slog.Logger, transactions transactions.ChallengeTransactions) ChallengeService {
	return ChallengeServiceImpl{
		logger: logger, transactions: transactions,
	}
}
