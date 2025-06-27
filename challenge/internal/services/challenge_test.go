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
	invasionState := SERVICE_CHALLENGE.GenerateUniqueChallenge(UUID)
	// First test that all default aliens are there
	assert.True(t, len(invasionState.DefaultAliens) == 3)
	assert.ElementsMatch(t, invasionState.DefaultAliens, []services.Alien{services.CreateRegularAlien(), services.CreateSwiftAlien(), services.CreateBossAlien()})
	// Next test that all default weapons are there
	assert.True(t, len(invasionState.DefaultWeapons) == 3)
	assert.ElementsMatch(t, invasionState.DefaultWeapons, []services.Weapon{services.CreateTurretWeapon(), services.CreateMachineGunWeapon(), services.CreateRayGunWeapon()})
	// Test that the waves generate the proper number of aliens
	assert.True(t, len(invasionState.Waves) == 3)
	assert.True(t, len(invasionState.Waves[0]) == 5)
	assert.True(t, len(invasionState.Waves[1]) == 10)
	assert.True(t, len(invasionState.Waves[2]) == 20)

	invasionState2 := SERVICE_CHALLENGE.GenerateUniqueChallenge(UUID)
	// Assert that the same uuid will always generate the same invasion state.
	assert.True(t, assert.ObjectsAreEqualValues(invasionState, invasionState2))
}
