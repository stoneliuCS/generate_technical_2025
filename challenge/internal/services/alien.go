package services

import (
	"generate_technical_challenge_2025/internal/utils"
	"math/rand"
	"slices"

	"github.com/samber/lo"
)

const (
	VOLLEY         = "volley"
	FOCUSED_VOLLEY = "focusedVolley"
	FOCUSED_SHOT   = "focusedShot"
)

// Represents the state of the invasion
type InvasionState struct {
	aliensLeft []Alien
	hpLeft     int
	commands   []string
}

func (i InvasionState) sortAliens() InvasionState {
	sortedAliens := slices.SortedFunc(slices.Values(i.aliensLeft), func(a1 Alien, a2 Alien) int {
		power1 := a1.Atk + a1.Hp
		power2 := a2.Atk + a2.Hp
		return power2 - power1 // Highest total power first
	})
	return InvasionState{
		aliensLeft: sortedAliens,
		hpLeft:     i.hpLeft,
		commands:   i.commands,
	}
}

func (i InvasionState) GetCommandsUsed() []string {
	return i.commands
}

func CreateInvasionState(aliens []Alien, startingHp int) InvasionState {
	// Sort all the aliens by attack power.
	return InvasionState{
		aliensLeft: aliens,
		hpLeft:     startingHp,
		commands:   []string{},
	}.sortAliens()
}

func RunCommandsToCompletion(startingState InvasionState, commands []string) *InvasionState {
	state := startingState
	mapFunc := func(s InvasionState, command string) InvasionState {
		mappings := map[string]func() InvasionState{
			VOLLEY:         s.AttackAliensModulo,
			FOCUSED_VOLLEY: s.AttackHighestDamagingHalf,
			FOCUSED_SHOT:   s.AttackHighestDamageAlien,
		}
		return mappings[command]().sortAliens().AliensAttack()
	}
	for _, command := range commands {
		if state.IsOver() {
			return nil
		}
		state = mapFunc(state, command)
	}
	return &state
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

func RunAllPossibleInvasionStatesToCompletion(initialState InvasionState) []InvasionState {
	endingStates := []InvasionState{}

	backtrack := func(currentState InvasionState) {}

	backtrack = func(currentState InvasionState) {
		if currentState.IsOver() {
			endingStates = append(endingStates, currentState)
		} else {
			volleyState := currentState.AttackAliensModulo().sortAliens().AliensAttack()
			focusedState := currentState.AttackHighestDamageAlien().sortAliens().AliensAttack()
			focusedVolleyState := currentState.AttackHighestDamagingHalf().sortAliens().AliensAttack()
			backtrack(volleyState)
			backtrack(focusedState)
			backtrack(focusedVolleyState)
		}
	}

	backtrack(initialState)
	return endingStates
}

func OracleSolution(initalState InvasionState) InvasionState {
	finalStates := RunAllPossibleInvasionStatesToCompletionGreedy(initalState)
	// Get the states with the smallest aliens left
	smallestAliensLeft := lo.MinBy(finalStates, func(s1 InvasionState, s2 InvasionState) bool {
		return s1.GetAliensLeft() < s2.GetAliensLeft()
	})
	smallestAlienStatesLeft := lo.Filter(finalStates, func(state InvasionState, idx int) bool {
		return state.GetAliensLeft() == smallestAliensLeft.GetAliensLeft()
	})
	// Get the states with the largest HP left
	largestHPLeft := lo.MaxBy(smallestAlienStatesLeft, func(s1 InvasionState, s2 InvasionState) bool {
		return s1.GetHpLeft() > s2.GetHpLeft()
	})
	largestHPStatesLeft := lo.Filter(smallestAlienStatesLeft, func(state InvasionState, idx int) bool {
		return state.GetHpLeft() == largestHPLeft.GetHpLeft()
	})
	return lo.MinBy(largestHPStatesLeft, func(s1 InvasionState, s2 InvasionState) bool {
		return s1.GetNumberOfCommandsUsed() < s2.GetNumberOfCommandsUsed()
	})
}

func RunAllPossibleInvasionStatesToCompletionGreedy(initialState InvasionState) []InvasionState {
	var backtrack func(currentState InvasionState)
	endingStates := []InvasionState{}

	backtrack = func(currentState InvasionState) {
		if currentState.IsOver() {
			endingStates = append(endingStates, currentState)
		} else {
			// It is never ideal to recur on a node of which the modulo is strictly less than or equal to the ceiling of aliens.
			if currentState.hpLeft%currentState.GetAliensLeft() >= (currentState.GetAliensLeft()+1)/2 {
				backtrack(currentState.AttackAliensModulo().sortAliens().AliensAttack())
			}
			backtrack(currentState.AttackHighestDamagingHalf().sortAliens().AliensAttack())
			backtrack(currentState.AttackHighestDamageAlien().sortAliens().AliensAttack())
		}
	}
	backtrack(initialState)
	return endingStates
}

// The invasion is over if and only if all aliens are dead or the remaining hp is empty.
func (i InvasionState) IsOver() bool {
	return len(i.aliensLeft) == 0 || i.hpLeft <= 0
}

func (i InvasionState) GetTotalAlienHPLeft() int {
	return lo.Reduce(i.aliensLeft, func(acc int, alien Alien, _ int) int {
		return acc + alien.Hp
	}, 0)
}

func (i InvasionState) GetTotalAlienAtkPower() int {
	return lo.Reduce(i.aliensLeft, func(acc int, alien Alien, _ int) int {
		return acc + alien.Atk
	}, 0)
}

// Returns the state of the invasion when the aliens attack
func (i InvasionState) AliensAttack() InvasionState {
	combinedAttack := lo.Reduce(i.aliensLeft, func(dmg int, alien Alien, idx int) int {
		return dmg + alien.Atk
	}, 0)
	return InvasionState{
		aliensLeft: i.aliensLeft,
		hpLeft:     i.hpLeft - combinedAttack,
		commands:   i.commands,
	}
}

// Returns the state of the invasion when attacking all aliens.
func (i InvasionState) AttackAliensModulo() InvasionState {
	// Attack all aliens to deal one damage
	portionOfAliensModulo := i.hpLeft % i.GetAliensLeft()
	aliensToHit := i.aliensLeft[:portionOfAliensModulo]
	restOfAliens := i.aliensLeft[portionOfAliensModulo:]

	newAliens := lo.Map(aliensToHit, func(alien Alien, _ int) Alien {
		return alien.TakeDamage(1)
	})
	// Filter out all the aliens that are dead.
	filteredAliens := lo.Filter(newAliens, func(alien Alien, _ int) bool {
		return alien.Hp > 0
	})
	return InvasionState{
		aliensLeft: append(filteredAliens, restOfAliens...),
		hpLeft:     i.hpLeft,
		commands:   append(slices.Clone(i.commands), VOLLEY),
	}
}

// Returns the state of the invasion when killing the highest damage alien.
func (i InvasionState) AttackHighestDamageAlien() InvasionState {
	// Invariant that we need to maintain, the highest damaging alien will the at the front of the list.
	if len(i.aliensLeft) == 0 {
		return i
	}
	return InvasionState{
		aliensLeft: i.aliensLeft[1:],
		hpLeft:     i.hpLeft,
		commands:   append(slices.Clone(i.commands), FOCUSED_SHOT),
	}
}

func (i InvasionState) AttackHighestDamagingHalf() InvasionState {
	// Pick the highest (Ceiling) 1/2 of the aliens
	mid := (len(i.aliensLeft) + 1) / 2
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
		commands:   append(slices.Clone(i.commands), FOCUSED_VOLLEY),
	}
}

// A Command consists one of:
// - Attack all aliens dealing 1 HP
// - Focus the highest ATK Alien, killing them instantly.
// - Focus the highest 1/2 (floor) Atk Aliens, dealing 2 hp.

const (
	UPPER_ALIEN_AMOUNT     = 20
	LOWER_ALIEN_AMOUNT     = 10
	ALIEN_ATK_HP_SPD_UPPER = 4 // [1, 4)
	ALIEN_ATK_HP_SPD_LOWER = 1
)

// Creates a random alien invasion, with aliens ranging from 10 to 20 aliens.
func GenerateAlienInvasion(rng *rand.Rand) []Alien {
	numAliens := utils.GenerateRandomNumWithinRange(rng, LOWER_ALIEN_AMOUNT, UPPER_ALIEN_AMOUNT)
	aliens := []Alien{}
	for range numAliens {
		alienHPVal := utils.GenerateRandomNumWithinRange(rng, ALIEN_ATK_HP_SPD_LOWER, ALIEN_ATK_HP_SPD_UPPER)
		alienAtkVal := utils.GenerateRandomNumWithinRange(rng, ALIEN_ATK_HP_SPD_LOWER, ALIEN_ATK_HP_SPD_UPPER)
		alien := CreateAlien(alienHPVal, alienAtkVal)
		aliens = append(aliens, alien)
	}
	return aliens
}

type Alien struct {
	Hp  int `json:"hp"`
	Atk int `json:"atk"`
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
