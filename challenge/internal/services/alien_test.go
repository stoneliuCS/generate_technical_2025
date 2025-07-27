package services_test

import (
	"generate_technical_challenge_2025/internal/services"
	"generate_technical_challenge_2025/internal/utils"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

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
			alien.Atk >= services.ALIEN_ATK_HP_SPD_LOWER &&
			alien.Atk <= services.ALIEN_ATK_HP_SPD_UPPER &&
			alien.Hp >= services.ALIEN_ATK_HP_SPD_LOWER &&
			alien.Hp <= services.ALIEN_ATK_HP_SPD_UPPER
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
	actualAlienStateAfterVolley := exampleAlienInvasion.AttackAliensModulo()
	assert.True(t, actualAlienStateAfterVolley.GetAliensLeft() == 3)
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

func TestAlienInvasionCeiling(t *testing.T) {
	a1 := services.CreateAlien(2, 2)
	a2 := services.CreateAlien(2, 2)
	a3 := services.CreateAlien(2, 2)
	a4 := services.CreateAlien(2, 3)
	a5 := services.CreateAlien(2, 3)
	exampleAlienInvasion := services.CreateInvasionState([]services.Alien{
		a1,
		a2,
		a3,
		a4,
		a5,
	}, 100)
	actualAliensAfterFocusedVolley := exampleAlienInvasion.AttackHighestDamagingHalf()
	// Should kill 3 aliens, 5 - 3 = 2 because 5/2 = 2.5 ceiling is 3.
	assert.Equal(t, actualAliensAfterFocusedVolley.GetAliensLeft(), 2)
}

// BEGIN ALGORITHM TESTING

func TestAlgorithmTimes(t *testing.T) {
	sampleAlienInvasion := services.GenerateAlienInvasion(RNG)
	done := make(chan bool)

	go func() {
		services.RunAllPossibleInvasionStatesToCompletionGreedy(services.CreateInvasionState(sampleAlienInvasion, 100))
		done <- true
	}()

	select {
	case <-done:
		// Function completed in time
		t.Log("Algorithm completed within the allotted time.")
		assert.True(t, true, "Successfully completed the algorithm within allotted time.")
	case <-time.After(1 * time.Second):
		t.Fatal("Algorithm did not complete within 1 seconds")
		assert.True(t, false)
	}
}

func TestAlgorithmCorrectness(t *testing.T) {
	sampleAlienInvasion := services.GenerateAlienInvasion(RNG)
	sampleInvasionState := services.CreateInvasionState(sampleAlienInvasion, 100)
	greedySol := services.RunAllPossibleInvasionStatesToCompletionGreedy(sampleInvasionState)
	bruteforceSol := services.RunAllPossibleInvasionStatesToCompletion(sampleInvasionState)
	// The greedy solution should always have the least number of aliens left over.
	bruteForceBestSolByAliensLeft := lo.MinBy(bruteforceSol, func(s1 services.InvasionState, s2 services.InvasionState) bool {
		return s1.GetAliensLeft() < s2.GetAliensLeft()
	})
	greedySolByAliensLeft := lo.MinBy(greedySol, func(s1 services.InvasionState, s2 services.InvasionState) bool {
		return s1.GetAliensLeft() < s2.GetAliensLeft()
	})
	allSolutionsBruteForceWithMinimalAliensLeft := lo.Filter(bruteforceSol, func(state services.InvasionState, idx int) bool {
		return state.GetAliensLeft() == bruteForceBestSolByAliensLeft.GetAliensLeft()
	})
	allSolutionsGreedyWithMinimalAliensLeft := lo.Filter(greedySol, func(state services.InvasionState, idx int) bool {
		return state.GetAliensLeft() == greedySolByAliensLeft.GetAliensLeft()
	})
	bruteForceBestSolByHP := lo.MaxBy(allSolutionsBruteForceWithMinimalAliensLeft, func(s1 services.InvasionState, s2 services.InvasionState) bool {
		return s1.GetHpLeft() > s2.GetHpLeft()
	})
	greedyBruteForceBestSolByHP := lo.MaxBy(allSolutionsGreedyWithMinimalAliensLeft, func(s1 services.InvasionState, s2 services.InvasionState) bool {
		return s1.GetHpLeft() > s2.GetHpLeft()
	})
	// Assert that the greedySol has the same minimal aliens left as the brute force.
	assert.LessOrEqual(t, greedyBruteForceBestSolByHP.GetAliensLeft(), bruteForceBestSolByHP.GetAliensLeft())
	// Assert that the greedySol has the same amount of hp as the brute force.
	assert.GreaterOrEqual(t, greedyBruteForceBestSolByHP.GetHpLeft(), bruteForceBestSolByHP.GetHpLeft())
}
