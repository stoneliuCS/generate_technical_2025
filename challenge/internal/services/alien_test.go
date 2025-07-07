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
