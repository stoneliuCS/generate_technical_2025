package services_test

import (
	"generate_technical_challenge_2025/internal/services"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var UUID = uuid.New()

func TestGenerateAlienInvasion(t *testing.T) {
	invasionState := services.GenerateInvasionState(UUID)
	invasionState2 := services.GenerateInvasionState(UUID)
	// Assert that the same uuid will always generate the same invasion state.
	assert.Equal(t, invasionState, invasionState2)
	assert.True(t, len(invasionState.Waves) == 3)
	assert.True(t, len(invasionState2.Waves) == 3)

	// Check that the generated aliens are indeed between the ranges.

	getAlienCountsPerWave := func(aliens []services.Alien, filterBy services.AlienType) int {
		return lo.Reduce(aliens, func(acc int, alien services.Alien, _ int) int {
			if alien.Type == filterBy {
				return acc + 1
			}
			return acc
		}, 0)
	}

	// TEST WAVE ONE
	numOfRegularAliens := getAlienCountsPerWave(invasionState.Waves[0], services.Regular)
	numOfSwiftAliens := getAlienCountsPerWave(invasionState.Waves[0], services.Swift)
	numOfBossAliens := getAlienCountsPerWave(invasionState.Waves[0], services.Boss)
	assert.True(t, numOfRegularAliens >= 3 && numOfRegularAliens <= 5)
	assert.True(t, numOfSwiftAliens >= 0 && numOfSwiftAliens <= 1)
	assert.True(t, numOfBossAliens == 0)

	// TEST WAVE TWO
	numOfRegularAliens = getAlienCountsPerWave(invasionState.Waves[1], services.Regular)
	numOfSwiftAliens = getAlienCountsPerWave(invasionState.Waves[1], services.Swift)
	numOfBossAliens = getAlienCountsPerWave(invasionState.Waves[1], services.Boss)
	assert.True(t, numOfRegularAliens >= 3 && numOfRegularAliens <= 5)
	assert.True(t, numOfSwiftAliens >= 3 && numOfSwiftAliens <= 5)
	assert.True(t, numOfBossAliens >= 0 && numOfBossAliens <= 1)

	// TEST WAVE THREE
	numOfRegularAliens = getAlienCountsPerWave(invasionState.Waves[2], services.Regular)
	numOfSwiftAliens = getAlienCountsPerWave(invasionState.Waves[2], services.Swift)
	numOfBossAliens = getAlienCountsPerWave(invasionState.Waves[2], services.Boss)
	assert.True(t, numOfRegularAliens >= 5 && numOfRegularAliens <= 10)
	assert.True(t, numOfSwiftAliens >= 5 && numOfSwiftAliens <= 7)
	assert.True(t, numOfBossAliens >= 1 && numOfBossAliens <= 3)
}

func TestGenerateAllPossibleWeaponPurchasesFromBudget(t *testing.T) {
	weapons := services.GenerateAllPossibleWeaponPurchasesFromBudget(100)
	for _, ws := range weapons {
		totalCost := lo.Reduce(ws, func(acc int, weapon services.Weapon, _ int) int {
			return acc + int(weapon.Cost)
		}, 0)
		assert.LessOrEqual(t, totalCost, 100)
	}

	// Budget of 10
	weapons = services.GenerateAllPossibleWeaponPurchasesFromBudget(10)
	assert.Len(t, weapons, 1)
	assert.Len(t, weapons[0], 1)
	assert.Equal(t, weapons[0][0].Type, services.Turret)

	// Budget of 20
	weapons = services.GenerateAllPossibleWeaponPurchasesFromBudget(20)
	assert.Len(t, weapons, 2)
}
