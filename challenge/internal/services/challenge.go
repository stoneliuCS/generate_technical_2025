package services

import (
	"generate_technical_challenge_2025/internal/transactions"
	"generate_technical_challenge_2025/internal/utils"
	"log/slog"
	"math"
	"net/url"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ChallengeService interface {
	GenerateUniqueAlienChallenge(id uuid.UUID) map[uuid.UUID]InvasionState
	GenerateUniqueFrontendChallenge(id uuid.UUID) []DetailedAlien
	ScoreMemberSubmission(memberID uuid.UUID, submission map[uuid.UUID]UserChallengeSubmission) OracleAnswer
	GenerateUniqueNgrokChallenge(memberID uuid.UUID) NgrokChallenge
	GradeNgrokServer(url url.URL, requests NgrokChallenge) NgrokChallengeScore
}

type UserChallengeSubmission struct {
	Hp         int
	Commands   []string
	AliensLeft int
}

type OracleAnswer struct {
	Score   int
	Message string
	Valid   bool
}

type ChallengeServiceImpl struct {
	logger       *slog.Logger
	transactions transactions.ChallengeTransactions
}

// ScoreMemberSubmission implements ChallengeService.
func (c ChallengeServiceImpl) ScoreMemberSubmission(memberID uuid.UUID, submission map[uuid.UUID]UserChallengeSubmission) OracleAnswer {
	challenges := c.GenerateUniqueAlienChallenge(memberID)
	// First check is that keys match.
	oracleChallengeKeys := mapset.NewSet(lo.Keys(challenges)...)
	memberChallengeKeys := mapset.NewSet(lo.Keys(submission)...)
	if !oracleChallengeKeys.Equal(memberChallengeKeys) {
		return OracleAnswer{Message: "Challenge IDs do not match.", Valid: false}
	}
	aggregatedAnswer := 0
	// If they do match, now run each simulation and check to see if it agrees with the oracles solution.
	for challengeID, submission := range submission {
		// Next check if the that each string in the commands field must match the commands in the spec.
		invalidCommands := lo.Filter(submission.Commands, func(command string, _ int) bool {
			return command != VOLLEY && command != FOCUSED_SHOT && command != FOCUSED_VOLLEY
		})
		if len(invalidCommands) > 0 {
			return OracleAnswer{Message: "Invalid commands detected for this challenge id: " + challengeID.String(), Valid: false}
		}
		state := challenges[challengeID]
		finalUserState := RunCommandsToCompletion(state, submission.Commands)
		// Check to see if the final HP and final remaining aliens match
		if finalUserState.GetAliensLeft() != submission.AliensLeft || finalUserState.GetHpLeft() != submission.Hp || finalUserState.GetNumberOfCommandsUsed() != len(submission.Commands) {
			return OracleAnswer{Message: "Submission HP, aliens, or commands left do not match for this challenge id: " + challengeID.String(), Valid: false}
		}
		// Run oracles algorithm
		finalOracleState := OracleSolution(state)
		// Take the absolute difference between oracle solution user solution
		//
		hpScore := int(math.Abs(float64(finalOracleState.GetHpLeft()) - float64(finalUserState.GetHpLeft())))
		alienScore := int(math.Abs(float64(finalOracleState.GetAliensLeft()) - float64(finalUserState.GetAliensLeft())))
		commandScore := int(math.Abs(float64(finalOracleState.GetNumberOfCommandsUsed()) - float64(finalUserState.GetNumberOfCommandsUsed())))
		aggregatedAnswer += hpScore + alienScore + commandScore
	}
	return OracleAnswer{Message: "Submission successfully recorded.", Score: aggregatedAnswer, Valid: true}
}

const (
	LOWER_HP_BOUND              = 50
	UPPER_HP_BOUND              = 100
	NUM_WAVES                   = 10
	LOWER_DETAILED_ALIEN_AMOUNT = 10
	UPPER_DETAILED_ALIEN_AMOUNT = 100
	NGROK_GET_REQUEST_COUNT      = 4
	NUM_NGROK_ALIENS_LOWER_BOUND = 10
	NUM_NGROK_ALIENS_UPPER_BOUND = NUM_NGROK_ALIENS_LOWER_BOUND + 5
	NGROK_POST_POINTS            = 20
	NGROK_GET_ALL_POINTS         = 15
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

func (c ChallengeServiceImpl) GradeNgrokServer(url url.URL, requests NgrokChallenge) NgrokChallengeScore {
	panic("Not implemented.")
}

func (c ChallengeServiceImpl) GenerateUniqueNgrokChallenge(memberID uuid.UUID) NgrokChallenge {
	// Use challenge ID as seed for deterministic but unique data.
	rng := utils.CreateRNGFromHash(memberID)

	count := utils.GenerateRandomNumWithinRange(rng, NUM_NGROK_ALIENS_LOWER_BOUND, NUM_NGROK_ALIENS_UPPER_BOUND)
	aliens := []DetailedAlien{}
	for alienIdx := range count {
		alien := GenerateDetailedAlien(rng, memberID, alienIdx)
		aliens = append(aliens, alien)
	}

	requests := []NgrokRequest{
		NgrokPostRequest{
			Name:   "POST all alien",
			Points: NGROK_POST_POINTS,
			Path:   "/api/aliens",
			Body:   aliens,
		},

		NgrokGetRequest{
			Name:          "GET all aliens",
			Points:        NGROK_GET_ALL_POINTS,
			Path:          "/api/aliens",
			ExpectedCount: len(aliens),
		},
	}

	requests = append(requests, generateRandomFilterTests(rng, aliens)...)

	return NgrokChallenge{Requests: requests}
}
