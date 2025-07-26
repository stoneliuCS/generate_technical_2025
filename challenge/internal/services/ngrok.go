package services

import (
	"bytes"
	"context"
	"encoding/json"
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
	Execute(client *http.Client, baseURL string) (pointsEarned int, err error)
	GetName() string
}

type NgrokDeleteRequest struct {
	Name   string
	Points int
	Path   string
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

func (t NgrokDeleteRequest) GetName() string {
	return t.Name
}

func (t NgrokGetRequest) GetName() string {
	return t.Name
}

func (t NgrokPostRequest) GetName() string {
	return t.Name
}

type Stats struct {
	Atk int `json:"atk"`
	HP  int `json:"hp"`
}

func (t NgrokDeleteRequest) Execute(client *http.Client, baseURL string) (int, error) {
	req, err := http.NewRequest("DELETE", baseURL+t.Path, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("ngrok-skip-browser-warning", "true")

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("expected status 204 or 200, got %d", resp.StatusCode)
	}

	return 0, nil
}

func (t NgrokPostRequest) Execute(client *http.Client, baseURL string) (int, error) {
	// Marshal body to JSON
	bodyBytes, err := json.Marshal(t.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", baseURL+t.Path, bytes.NewReader(bodyBytes))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("ngrok-skip-browser-warning", "true")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code (expect 201 Created for POST)
	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	return NGROK_POST_POINTS, nil
}

func (t NgrokGetRequest) Execute(client *http.Client, baseURL string) (int, error) {
	// Create request
	req, err := http.NewRequest("GET", baseURL+t.Path, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("ngrok-skip-browser-warning", "true")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Parse response
	var actualAliens []DetailedAlien
	if err := json.NewDecoder(resp.Body).Decode(&actualAliens); err != nil {
		return 0, fmt.Errorf("invalid JSON response: %w", err)
	}

	// Calculate distance and adjust points.
	distance := calculateAlienDistance(t.ExpectedAliens, actualAliens)

	if distance == 0 {
		// Perfect match.
		return t.Points, nil
	} else {
		// Partial credit.
		adjustedPoints := t.Points - distance
		if adjustedPoints < 0 {
			adjustedPoints = 0
		}

		fmt.Printf("distance_error:%d:expected %d aliens, got %d aliens with %d differences\n",
			distance, len(t.ExpectedAliens), len(actualAliens), distance)
		return adjustedPoints, nil
	}
}

func calculateAlienDistance(expected, actual []DetailedAlien) int {
	distance := 0

	lengthDiff := abs(len(expected) - len(actual))
	distance += lengthDiff

	expectedMap := make(map[string]DetailedAlien)
	actualMap := make(map[string]DetailedAlien)

	for _, alien := range expected {
		expectedMap[alien.ID] = alien
	}
	for _, alien := range actual {
		actualMap[alien.ID] = alien
	}

	contentDiffs := 0

	// Check aliens that exist in both sets.
	for _, expectedAlien := range expected {
		if actualAlien, exists := actualMap[expectedAlien.ID]; exists {
			if !aliensEqual(expectedAlien, actualAlien) {
				contentDiffs++
			}
		}
	}

	distance += contentDiffs

	return distance
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func aliensEqual(a, b DetailedAlien) bool {
	return a.ID == b.ID &&
		a.FirstName == b.FirstName &&
		a.LastName == b.LastName &&
		a.Type == b.Type &&
		a.Spd == b.Spd &&
		a.ProfileURL == b.ProfileURL
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

func health(ctx context.Context, url string) (ok bool, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/health", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("ngrok-skip-browser-warning", "true")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}
