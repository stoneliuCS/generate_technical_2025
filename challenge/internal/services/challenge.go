package services

import (
	"generate_technical_challenge_2025/internal/transactions"
	"generate_technical_challenge_2025/internal/utils"
	"log/slog"
	"net/url"

	"github.com/google/uuid"
)

type ChallengeService interface {
	GenerateUniqueAlienChallenge(id uuid.UUID) []InvasionState
	SolveAlienChallenge(state InvasionState) InvasionState
	GenerateUniqueNgrokChallenge(id uuid.UUID) NgrokChallenge
	GradeNgrokServer(url url.URL, requests NgrokChallenge) NgrokChallengeScore
}

type ChallengeServiceImpl struct {
	logger       *slog.Logger
	transactions transactions.ChallengeTransactions
}

const (
	LOWER_HP_BOUND               = 50
	UPPER_HP_BOUND               = 100
	NUM_WAVES_LOWER_BOUND        = 5
	NUM_WAVES_UPPER_BOUND        = 10
	NGROK_GET_REQUEST_COUNT      = 4
	NUM_NGROK_ALIENS_LOWER_BOUND = 10
	NUM_NGROK_ALIENS_UPPER_BOUND = NUM_NGROK_ALIENS_LOWER_BOUND + 5
	NGROK_POST_POINTS            = 20
	NGROK_GET_ALL_POINTS         = 15
)

func (c ChallengeServiceImpl) GradeNgrokServer(url url.URL, requests NgrokChallenge) NgrokChallengeScore {
	panic("Not implemented.")
}

func (c ChallengeServiceImpl) GenerateUniqueNgrokChallenge(id uuid.UUID) NgrokChallenge {
	// Use challenge ID as seed for deterministic but unique data
	rng := utils.CreateRNGFromHash(id)

	count := utils.GenerateRandomNumWithinRange(rng, NUM_NGROK_ALIENS_LOWER_BOUND, NUM_NGROK_ALIENS_UPPER_BOUND)
	aliens := generateRandomAliens(rng, count)

	requests := []NgrokRequest{
		NgrokPostRequest{
			Name:   "POST all alien",
			Points: NGROK_POST_POINTS,
			Path:   "/todo",
			Body:   aliens,
		},

		NgrokGetRequest{
			Name:          "GET all aliens",
			Points:        NGROK_GET_ALL_POINTS,
			Path:          "/todo",
			ExpectedCount: len(aliens),
		},
	}

	requests = append(requests, generateRandomFilterTests(rng, aliens)...)

	return NgrokChallenge{Requests: requests}
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
