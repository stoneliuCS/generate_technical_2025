package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
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
	Execute(ctx context.Context, client *http.Client, baseURL string) (pointsEarned int, err error)
	GetName() string
	GetTotalPossiblePoints() int
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

func (t NgrokDeleteRequest) GetTotalPossiblePoints() int {
	return t.Points
}

func (t NgrokGetRequest) GetTotalPossiblePoints() int {
	return t.Points
}

func (t NgrokPostRequest) GetTotalPossiblePoints() int {
	return t.Points
}

type Stats struct {
	Atk int `json:"atk"`
	HP  int `json:"hp"`
}

const (
	SPD = "spd"
	ATK = "atk"
	HP  = "hp"
)

func (t NgrokDeleteRequest) Execute(ctx context.Context, client *http.Client, baseURL string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "DELETE", baseURL+t.Path, nil)
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

func (t NgrokPostRequest) Execute(ctx context.Context, client *http.Client, baseURL string) (int, error) {
	// Marshal body to JSON
	bodyBytes, err := json.Marshal(t.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+t.Path, bytes.NewReader(bodyBytes))
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

func (t NgrokGetRequest) Execute(ctx context.Context, client *http.Client, baseURL string) (int, error) {
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+t.Path, nil)
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
	distance := CalculateAlienDistance(t.ExpectedAliens, actualAliens)

	if distance == 0 {
		// Perfect match.
		return t.Points, nil
	} else {
		// Partial credit.
		adjustedPoints := max(t.Points-distance, 0)

		if VERBOSE {
			fmt.Printf("distance_error:%d:expected %d aliens, got %d aliens with %d differences\n",
				distance, len(t.ExpectedAliens), len(actualAliens), distance)
		}
		return adjustedPoints, nil
	}
}

// Returns the 'distance' between the expected and actual alien sets.
// This is a positive score that is to be subtracted from some total possible score.
// A distance is calculated:
// - For each alien in expected, if there exists an alien in actual with the same exact
// same values, then those two aliens have a distance of 0.
// - If there exists an alien with the same ID but any number of other differing values,
// then those two aliens have a distance of 1.
// - Any alien in actual but not in expected has a distance of 1.
// - There is no double-counting of distance--if there is an alien in actual with the same ID as one in
// expected but different Atk and Spd, it will still only have a distance of 1.
func CalculateAlienDistance(expected, actual []DetailedAlien) int {
	distance := 0

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

			if VERBOSE {
				fmt.Println()
				fmt.Printf("Expected alien: %+v\n", expectedAlien)
				fmt.Printf("Actual alien: %+v\n", actualAlien)
				fmt.Println()
			}

			if !aliensEqual(expectedAlien, actualAlien) {
				contentDiffs++
			}
		} else {
			// Key/value pair doesn't exist in the actual set.
			// Therefore, either:
			// - The candidate doesn't have this alien at all.
			// - The candidate does have this alien but with a different ID.
			contentDiffs++
			if VERBOSE {
				fmt.Printf("Not found: expected: %+v actual: %+v\n", expectedAlien, actualAlien)
			}
		}
	}

	lengthDiff := max(0, len(actual)-len(expected))
	distance += lengthDiff

	distance += contentDiffs

	return distance
}

func aliensEqual(a, b DetailedAlien) bool {
	return a.ID == b.ID &&
		a.FirstName == b.FirstName &&
		a.LastName == b.LastName &&
		a.Type == b.Type &&
		a.Spd == b.Spd &&
		a.ProfileURL == b.ProfileURL
}

type FilterFunc func([]DetailedAlien, string) []DetailedAlien

var supportedFilters = map[string]FilterFunc{
	"type": filterAliensByType,
	"atk_lte": func(aliens []DetailedAlien, value string) []DetailedAlien {
		return filterAliensByNumericField(aliens, ATK, LTE, value)
	},
	"atk_gte": func(aliens []DetailedAlien, value string) []DetailedAlien {
		return filterAliensByNumericField(aliens, ATK, GTE, value)
	},
	"spd_lte": func(aliens []DetailedAlien, value string) []DetailedAlien {
		return filterAliensByNumericField(aliens, SPD, LTE, value)
	},
	"spd_gte": func(aliens []DetailedAlien, value string) []DetailedAlien {
		return filterAliensByNumericField(aliens, SPD, GTE, value)
	},
	"hp_lte": func(aliens []DetailedAlien, value string) []DetailedAlien {
		return filterAliensByNumericField(aliens, HP, LTE, value)
	},
	"hp_gte": func(aliens []DetailedAlien, value string) []DetailedAlien {
		return filterAliensByNumericField(aliens, HP, GTE, value)
	},
}

type ComparisonOp string

const (
	LTE ComparisonOp = "lte"
	GTE ComparisonOp = "gte"
)

// Generates:
// - GET filter by a random type.
// - GET filter by random max/min SPD.
// - GET filter by random max/min ATK.
// - GET filter by random max/min HP.
// - GET filter by contradicting ATK/SPD/HP (e.g. ?atk_lte=2&atkgte=3).
func GenerateRandomFilterTests(rng *rand.Rand, aliens []DetailedAlien) []NgrokRequest {
	requests := []NgrokRequest{}

	// Test 1: Filter by random type
	requests = append(requests, generateTypeFilterTest(rng, aliens))

	// Test 2: Filter by random SPD (gte or lte)
	requests = append(requests, generateNumericFilterTest(rng, aliens, SPD, NGROK_FILTER_SPD_POINTS))

	// Test 3: Filter by random ATK (gte or lte)
	requests = append(requests, generateNumericFilterTest(rng, aliens, ATK, NGROK_FILTER_ATK_POINTS))

	// Test 4: Filter by random HP (gte or lte)
	requests = append(requests, generateNumericFilterTest(rng, aliens, HP, NGROK_FILTER_HP_POINTS))

	// Test 5: Filter by contradicting filters (randomly pick ATK, SPD, or HP)
	fields := []string{ATK, SPD, HP}
	contradictField := fields[rng.Intn(len(fields))]
	requests = append(requests, generateContradictingFilterTest(rng, aliens, contradictField))

	return requests
}

func generateTypeFilterTest(rng *rand.Rand, aliens []DetailedAlien) NgrokRequest {
	randomType := alienTypes[rng.Intn(len(alienTypes))]

	queryParams := []string{fmt.Sprintf("type=%s", randomType)}
	expectedAliens := applyFilters(aliens, map[string]string{"type": string(randomType)})

	return NgrokGetRequest{
		Name:           fmt.Sprintf("Filter by type=%s", randomType),
		Points:         NGROK_FILTER_TYPE_POINTS,
		Path:           NGROK_PATH + "?" + strings.Join(queryParams, "&"),
		ExpectedAliens: expectedAliens,
	}
}

func generateNumericFilterTest(rng *rand.Rand, aliens []DetailedAlien, field string, points int) NgrokRequest {
	isGte := rng.Intn(2) == 0
	value := rng.Intn(ALIEN_ATK_HP_SPD_UPPER)

	var filterKey, queryParam, description string
	if isGte {
		filterKey = field + "_gte"
		queryParam = fmt.Sprintf("%s_gte=%d", field, value)
		description = fmt.Sprintf("Filter by %s>=%d", field, value)
	} else {
		filterKey = field + "_lte"
		queryParam = fmt.Sprintf("%s_lte=%d", field, value)
		description = fmt.Sprintf("Filter by %s<=%d", field, value)
	}

	queryParams := []string{queryParam}
	expectedAliens := applyFilters(aliens, map[string]string{filterKey: strconv.Itoa(value)})

	return NgrokGetRequest{
		Name:           description,
		Points:         points,
		Path:           NGROK_PATH + "?" + strings.Join(queryParams, "&"),
		ExpectedAliens: expectedAliens,
	}
}

func generateContradictingFilterTest(rng *rand.Rand, aliens []DetailedAlien, field string) NgrokRequest {
	// Generate contradicting values: lte < gte
	lteValue := rng.Intn(ALIEN_ATK_HP_SPD_UPPER / 2)              // Lower half
	gteValue := lteValue + rng.Intn(ALIEN_ATK_HP_SPD_UPPER/2) + 1 // Higher value

	queryParams := []string{
		fmt.Sprintf("%s_lte=%d", field, lteValue),
		fmt.Sprintf("%s_gte=%d", field, gteValue),
	}

	// Apply contradicting filters--should result in empty array.
	filters := map[string]string{
		field + "_lte": strconv.Itoa(lteValue),
		field + "_gte": strconv.Itoa(gteValue),
	}
	expectedAliens := applyFilters(aliens, filters)

	return NgrokGetRequest{
		Name:           fmt.Sprintf("Filter by %s contradict (lte=%d, gte=%d)", field, lteValue, gteValue),
		Points:         NGROK_FILTER_CONTRADICT,
		Path:           NGROK_PATH + "?" + strings.Join(queryParams, "&"),
		ExpectedAliens: expectedAliens,
	}
}

// Filter Helpers

func applyFilters(aliens []DetailedAlien, filters map[string]string) []DetailedAlien {
	result := aliens

	for filterKey, value := range filters {
		if filterFunc, exists := supportedFilters[filterKey]; exists {
			result = filterFunc(result, value)
		}
	}

	return result
}

func filterAliensByType(aliens []DetailedAlien, alienType string) []DetailedAlien {
	var filtered []DetailedAlien
	for _, alien := range aliens {
		if string(alien.Type) == alienType {
			filtered = append(filtered, alien)
		}
	}
	return filtered
}

func filterAliensByNumericField(aliens []DetailedAlien, field string, op ComparisonOp, value string) []DetailedAlien {
	targetValue, _ := strconv.Atoi(value)
	var filtered []DetailedAlien

	for _, alien := range aliens {
		var fieldValue int

		// Get the field value:
		switch field {
		case ATK:
			fieldValue = alien.BaseAlien.Atk
		case SPD:
			fieldValue = alien.Spd
		case HP:
			fieldValue = alien.BaseAlien.Hp
		default:
			continue
		}

		// Apply comparison:
		var matches bool
		switch op {
		case LTE:
			matches = fieldValue <= targetValue
		case GTE:
			matches = fieldValue >= targetValue
		}

		if matches {
			filtered = append(filtered, alien)
		}
	}

	return filtered
}

func health(ctx context.Context, url string) (ok bool, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/healthcheck", nil)
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
