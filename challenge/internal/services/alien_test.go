package services_test

import (
	"generate_technical_challenge_2025/internal/services"
	"generate_technical_challenge_2025/internal/utils"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

const STARTING_HP = 1000

var (
	UUID = uuid.New()
	RNG  = utils.CreateRNGFromHash(UUID)
)

func TestGenerateAlienInvasion(t *testing.T) {
	sampleAlienInvasion := services.GenerateAlienInvasion(RNG)
	actualSizeOfAlienInvasion := len(sampleAlienInvasion)
	// Assert that the size of the alien invasion must always be within the bounds of
	// the generated bounds
	assert.True(t, actualSizeOfAlienInvasion >= services.LOWER_ALIEN_AMOUNT)
	assert.True(t, actualSizeOfAlienInvasion <= services.UPPER_ALIEN_AMOUNT)
	// Assert that for each alien generated, their HP and ATTACK are always within the bounds.
	withinBounds := lo.Reduce(sampleAlienInvasion, func(flag bool, alien services.Alien, _ int) bool {
		return flag &&
			alien.Atk >= services.ALIEN_ATK_HP_LOWER &&
			alien.Atk <= services.ALIEN_ATK_HP_UPPER &&
			alien.Hp >= services.ALIEN_ATK_HP_LOWER &&
			alien.Hp <= services.ALIEN_ATK_HP_UPPER
	}, true)
	assert.True(t, withinBounds)
}

func TestAlienInvasionIsOver(t *testing.T) {
	// Case one there are aliens left but the HP is gone.
	case1 := services.CreateInvasionState([]services.Alien{
		// One alien is left.
		services.CreateAlien(3, 3),
	}, 0)
	assert.True(t, case1.IsOver())
	// Case two there are no aliens left and Hp is still above 0
	case2 := services.CreateInvasionState([]services.Alien{
		// One alien is left.
	}, 1)
	assert.True(t, case2.IsOver())
	// Case three there are aliens left and the HP is overkilled to negative
	case3 := services.CreateInvasionState([]services.Alien{
		// One alien is left.
		services.CreateAlien(3, 3),
	}, -100)
	assert.True(t, case3.IsOver())
}

func TestAlienInvasionVolley(t *testing.T) {
	exampleAlienInvasion := services.CreateInvasionState([]services.Alien{
		services.CreateAlien(1, 2),
		services.CreateAlien(2, 2),
		services.CreateAlien(3, 3),
	}, 100)

	assert.True(t, exampleAlienInvasion.GetAliensLeft() == 3)
	actualAlienStateAfterVolley := exampleAlienInvasion.AttackAllAliens()
	// The first alien must have died so it should be reflected in the final aliensLeft count
	assert.True(t, actualAlienStateAfterVolley.GetAliensLeft() == 2)
	assert.True(t, !actualAlienStateAfterVolley.IsOver())
}

func TestAlienInvasionFocusedShot(t *testing.T) {
	exampleAlienInvasion := services.CreateInvasionState([]services.Alien{
		services.CreateAlien(1, 2),
		services.CreateAlien(2, 2),
		services.CreateAlien(3, 3),
	}, 100)
	// Current highest damaging alien has attack of 3
	assert.True(t, exampleAlienInvasion.GetCurrentHighestDamagingAlien().Atk == 3)
	actualAliensAfterFocusedShot := exampleAlienInvasion.AttackHighestDamageAlien()
	assert.True(t, actualAliensAfterFocusedShot.GetAliensLeft() == 2)
	// Now it should be 2
	assert.True(t, actualAliensAfterFocusedShot.GetCurrentHighestDamagingAlien().Atk == 2)
	assert.True(t, !actualAliensAfterFocusedShot.IsOver())
}

func TestAlienInvasionFocusedVolley(t *testing.T) {
	a1 := services.CreateAlien(1, 2)
	a2 := services.CreateAlien(2, 2)
	a3 := services.CreateAlien(2, 2)
	a4 := services.CreateAlien(3, 3)
	exampleAlienInvasion := services.CreateInvasionState([]services.Alien{
		a1,
		a2,
		a3,
		a4,
	}, 100)
	actualAliensAfterFocusedVolley := exampleAlienInvasion.AttackHighestDamagingHalf()
	// It should be the case that 1 alien has died, and the other alien the 3 by 3 one has lost 2 hp
	assert.True(t, actualAliensAfterFocusedVolley.GetAliensLeft() == 3)
	assert.True(t, actualAliensAfterFocusedVolley.GetCurrentHighestDamagingAlien().Atk == 3)
	// When we view the aliens, they should be sorted in order because of how our invasion state manages aliens internally.
	assert.True(t, slices.EqualFunc(
		actualAliensAfterFocusedVolley.SurveyRemainingAlienInvasion(),
		[]services.Alien{
			a4.TakeDamage(2),
			a2,
			a1,
		},
		func(alien1, alien2 services.Alien) bool {
			return alien1.Hp == alien2.Hp && alien1.Atk == alien2.Atk
		},
	))
}

// BEGIN ALGORITHM TESTING

func TestAlgorithm(t *testing.T) {
	sampleAlienInvasion := services.GenerateAlienInvasion(RNG)
	// A State with only 100 HP should end in at most 10 rounds if generating atleast 10 aliens.
	services.RunAllPossibleInvasionStatesToCompletion(services.CreateInvasionState(sampleAlienInvasion, 500))
}
