package utils

import (
	"hash/fnv"
	"math/rand"

	"github.com/google/uuid"
)

// Creates a seeded random number generator from the given ID hash.
func CreateRNGFromHash(id uuid.UUID) *rand.Rand {
	h := fnv.New64a()
	h.Write(id[:])
	hash := int64(h.Sum64())
	return rand.New(rand.NewSource(hash))
}

func GenerateRandomNumWithinRange(rng *rand.Rand, lowerBound int, upperBound int) int {
	hp := rng.Intn(upperBound-lowerBound) + lowerBound
	return hp
}
