package services_test

import (
	"generate_technical_challenge_2025/internal/services"
	"generate_technical_challenge_2025/internal/utils"
	// "slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNgrokAliens(t *testing.T) {
	firstRNG := utils.CreateRNGFromHash(UUID)
	aliens := services.GenerateNgrokAliens(firstRNG, UUID)
	assert.True(t, len(aliens) <= services.NUM_NGROK_ALIENS_UPPER_BOUND)
	assert.True(t, len(aliens) >= services.NUM_NGROK_ALIENS_LOWER_BOUND)

	// Aliens generated from the same UUID should be the same in length and content.
	duplicateRng := utils.CreateRNGFromHash(UUID)
	aliensAgain := services.GenerateNgrokAliens(duplicateRng, UUID)
	assert.Equal(t, aliens, aliensAgain)
}

func TestGenerateRandomFilterTests(t *testing.T) {
	aliens := services.GenerateNgrokAliens(RNG, UUID)

	requests := services.GenerateRandomFilterTests(RNG, aliens)
	assert.Len(t, requests, 5)
}

// func TestCalculateAlienDistance(t *testing.T) {
// 	firstRNG := utils.CreateRNGFromHash(UUID)
// 	aliens := services.GenerateNgrokAliens(firstRNG, UUID)
// 	// Aliens are 0 distance from themselves.
// 	dist := services.CalculateAlienDistance(aliens, aliens)
// 	assert.True(t, dist == 0)
//
// 	distWithEmpty := services.CalculateAlienDistance(aliens, []services.DetailedAlien{})
// 	assert.True(t, distWithEmpty == len(aliens))
//
// 	distWithOneMissing := services.CalculateAlienDistance(aliens, aliens[:len(aliens)-1])
// 	assert.True(t, distWithOneMissing == 1)
//
// 	// Add an extra element.
// 	extraAlien := services.GenerateDetailedAlien(firstRNG, UUID, 3)
// 	longer := make([]services.DetailedAlien, 0, len(aliens)+1)
// 	longer = append(longer, extraAlien)
// 	longer = append(longer, aliens...)
// 	distWithOneAdded := services.CalculateAlienDistance(aliens, longer)
// 	assert.True(t, distWithOneAdded == 1)
//
// 	// Change the first element to have a different SPD value.
// 	aliensToBeModified := slices.Clone(aliens)
// 	tmp := aliensToBeModified[0]
// 	tmp.Spd = min(1, services.ALIEN_ATK_HP_SPD_UPPER%(tmp.Spd+1))
// 	aliensToBeModified[0] = tmp
//
// 	distWithOneChange := services.CalculateAlienDistance(aliens, aliensToBeModified)
// 	assert.True(t, distWithOneChange == 1)
//
// 	// With two changes to a single alien, don't double-count it.
// 	tmp.BaseAlien.Atk = min(1, services.ALIEN_ATK_HP_SPD_UPPER%(tmp.BaseAlien.Atk+1))
// 	aliensToBeModified[0] = tmp
// 	distWithTwoChanges := services.CalculateAlienDistance(aliens, aliensToBeModified)
// 	assert.True(t, distWithTwoChanges == 1)
//
// 	// Changing the ID on an alien:
// 	tmp.ID = tmp.ID + "1" // Guaranteed to change the ID--not a realistic ID but definitely different.
// 	aliensToBeModified[0] = tmp
// 	distWithIDChange := services.CalculateAlienDistance(aliens, aliensToBeModified)
// 	assert.True(t, distWithIDChange == 1)
//
// 	allTheSameAlien := make([]services.DetailedAlien, len(aliens))
// 	for idx := range aliens {
// 		allTheSameAlien[idx] = aliens[0]
// 	}
// 	distWithAllTheSameAlien := services.CalculateAlienDistance(aliens, allTheSameAlien)
// 	assert.True(t, distWithAllTheSameAlien == len(aliens)-1)
// }
