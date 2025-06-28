package services_test

import (
	"generate_technical_challenge_2025/internal/services"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	SERVICE_CHALLENGE = services.CreateChallengeService(nil, nil)
	UUID              = uuid.New()
)

func TestGenerateAlienInvasion(t *testing.T) {
	invasionState := services.GenerateInvasionState(UUID)
	// Test that the waves generate the proper number of aliens
	assert.True(t, len(invasionState.Waves) == 3)
	assert.True(t, len(invasionState.Waves[0]) == 5)
	assert.True(t, len(invasionState.Waves[1]) == 10)
	assert.True(t, len(invasionState.Waves[2]) == 20)

	invasionState2 := services.GenerateInvasionState(UUID)
	// Assert that the same uuid will always generate the same invasion state.
	assert.Equal(t, invasionState, invasionState2)
}
