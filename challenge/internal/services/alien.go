package services

import (
	"generate_technical_challenge_2025/internal/utils"
	"hash/fnv"
	"maps"
	"math/rand"
	"slices"

	"github.com/etnz/permute"
	"github.com/samber/lo"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
)

// Represents the state of the invasion
type InvasionState struct {
	Budget         int
	WallDurability int
	Waves          [][]Alien
	GunsPurchased  []Weapon
	GunQueue       map[uint]map[Weapon][]Alien
}

// GenerateUniqueChallenge implements ChallengeService.
func GenerateInvasionState(id uuid.UUID) InvasionState {
	// Create a unique hash for the id so we can generate deterministically per id.
	h := fnv.New64a()
	h.Write(id[:])
	hash := int64(h.Sum64())
	rng := rand.New(rand.NewSource(hash))
	uuid.SetRand(rng)

	waves := generateWaves(rng)
	// Build the Invasion State
	state := InvasionState{}
	state.Waves = waves
	state.Budget = 100
	state.WallDurability = 100
	state.GunsPurchased = []Weapon{}
	state.GunQueue = map[uint]map[Weapon][]Alien{}
	uuid.SetRand(nil)
	return state
}

func generateWaves(rng *rand.Rand) [][]Alien {
	var waves [][]Alien
	// Wave 1, generate 5 Aliens, 3 to 5 Regular Aliens and 0 to 1 Swift Aliens
	waveOneBounds := []AlienGenerator{{lower: 3, upper: 5, supplier: CreateAlien}}
	// Wave 2, generate 10 Aliens, 3 to 5 Regular Aliens and 3 to 5 Swift Aliens and 0 to 1 Boss Aliens
	waveTwoBounds := []AlienGenerator{{lower: 3, upper: 5, supplier: CreateAlien}}
	// Wave 3, generate 20 Aliens, 5 to 10 Regular Aliens, 5 to 9 Swift Aliens and 1 to 3 Boss Aliens
	waveThreeBounds := []AlienGenerator{{lower: 3, upper: 5, supplier: CreateAlien}}

	waveOneAliens := generateAliens(rng, waveOneBounds)
	waveTwoAliens := generateAliens(rng, waveTwoBounds)
	waveThreeAliens := generateAliens(rng, waveThreeBounds)

	waves = append(waves, waveOneAliens)
	waves = append(waves, waveTwoAliens)
	waves = append(waves, waveThreeAliens)
	return waves
}

func generateAliens(rng *rand.Rand, ranges []AlienGenerator) []Alien {
	var aliens []Alien
	for _, alienGenerator := range ranges {
		numToGenerate := rng.Intn(int(alienGenerator.upper)+1-int(alienGenerator.lower)) + int(alienGenerator.lower)
		i := 0
		for i < numToGenerate {
			aliens = append(aliens, alienGenerator.supplier())
			i = i + 1
		}
	}
	return aliens
}

// Generate all possible weapon purchases given the budget.
// NOTE: This can include equivalent purchases such as Weapon 1 and Weapon 2 aswell as Weapon 2 and Weapon 1
func GenerateAllPossibleWeaponPurchasesFromBudget(budget int) [][]Weapon {
	// Begin with a decision tree, You can either purchase 3 Guns
	// Recursively call your backtracking algorithm. At each stage you can purchase one of three items, unless you run out of money then it stops.
	weapons := [][]Weapon{}
	var backtrack func(candidates []Weapon, remainingBudget int)
	weaponsSupplier := []func() Weapon{CreateTurretWeapon, CreateMachineGunWeapon, CreateRayGunWeapon}
	backtrack = func(candidates []Weapon, remainingBudget int) {
		for _, supplier := range weaponsSupplier {
			w := supplier()
			// Attempt to buy that weapon, only add to the final list and recur if we can afford it.
			if int(w.Cost) <= remainingBudget {
				newCandidates := append(slices.Clone(candidates), w)
				if !lo.ContainsBy(weapons, func(wpns []Weapon) bool {
					weaponSet1 := mapset.NewSet(wpns...)
					weaponSet2 := mapset.NewSet(newCandidates...)
					return weaponSet1.Equal(weaponSet2)
				}) {
					weapons = append(weapons, newCandidates)
				}
				backtrack(newCandidates, remainingBudget-int(w.Cost))
			}
		}
	}
	backtrack([]Weapon{}, budget)
	return weapons
}

func FindWeaponQueueAssignments(weapons []Weapon, aliens []Alien) []map[Weapon][]Alien {
	assignments := []map[Weapon][]Alien{}
	var backtrack func(aliensLeft []Alien, candidates map[Weapon][]Alien)
	backtrack = func(aliensLeft []Alien, candidates map[Weapon][]Alien) {
		if len(aliensLeft) == 0 {
			assignments = append(assignments, candidates)
			return
		}
		alien := aliensLeft[0]
		for wpn, queue := range candidates {
			if len(queue) >= 3 {
				continue
			}
			copied := lo.Assign(map[Weapon][]Alien{}, candidates)
			copied[wpn] = append(queue, alien)
			backtrack(aliensLeft[1:], copied)
		}
	}
	backtrack(aliens, map[Weapon][]Alien{})
	return assignments
}

// From the given weapons and aliens, compute all possible arrangement of valid weapon queues.
// Invariant: The assigned queues must have only distinct aliens
func GenerateAllPossibleValidWeaponQueuesCrazy(weapons []Weapon, aliens []Alien) []map[Weapon][]Alien {
	queues := []map[Weapon][]Alien{}
	// Algorithm:
	// Generate all possible subsets from the given aliens.
	// From there, get every single possible list of sets of size exactly of length weapons
	// This is a combinatorial problem, essential boils down to the number of ways we can choose 4 queues from each arrangement of aliens.
	// From there filter out all of the list of sets whose union is not equal to the original alien set
	// and whose intersection is not all empty.
	// This should leave us with all the sets that satisify the invariants.

	originalAlienSet := mapset.NewSet(aliens...)
	alienPowerSet := utils.PowerSet[Alien](aliens)
	// Next convert each alien sublist to a set so order doesn't matter
	convertedAlienPowerSet := lo.Map(alienPowerSet, func(items []Alien, _ int) mapset.Set[Alien] {
		return mapset.NewSet(items...)
	})
	allPossibleSubSets := permute.Combinations(len(weapons), convertedAlienPowerSet)
	allPossibleValidSubSets := [][]mapset.Set[Alien]{}
	for combination := range allPossibleSubSets {
		// Condition One, The union of all aliens must be equal to the original Alien Set
		unionedAlienSet := lo.Reduce(combination, func(unionAcc mapset.Set[Alien], curr mapset.Set[Alien], _ int) mapset.Set[Alien] {
			return unionAcc.Union(curr)
		}, mapset.NewSet[Alien]())
		if !unionedAlienSet.Equal(originalAlienSet) {
			continue
		}
		// Condition Two, all sets must be mutually exclusive to one another
		mutuallyExclusiveCheck := lo.Reduce(combination, func(intersectAcc mapset.Set[Alien], curr mapset.Set[Alien], _ int) mapset.Set[Alien] {
			return intersectAcc.Intersect(curr)
		}, mapset.NewSet[Alien]())
		if mutuallyExclusiveCheck.Cardinality() != 0 {
			continue
		}
		// Finally once these conditions are met, we now have the valid subset
		allPossibleValidSubSets = append(allPossibleValidSubSets, combination)
	}

	for _, validSubSets := range allPossibleValidSubSets {
		// We have guranteed that each of our validSubSet lists are of size weapon len
		// However we care about the order so lets convert these sets back into lists and take the
		// permutation of each valid list
		validLists := lo.Map(validSubSets, func(alienSet mapset.Set[Alien], _ int) []Alien {
			return alienSet.ToSlice()
		})
		weaponQueue := map[Weapon][]Alien{}
		for _, weapon := range weapons {
			// Create the hasmap of all weapons
			for _, validList := range validLists {
				// There are validLists however we need to permute them since order matters
				allPossibleLists := permute.Permutations(validList)
				for _, alienQueue := range allPossibleLists {
					weaponQueue[weapon] = alienQueue
				}
			}
		}
		queues = append(queues, weaponQueue)
	}
	return queues
}

// Determines if the invasion is over on these conditions:
// The WallDurability is <= 0
// All aliens of all waves have been exhausted.
func (s InvasionState) IsOver() bool {
	return s.WallDurability <= 0 || len(s.Waves) == 0
}

// Purchases a weapon and returns a new Invasion state with either that weapon
// purchased and the remaining budget afterwards, or if the weapon cannot be
// afforded the same Invasion state.
func (s InvasionState) PurchaseWeapon(weapon Weapon) InvasionState {
	remainingBudget := s.Budget - int(weapon.Cost)
	if remainingBudget < 0 {
		panic("Attempted to spend more than the budget.")
	}
	return InvasionState{
		Budget:         remainingBudget,
		WallDurability: s.WallDurability,
		Waves:          s.Waves,
		GunsPurchased:  append(slices.Clone(s.GunsPurchased), weapon),
		GunQueue:       s.GunQueue,
	}
}

// Invariant: The Gun given MUST Be a gun in the GunsPurchased or will panic.
// Invariant: The queue assigned must be a unique selection of aliens.
// Invariant: The queue assigned must only contain aliens from the current wave.
// Invariant: The queue assigned must have DISTINCT ALIENS FROM ALL OTHER SELECTIONS.
func (s InvasionState) ConfigureWeapon(wave uint, weapon Weapon, queue []Alien) InvasionState {
	if !slices.Contains(s.GunsPurchased, weapon) {
		panic("Attempted to use a gun that we did not purchase.")
	}
	if wave > uint(len(s.Waves)) {
		panic("Attempted to access a wave not defined.")
	}
	currentWaveAliens := s.Waves[wave]
	alienSet := mapset.NewSet(currentWaveAliens...)
	queueSet := mapset.NewSet(queue...)
	if !queueSet.IsSubset(alienSet) {
		panic("Attempted to give a queue of aliens that was not a subset of the current alien wave.")
	}
	otherConfigurations := s.GunQueue[wave]
	for otherWeapon, weaponQueue := range otherConfigurations {
		if weapon == otherWeapon {
			panic("Cannot reassign a weapon in the queue.")
		}
		otherQueueSet := mapset.NewSet(weaponQueue...)
		if otherQueueSet.Intersect(queueSet).Cardinality() != 0 {
			panic("Cannot have shared target aliens across queues.")
		}
	}
	gunQueueClone := maps.Clone(s.GunQueue)
	gunQueueClone[wave] = map[Weapon][]Alien{weapon: queue}
	return InvasionState{
		Budget:         s.Budget,
		WallDurability: s.WallDurability,
		Waves:          slices.Clone(s.Waves),
		GunsPurchased:  slices.Clone(s.GunsPurchased),
		GunQueue:       gunQueueClone,
	}
}

// Begin Alien Data Definitions

type Alien struct {
	Hp  int
	Atk uint
	ID  uuid.UUID
}

func CreateAlien() Alien {
	return Alien{
		Hp:  5,
		Atk: 5,
		ID:  uuid.New(),
	}
}

func (a Alien) TakeDamage(damage int) Alien {
	return Alien{
		Hp:  a.Hp - damage,
		Atk: a.Atk,
		ID:  a.ID,
	}
}

// Instructs the generator to generate a range of a type of alien.
type AlienGenerator struct {
	lower    uint // Inclusive
	upper    uint // Inclusive
	supplier func() Alien
}

// Begin Weapon Data Definitions

// A Weapon is one of:
// Turret
// MachineGun
// RayGun

type WeaponType int

const (
	Turret WeaponType = iota
	MachineGun
	RayGun
)

// Internally we assign each Weapon a unique id, so that we can differentiate guns.
type Weapon struct {
	Atk  uint
	Cost uint
	Type WeaponType
	ID   uuid.UUID
}

// Given the alien queue, determines the number of turns it would take to kill all aliens.
func (w Weapon) MaxTurnsToKill(aliens []Alien) int {
	turns := 0
	currentAliens := aliens

	for len(currentAliens) != 0 {
		// Get the front of the alien
		alien := currentAliens[0].TakeDamage(int(w.Atk))
		if alien.Hp <= 0 {
			currentAliens = currentAliens[1:]
		} else {
			currentAliens[0] = alien
		}
		turns += 1
	}
	return turns
}

func CreateTurretWeapon() Weapon {
	return Weapon{
		Atk:  1,
		Cost: 10,
		Type: Turret,
		ID:   uuid.New(),
	}
}

func CreateMachineGunWeapon() Weapon {
	return Weapon{
		Atk:  3,
		Cost: 30,
		Type: MachineGun,
		ID:   uuid.New(),
	}
}

func CreateRayGunWeapon() Weapon {
	return Weapon{
		Atk:  5,
		Cost: 50,
		Type: RayGun,
		ID:   uuid.New(),
	}
}
