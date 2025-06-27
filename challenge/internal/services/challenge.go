package services

import (
	"generate_technical_challenge_2025/internal/transactions"
	"hash/fnv"
	"log/slog"
	"math/rand"

	"github.com/google/uuid"
)

// Represents the state of the invasion
type InvasionState struct {
	Budget         uint
	WallDurability uint
	DefaultWeapons []Weapon
	DefaultAliens  []Alien
	Waves          [][]Alien
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
}

func CreateRegularAlien() Alien {
	return Alien{
		Hp:   5,
		Atk:  1,
		Type: Regular,
	}
}

func CreateSwiftAlien() Alien {
	return Alien{
		Hp:   3,
		Atk:  5,
		Type: Swift,
	}
}

func CreateBossAlien() Alien {
	return Alien{
		Hp:   10,
		Atk:  10,
		Type: Boss,
	}
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

type Weapon struct {
	Atk  uint
	Cost uint
	Type WeaponType
}

func CreateTurretWeapon() Weapon {
	return Weapon{
		Atk:  1,
		Cost: 10,
		Type: Turret,
	}
}

func CreateMachineGunWeapon() Weapon {
	return Weapon{
		Atk:  3,
		Cost: 30,
		Type: MachineGun,
	}
}

func CreateRayGunWeapon() Weapon {
	return Weapon{
		Atk:  5,
		Cost: 50,
		Type: RayGun,
	}
}

type ChallengeService interface {
	GenerateUniqueChallenge(id uuid.UUID) InvasionState
}

type ChallengeServiceImpl struct {
	logger       *slog.Logger
	transactions transactions.ChallengeTransactions
}

// Instructs the generator to generate a range of a type of alien.
type AlienGenerator struct {
	lower    uint // Inclusive
	upper    uint // Inclusive
	supplier func() Alien
}

// GenerateUniqueChallenge implements ChallengeService.
func (c ChallengeServiceImpl) GenerateUniqueChallenge(id uuid.UUID) InvasionState {
	// Create a unique hash for the id so we can generate deterministically per id.
	h := fnv.New64a()
	h.Write(id[:])
	hash := int64(h.Sum64())
	rng := rand.New(rand.NewSource(hash))

	// Generate Alien Wave Data
	defaultAliens := []Alien{CreateRegularAlien(), CreateSwiftAlien(), CreateBossAlien()}
	defaultWeapons := []Weapon{CreateTurretWeapon(), CreateMachineGunWeapon(), CreateRayGunWeapon()}
	waves := c.generateWaves(rng)

	// Build the Invasion State
	state := InvasionState{}
	state.DefaultAliens = defaultAliens
	state.DefaultWeapons = defaultWeapons
	state.Waves = waves
	state.Budget = 100
	state.WallDurability = 100
	return state
}

func (c ChallengeServiceImpl) generateWaves(rng *rand.Rand) [][]Alien {
	var waves [][]Alien
	const WAVE_1_NUM_ALIENS = 5
	const WAVE_2_NUM_ALIENS = WAVE_1_NUM_ALIENS * 2
	const WAVE_3_NUM_ALIENS = WAVE_2_NUM_ALIENS * 2

	// Wave 1, generate 5 Aliens, 3 to 5 Regular Aliens and 0 to 1 Swift Aliens
	waveOneBounds := []AlienGenerator{{lower: 3, upper: 5, supplier: CreateRegularAlien}, {lower: 0, upper: 1, supplier: CreateSwiftAlien}}
	// Wave 2, generate 10 Aliens, 3 to 5 Regular Aliens and 3 to 5 Swift Aliens and 0 to 1 Boss Aliens
	waveTwoBounds := []AlienGenerator{{lower: 3, upper: 5, supplier: CreateRegularAlien}, {lower: 3, upper: 5, supplier: CreateSwiftAlien}, {lower: 0, upper: 1, supplier: CreateBossAlien}}
	// Wave 3, generate 20 Aliens, 5 to 10 Regular Aliens, 5 to 9 Swift Aliens and 1 to 3 Boss Aliens
	waveThreeBounds := []AlienGenerator{{lower: 5, upper: 10, supplier: CreateRegularAlien}, {lower: 5, upper: 9, supplier: CreateSwiftAlien}, {lower: 1, upper: 3, supplier: CreateBossAlien}}

	waveOneAliens := c.generateAliens(rng, WAVE_1_NUM_ALIENS, waveOneBounds)
	waveTwoAliens := c.generateAliens(rng, WAVE_2_NUM_ALIENS, waveTwoBounds)
	waveThreeAliens := c.generateAliens(rng, WAVE_3_NUM_ALIENS, waveThreeBounds)

	waves = append(waves, waveOneAliens)
	waves = append(waves, waveTwoAliens)
	waves = append(waves, waveThreeAliens)
	return waves
}

func (c ChallengeServiceImpl) generateAliens(rng *rand.Rand, numberOfAliens uint, ranges []AlienGenerator) []Alien {
	var aliens []Alien
	for uint(len(aliens)) < numberOfAliens {
		for _, alienGenerator := range ranges {
			numToGenerate := rng.Intn(int(alienGenerator.upper)-int(alienGenerator.lower)+1) + int(alienGenerator.lower)
			i := 0
			for i < numToGenerate && len(aliens) < int(numberOfAliens) {
				aliens = append(aliens, alienGenerator.supplier())
				i = i + 1
			}
		}
	}
	return aliens
}

func CreateChallengeService(logger *slog.Logger, transactions transactions.ChallengeTransactions) ChallengeService {
	return ChallengeServiceImpl{
		logger: logger, transactions: transactions,
	}
}
