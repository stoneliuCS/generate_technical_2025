package services

import (
	"generate_technical_challenge_2025/internal/transactions"
	"hash/fnv"
	"log/slog"
	"math/rand"

	"github.com/google/uuid"
)

type ChallengeService interface {
	GenerateUniqueAlienChallenge(id uuid.UUID) InvasionState
	SolveAlienChallenge(state InvasionState) InvasionState
}

type ChallengeServiceImpl struct {
	logger       *slog.Logger
	transactions transactions.ChallengeTransactions
}

// SolveChallenge implements ChallengeService.
func (c ChallengeServiceImpl) SolveAlienChallenge(state InvasionState) InvasionState {
	return InvasionState{}
}

// GenerateUniqueChallenge implements ChallengeService.
func (c ChallengeServiceImpl) GenerateUniqueAlienChallenge(id uuid.UUID) InvasionState {
	// Create a unique hash for the id so we can generate deterministically per id.
	h := fnv.New64a()
	h.Write(id[:])
	hash := int64(h.Sum64())
	rng := rand.New(rand.NewSource(hash))

	waves := c.generateWaves(rng)
	// Build the Invasion State
	state := InvasionState{}
	state.Waves = waves
	state.Budget = 100
	state.WallDurability = 100
	state.GunsPurchased = []Weapon{}
	state.GunQueue = map[uint]map[Weapon][]Alien{}
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
