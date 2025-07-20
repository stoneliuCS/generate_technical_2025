package services_test

import (
	"generate_technical_challenge_2025/internal/data"
	"generate_technical_challenge_2025/internal/services"
	"generate_technical_challenge_2025/internal/transactions"
	"log/slog"
	"math/rand"
	"net/url"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var alienTypes = []services.AlienType{
	services.AlienTypeRegular,
	services.AlienTypeElite,
	services.AlienTypeBoss,
}

// Helper for testing that a DetailedAlien is valid.
func assertValidAlien(t *testing.T, sampleAlien services.DetailedAlien) {
	assert.True(t, sampleAlien.BaseAlien.Atk >= services.ALIEN_ATK_HP_SPD_LOWER)
	assert.True(t, sampleAlien.BaseAlien.Atk <= services.ALIEN_ATK_HP_SPD_UPPER)

	validType := false

	for _, alienType := range alienTypes {
		if sampleAlien.Type == alienType {
			validType = true
		}
	}
	assert.True(t, validType)

	// Valid, parseable profile photo URL.
	_, err := url.Parse(sampleAlien.ProfileURL)
	assert.True(t, err == nil)

	// Valid, non-nil uuid ID.
	assert.True(t, sampleAlien.ID != uuid.Nil)

	assert.True(t, slices.Contains(data.AlienNames, sampleAlien.Name))
}

func TestGenerateDetailedAlien(t *testing.T) {
	alienIndex := rand.Int()
	sampleAlien := services.GenerateDetailedAlien(RNG, UUID, alienIndex)
	assertValidAlien(t, sampleAlien)
}

var (
	LOGGER                 = slog.New(slog.Default().Handler())
	CHALLENGE_SERVICE_IMPL = services.CreateChallengeService(LOGGER,
		transactions.CreateChallengeTransactions(LOGGER, nil)) // nil DB for testing.
)

func TestGenerateUniqueFrontendChallenge(t *testing.T) {
	firstAliens := CHALLENGE_SERVICE_IMPL.GenerateUniqueFrontendChallenge(UUID)
	secondAliens := CHALLENGE_SERVICE_IMPL.GenerateUniqueFrontendChallenge(UUID)

	assert.Equal(t, firstAliens, secondAliens)

	for _, alien := range firstAliens {
		assertValidAlien(t, alien)
	}
}
