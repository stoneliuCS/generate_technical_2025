package services

import (
	"math/rand"
	"slices"

	"github.com/samber/lo"
)

// Represents the state of the invasion
type InvasionState struct {
	aliensLeft []Alien
	hpLeft     int
	commands   []string
}

func CreateInvasionState(aliens []Alien, startingHp int) InvasionState {
	// Sort all the aliens by attack power.
	sortedAliens := slices.SortedFunc(slices.Values(aliens), func(a1 Alien, a2 Alien) int {
		power1 := a1.Atk + a1.Hp
		power2 := a2.Atk + a2.Hp
		return power2 - power1 // Highest total power first
	})
	return InvasionState{
		aliensLeft: sortedAliens,
		hpLeft:     startingHp,
		commands:   []string{},
	}
}

func (i InvasionState) GetNumberOfCommandsUsed() int {
	return len(i.commands)
}

func (i InvasionState) SurveyRemainingAlienInvasion() []Alien {
	return slices.Clone(i.aliensLeft)
}

func (i InvasionState) GetAliensLeft() int {
	return len(i.aliensLeft)
}

func (i InvasionState) GetHpLeft() int {
	return i.hpLeft
}

func (i InvasionState) GetCurrentHighestDamagingAlien() Alien {
	if len(i.aliensLeft) == 0 {
		panic("No more aliens left.")
	}
	return i.aliensLeft[0]
}

// From all possible states, filter out the invasion states by these criteria
// The states with the least amount of aliens.
// The states with the highest hp left over.
// The states with the least number of commands left over.
func FilterToFindTheMostOptimalInvasions(states []InvasionState) []InvasionState {
	filterByMostAliensKilled := func(remainingStates []InvasionState) []InvasionState {
		leastNumAliensLeft := lo.MaxBy(remainingStates, func(s1 InvasionState, s2 InvasionState) bool {
			return s1.GetAliensLeft() < s2.GetAliensLeft()
		}).GetAliensLeft()

		filteredStatesByMostAliensKilled := lo.Filter(remainingStates, func(state InvasionState, _ int) bool {
			return state.GetAliensLeft() == leastNumAliensLeft
		})
		return filteredStatesByMostAliensKilled
	}

	filteredByHighestHpLeftOver := func(remainingStates []InvasionState) []InvasionState {
		highestHpLeft := lo.MaxBy(remainingStates, func(s1 InvasionState, s2 InvasionState) bool {
			return s1.hpLeft > s2.hpLeft
		}).GetHpLeft()

		filteredStatesByHighestHp := lo.Filter(remainingStates, func(state InvasionState, _ int) bool {
			return state.hpLeft == highestHpLeft
		})
		return filteredStatesByHighestHp
	}

	filteredByLeastNumberOfCommandsUsed := func(remainingStates []InvasionState) []InvasionState {
		leastNumberOfCommandsUsed := lo.MaxBy(remainingStates, func(s1 InvasionState, s2 InvasionState) bool {
			return s1.GetNumberOfCommandsUsed() < s2.GetNumberOfCommandsUsed()
		}).GetHpLeft()

		filteredStatesByLeastCommands := lo.Filter(remainingStates, func(state InvasionState, _ int) bool {
			return state.GetNumberOfCommandsUsed() == leastNumberOfCommandsUsed
		})
		return filteredStatesByLeastCommands
	}

	filters := []func(remainingStates []InvasionState) []InvasionState{filterByMostAliensKilled, filteredByHighestHpLeftOver, filteredByLeastNumberOfCommandsUsed}

	currentStates := states
	for _, filter := range filters {
		if len(currentStates) == 0 {
			break
		}
		currentStates = filter(currentStates)
	}
	return currentStates
}

func RunAllPossibleInvasionStatesToCompletion(initialState InvasionState) []InvasionState {
	endingStates := []InvasionState{}

	backtrack := func(currentState InvasionState) {}

	backtrack = func(currentState InvasionState) {
		volleyState := currentState.AttackAllAliens().AliensAttack()
		focusedState := currentState.AttackHighestDamageAlien().AliensAttack()
		focusedVolleyState := currentState.AttackHighestDamagingHalf().AliensAttack()
		if !volleyState.IsOver() {
			backtrack(volleyState)
		} else {
			endingStates = append(endingStates, volleyState)
		}
		if !focusedState.IsOver() {
			backtrack(focusedState)
		} else {
			endingStates = append(endingStates, focusedState)
		}
		if !focusedVolleyState.IsOver() {
			backtrack(focusedVolleyState)
		} else {
			endingStates = append(endingStates, focusedVolleyState)
		}
	}

	backtrack(initialState)
	return endingStates
}

// The invasion is over if and only if all aliens are dead or the remaining hp is empty.
func (i InvasionState) IsOver() bool {
	return len(i.aliensLeft) == 0 || i.hpLeft <= 0
}

// Returns the state of the invasion when the aliens attack
func (i InvasionState) AliensAttack() InvasionState {
	if len(i.aliensLeft) == 0 {
		return i
	}
	combinedAttack := lo.Reduce(i.aliensLeft, func(dmg int, alien Alien, idx int) int {
		return dmg + alien.Atk
	}, 0)
	return InvasionState{
		aliensLeft: i.aliensLeft,
		hpLeft:     i.hpLeft - combinedAttack,
		commands:   append(slices.Clone(i.commands), "alienAttack"),
	}
}

// Returns the state of the invasion when attacking all aliens.
func (i InvasionState) AttackAllAliens() InvasionState {
	// Attack all aliens to deal one damage
	newAliens := lo.Map(i.aliensLeft, func(alien Alien, _ int) Alien {
		return alien.TakeDamage(1)
	})
	// Filter out all the aliens that are dead.
	filteredAliens := lo.Filter(newAliens, func(alien Alien, _ int) bool {
		return alien.Hp > 0
	})
	return InvasionState{
		aliensLeft: filteredAliens,
		hpLeft:     i.hpLeft,
		commands:   append(slices.Clone(i.commands), "volley"),
	}
}

// Returns the state of the invasion when killing the highest damage alien.
func (i InvasionState) AttackHighestDamageAlien() InvasionState {
	// Invariant that we need to maintain, the highest damaging alien will the at the front of the list.
	return InvasionState{
		aliensLeft: i.aliensLeft[1:],
		hpLeft:     i.hpLeft,
		commands:   append(slices.Clone(i.commands), "focusedShot"),
	}
}

func (i InvasionState) AttackHighestDamagingHalf() InvasionState {
	// Pick the highest 1/2 of the aliens
	mid := len(i.aliensLeft) / 2
	// The first half will be all the aliens sorted by attack power.
	firstHalf := i.aliensLeft[:mid]
	// The second half will be spared.
	secondHalf := i.aliensLeft[mid:]
	firstHalfDamaged := lo.Map(firstHalf, func(alien Alien, _ int) Alien {
		return alien.TakeDamage(2)
	})
	filteredFirstHalfDamaged := lo.Filter(firstHalfDamaged, func(alien Alien, _ int) bool {
		return alien.Hp > 0
	})
	aliensLeft := append(filteredFirstHalfDamaged, secondHalf...)
	return InvasionState{
		aliensLeft: aliensLeft,
		hpLeft:     i.hpLeft,
		commands:   append(slices.Clone(i.commands), "focusedVolley"),
	}
}

// A Command consists one of:
// - Attack all aliens dealing 1 HP
// - Focus the highest ATK Alien, killing them instantly.
// - Focus the highest 1/2 (floor) Atk Aliens, dealing 2 hp.

const (
	UPPER_ALIEN_AMOUNT = 20
	LOWER_ALIEN_AMOUNT = 10
	ALIEN_ATK_HP_UPPER = 4
	ALIEN_ATK_HP_LOWER = 1
)

// Creates a random alien invasion, with aliens ranging from 10 to 20 aliens.
func GenerateAlienInvasion(rng *rand.Rand) []Alien {
	numAliens := rng.Intn(UPPER_ALIEN_AMOUNT-LOWER_ALIEN_AMOUNT) + LOWER_ALIEN_AMOUNT
	aliens := []Alien{}

	for range numAliens {
		alienHPVal := rng.Intn(ALIEN_ATK_HP_UPPER-ALIEN_ATK_HP_LOWER) + ALIEN_ATK_HP_LOWER
		alienAtkVal := rng.Intn(ALIEN_ATK_HP_UPPER-ALIEN_ATK_HP_LOWER) + ALIEN_ATK_HP_LOWER
		alien := CreateAlien(alienHPVal, alienAtkVal)
		aliens = append(aliens, alien)
	}
	return aliens
}

type Alien struct {
	Hp  int
	Atk int
}

// Creates an Alien with HP and ATK ranging from the upper and lower bounds
func CreateAlien(hp int, atk int) Alien {
	// Temporarily make the uuid generation dependent on the rng.
	return Alien{Hp: hp, Atk: atk}
}

func (a Alien) TakeDamage(dmg int) Alien {
	return Alien{
		Hp:  a.Hp - dmg,
		Atk: a.Atk,
	}
}
