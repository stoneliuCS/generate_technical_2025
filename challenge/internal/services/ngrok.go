package services

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

type NgrokChallenge struct {
	Requests []NgrokRequest
}

type NgrokChallengeScore struct {
	Valid  bool
	Score  int
	Reason string // optional, only set when Valid = false
}

type NgrokRequest interface {
	Execute(client *http.Client, baseURL string) error
}

type NgrokPostRequest struct {
	Name   string
	Points int
	Path   string
	Body   []DetailedAlien
}

type NgrokGetRequest struct {
	Name           string
	Points         int
	Path           string
	ExpectedAliens []DetailedAlien
}

type Stats struct {
	Atk int `json:"atk"`
	HP  int `json:"hp"`
}

func (t NgrokPostRequest) Execute(client *http.Client, baseURL string) error {
	return nil
}
func (t NgrokGetRequest) Execute(client *http.Client, baseURL string) error {
	return nil
}

func generateRandomFilterTests(rng *rand.Rand, aliens []DetailedAlien) []NgrokRequest {
	requests := []NgrokRequest{}
	// Filter by type:
	randomType := alienTypes[rng.Intn(len(alienTypes))]

	var queryParamsType []string
	queryParamsType = append(queryParamsType, fmt.Sprintf("type=%s", randomType))

	expectedAliensType := filterAliensByType(aliens, randomType)

	request := NgrokGetRequest{
		Name:           fmt.Sprintf("Filter by type=%s", randomType),
		Points:         NGROK_FILTER_TYPE_POINTS,
		Path:           "/api/aliens?" + strings.Join(queryParamsType, "&"),
		ExpectedAliens: expectedAliensType,
	}

	requests = append(requests, request)

	// Filter by SPD
	minSPD := rng.Intn(ALIEN_ATK_HP_SPD_UPPER + 1)

	var queryParamsSPD []string
	queryParamsSPD = append(queryParamsSPD, fmt.Sprintf("spd_gte=%s", randomType))

	expectedAliensSPD := filterAliensByMinSPD(aliens, minSPD)

	request = NgrokGetRequest{
		Name:           fmt.Sprintf("Filter by spd>=%s", randomType),
		Points:         NGROK_FILTER_SPD_POINTS,
		Path:           "/api/aliens?" + strings.Join(queryParamsSPD, "&"),
		ExpectedAliens: expectedAliensSPD,
	}

	requests = append(requests, request)

	var queryParamsATK []string
	queryParamsATK = append(queryParamsATK, fmt.Sprintf("atk_gte=%s", 3))
	queryParamsATK = append(queryParamsATK, fmt.Sprintf("atk_lt=%s", 2))

	// The expected aliens is an empty array--the filters contradict each other.
	expectedAliensATK := []DetailedAlien{}

	request = NgrokGetRequest{
		Name:           "Filter by atk contradict",
		Points:         NGROK_FILTER_SPD_POINTS,
		Path:           "/api/aliens?" + strings.Join(queryParamsSPD, "&"),
		ExpectedAliens: expectedAliensATK,
	}

	requests = append(requests, request)

	return requests
}

// Filter Helpers

func filterAliensByType(aliens []DetailedAlien, alienType AlienType) []DetailedAlien {
	var filtered []DetailedAlien
	for _, alien := range aliens {
		if alien.Type == alienType {
			filtered = append(filtered, alien)
		}
	}
	return filtered
}

func filterAliensByMinSPD(aliens []DetailedAlien, minSPD int) []DetailedAlien {
	var filtered []DetailedAlien
	for _, alien := range aliens {
		if alien.Spd >= minSPD {
			filtered = append(filtered, alien)
		}
	}
	return filtered
}
