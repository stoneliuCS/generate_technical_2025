package services_test

import (
	"generate_technical_challenge_2025/internal/services"
	"generate_technical_challenge_2025/internal/utils"
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
	assert.Len(t, requests, 3)

}
