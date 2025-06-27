package services

import (
	"maps"
	"slices"

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
		return s
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

// An AlienType is One Of:
// Regular
// Swift
// Boss

type AlienType int

const (
	Regular AlienType = iota
	Swift
	Boss
)

type Alien struct {
	Hp   uint
	Atk  uint
	Type AlienType
	ID   uuid.UUID
}

func CreateRegularAlien() Alien {
	return Alien{
		Hp:   5,
		Atk:  1,
		Type: Regular,
		ID:   uuid.New(),
	}
}

func CreateSwiftAlien() Alien {
	return Alien{
		Hp:   3,
		Atk:  5,
		Type: Swift,
		ID:   uuid.New(),
	}
}

func CreateBossAlien() Alien {
	return Alien{
		Hp:   10,
		Atk:  10,
		Type: Boss,
		ID:   uuid.New(),
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
