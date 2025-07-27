package services

import (
	"generate_technical_challenge_2025/internal/database/models"
	"generate_technical_challenge_2025/internal/transactions"
	"generate_technical_challenge_2025/internal/utils"
	"log/slog"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ChallengeService interface {
	GenerateUniqueAlienChallenge(id uuid.UUID) map[uuid.UUID]InvasionState
	SolveAlienChallenge(state InvasionState) InvasionState
	GenerateUniqueFrontendChallenge(id uuid.UUID) []DetailedAlien
	SaveAlienChallengeAnswers(sols []models.AlienChallengeSolution) error
	CheckIfMemberHasAlienChallengeSolved(memberID uuid.UUID, challengeID uuid.UUID) (bool, error)
}

type ChallengeServiceImpl struct {
	logger       *slog.Logger
	transactions transactions.ChallengeTransactions
}

// CheckIfMemberHasAlienChallengeSolved implements ChallengeService.
func (c ChallengeServiceImpl) CheckIfMemberHasAlienChallengeSolved(memberID uuid.UUID, challengeID uuid.UUID) (bool, error) {
	panic("unimplemented")
}

// SaveAlienChallengeAnswers implements ChallengeService.
func (c ChallengeServiceImpl) SaveAlienChallengeAnswers(sols []models.AlienChallengeSolution) error {
	return c.transactions.SaveAlienChallengeSolutionsForMember(sols)
}

const (
	LOWER_HP_BOUND              = 50
	UPPER_HP_BOUND              = 100
	NUM_WAVES                   = 10
	LOWER_DETAILED_ALIEN_AMOUNT = 10
	UPPER_DETAILED_ALIEN_AMOUNT = 100
	NAMESPACE                   = "ALIEN"
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
func (c ChallengeServiceImpl) GenerateUniqueAlienChallenge(id uuid.UUID) map[uuid.UUID]InvasionState {
	rng := utils.CreateRNGFromHash(id)
	maps := map[uuid.UUID]InvasionState{}
	uuid.SetRand(rng)
	for range NUM_WAVES {
		aliens := GenerateAlienInvasion(rng)
		hp := utils.GenerateRandomNumWithinRange(rng, LOWER_HP_BOUND, UPPER_HP_BOUND)
		invasionState := CreateInvasionState(aliens, hp)
		challengeUUID := uuid.New()
		maps[challengeUUID] = invasionState
	}
	uuid.SetRand(nil)
	return maps
}

// SolveChallenge implements ChallengeService.
func (c ChallengeServiceImpl) SolveAlienChallenge(state InvasionState) InvasionState {
	candidates := RunAllPossibleInvasionStatesToCompletionGreedy(state)
	// We have all the candidates, so simply pick the one that has the smallest number of aliens alive and has the most amount of HP left over.
	// Filter by smallest aliens alive.
	smallestAliensLeft := lo.MinBy(candidates, func(s1 InvasionState, s2 InvasionState) bool {
		return s1.GetAliensLeft() < s2.GetAliensLeft()
	}).GetAliensLeft()
	filteredSmallestAliens := lo.Filter(candidates, func(state InvasionState, idx int) bool {
		return state.GetAliensLeft() == smallestAliensLeft
	})
	// Filter by largest HP remaining. (Any will do)
	bestHP := lo.MaxBy(filteredSmallestAliens, func(s1 InvasionState, s2 InvasionState) bool {
		return s1.GetHpLeft() > s2.GetHpLeft()
	}).GetHpLeft()
	greatestHPCandidates := lo.Filter(filteredSmallestAliens, func(state InvasionState, idx int) bool {
		return state.GetHpLeft() == bestHP
	})
	// Finally filter by the least number of commands used.
	idealCandidate := lo.MinBy(greatestHPCandidates, func(s1 InvasionState, s2 InvasionState) bool {
		return s1.GetNumberOfCommandsUsed() < s2.GetNumberOfCommandsUsed()
	})
	return idealCandidate
}

func CreateChallengeService(logger *slog.Logger, transactions transactions.ChallengeTransactions) ChallengeService {
	return ChallengeServiceImpl{
		logger: logger, transactions: transactions,
	}
}
