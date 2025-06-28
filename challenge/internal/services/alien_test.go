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
	i := 0
	for i < 10 {
		// No Weapons can be purchased from 0 to 9
		weapons := services.GenerateAllPossibleWeaponPurchasesFromBudget(i)
		assert.True(t, len(weapons) == 0)
		i = i + 1
	}

	i = 10
	for i < 20 {
		weapons := services.GenerateAllPossibleWeaponPurchasesFromBudget(i)
		// Possibly 2 purchases, 1 or 2 turrets 
		assert.True(t, len(weapons) == 1 || len(weapons) == 2)
		i = i + 1
	}

	i = 20
	for i < 30 {
		weapons := services.GenerateAllPossibleWeaponPurchasesFromBudget(i)
		// Possibly 3 purchaes, 1, 2, or 3 turrets
		assert.True(t, len(weapons) == 2 || len(weapons) == 1 || len(weapons) == 3)
		i = i + 1
	}

	i = 30
	for i < 40 {
		weapons := services.GenerateAllPossibleWeaponPurchasesFromBudget(i)
		// Possibly 4 purchases, 1 2,3 or 4 turrets. Or 1 Machine Gun and 1 turret
		assert.True(t, len(weapons) == 2 || len(weapons) == 1 || len(weapons) == 3 || len(weapons) == 4)
		i = i + 1
	}

	weapons := services.GenerateAllPossibleWeaponPurchasesFromBudget(100)
	for _, weapons := range weapons {
		totalCost := lo.Reduce(weapons, func(acc int, weapon services.Weapon, _ int) int {
			return acc + int(weapon.Cost)
		}, 0)
		assert.LessOrEqual(t, totalCost, 100)
	}
}
