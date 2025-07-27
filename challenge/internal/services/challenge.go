package services

import (
	"context"
	"fmt"
	"generate_technical_challenge_2025/internal/transactions"
	"generate_technical_challenge_2025/internal/utils"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ChallengeService interface {
	GenerateUniqueAlienChallenge(memberID uuid.UUID) []InvasionState
	SolveAlienChallenge(state InvasionState) InvasionState
	GenerateUniqueFrontendChallenge(memberID uuid.UUID) []DetailedAlien
	GenerateUniqueNgrokChallenge(memberID uuid.UUID) NgrokChallenge
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
	LOWER_DETAILED_ALIEN_AMOUNT  = 10
	UPPER_DETAILED_ALIEN_AMOUNT  = 100
	NGROK_GET_REQUEST_COUNT      = 4
	NUM_NGROK_ALIENS_LOWER_BOUND = 10
	NUM_NGROK_ALIENS_UPPER_BOUND = NUM_NGROK_ALIENS_LOWER_BOUND + 5
	NGROK_POST_POINTS            = 20
	NGROK_GET_ALL_POINTS         = 15
	NGROK_FILTER_TYPE_POINTS     = 15
	NGROK_FILTER_SPD_POINTS      = 15
	VERBOSE                      = false
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

func (c ChallengeServiceImpl) GradeNgrokServer(url url.URL, requests NgrokChallenge) NgrokChallengeScore {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	baseURL := url.String()

	ok, err := health(ctx, baseURL)
	if err != nil || !ok {
		return NgrokChallengeScore{
			Valid:  false,
			Reason: "Health check failed - server unreachable",
		}
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	totalPossiblePoints := lo.Reduce(requests.Requests, func(acc int, req NgrokRequest, index int) int {
		return acc + req.GetTotalPossiblePoints()
	}, 0)

	totalScore := 0

	var deleteRequest NgrokRequest
	var postRequest NgrokRequest
	var getRequests []NgrokRequest

	for _, request := range requests.Requests {
		switch req := request.(type) {
		case NgrokDeleteRequest:
			deleteRequest = req
		case NgrokPostRequest:
			postRequest = req
		case NgrokGetRequest:
			getRequests = append(getRequests, req)
		}
	}

	// 1. Clean up any old data.
	if deleteRequest != nil {
		_, err := deleteRequest.Execute(client, baseURL)
		if err != nil {
			if VERBOSE {
				fmt.Printf("DELETE request failed: %s\n", err.Error())
			}
			// Don't fail grading if DELETE fails, because it might be the first run.
		} else {
			if VERBOSE {
				fmt.Println("DELETE request succeeded")
			}
		}
		time.Sleep(200 * time.Millisecond)
	}

	// 2) Populate data.
	if postRequest != nil {
		points, err := postRequest.Execute(client, baseURL)
		if err != nil {
			if VERBOSE {
				fmt.Printf("POST request failed: %s\n", err.Error())
			}
			return NgrokChallengeScore{
				Valid:  false,
				Reason: fmt.Sprintf("POST request failed - %s", err.Error()),
			}
		} else {
			totalScore += points
			if VERBOSE {
				fmt.Printf("POST request succeeded (+%d points out of %d total points)\n", points, postRequest.GetTotalPossiblePoints())
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	// 3) Make GET requests.
	for _, getRequest := range getRequests {
		if VERBOSE {
			fmt.Printf("Request: %s\n", getRequest.GetName())
		}
		points, err := getRequest.Execute(client, baseURL)
		if err != nil {
			if VERBOSE {
				fmt.Printf("GET request failed: %s\n", err.Error())
			}
		} else {
			totalScore += points
			if VERBOSE {
				fmt.Printf("GET request succeeded (+%d points out of %d total points)\n", points, getRequest.GetTotalPossiblePoints())
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	return NgrokChallengeScore{
		Valid: true,
		Score: totalPossiblePoints - totalScore,
	}
}

// An NgrokChallenge consists of:
//  1. a DELETE request that deletes all aliens in the candidate's DB.
//  2. a POST request that creates a deterministic set of aliens in the candidate's DB.
//  3. A GET request with no filters that retrieves all aliens from the candidate's DB.
//  4. A GET request with randomized filters:
//     A filter is one of:
//     - type=
//     - atk_lte=
//     - atk_gte=
//     - spd_lte=
//     - spd_gte=
//     - hp_lte=
//     - hp_gte=
func (c ChallengeServiceImpl) GenerateUniqueNgrokChallenge(memberID uuid.UUID) NgrokChallenge {
	rng := utils.CreateRNGFromHash(memberID)
	aliens := GenerateNgrokAliens(rng, memberID)

	requests := []NgrokRequest{
		NgrokDeleteRequest{
			Name:   "POST all alien",
			Points: 0,
			Path:   "/api/aliens",
		},

		NgrokPostRequest{
			Name:   "POST all alien",
			Points: NGROK_POST_POINTS,
			Path:   "/api/aliens",
			Body:   slices.Clone(aliens),
		},

		NgrokGetRequest{
			Name:           "GET all aliens",
			Points:         NGROK_GET_ALL_POINTS,
			Path:           "/api/aliens",
			ExpectedAliens: slices.Clone(aliens),
		},
	}

	requests = append(requests, GenerateRandomFilterTests(rng, slices.Clone(aliens))...)

	return NgrokChallenge{Requests: requests}
}

func GenerateNgrokAliens(rng *rand.Rand, memberID uuid.UUID) []DetailedAlien {
	// Use challenge ID as seed for deterministic but unique data.
	count := utils.GenerateRandomNumWithinRange(rng, NUM_NGROK_ALIENS_LOWER_BOUND, NUM_NGROK_ALIENS_UPPER_BOUND)
	aliens := []DetailedAlien{}
	for alienIdx := range count {
		alien := GenerateDetailedAlien(rng, memberID, alienIdx)
		aliens = append(aliens, alien)
	}

	return aliens
}
