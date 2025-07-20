package services

import (
	"generate_technical_challenge_2025/internal/api"
	"generate_technical_challenge_2025/internal/data"
	"generate_technical_challenge_2025/internal/utils"
	"math/rand"

	"github.com/google/uuid"
)

type AlienType string

const (
	AlienTypeRegular AlienType = "Regular"
	AlienTypeElite   AlienType = "Elite"
	AlienTypeBoss    AlienType = "Boss"
)

type DetailedAlien struct {
	ID         uuid.UUID
	BaseAlien  Alien
	Name       string
	Type       AlienType
	Spd        int
	ProfileURL string
}

// Creates a DetailedAlien with ID, HP, ATK, and SPD stats, with a name, profile URL, and AlienType.
func CreateDetailedAlien(id uuid.UUID, hp int, atk int, spd int,
	name string, alienType AlienType, profileURL string) DetailedAlien {
	return DetailedAlien{
		ID:         id,
		BaseAlien:  Alien{Hp: hp, Atk: atk},
		Name:       name,
		Type:       alienType,
		Spd:        spd,
		ProfileURL: profileURL}
}

func (at AlienType) ToAPI() api.APIV1ChallengeFrontendIDAliensGetOKItemType {
	switch at {
	case AlienTypeRegular:
		return api.APIV1ChallengeFrontendIDAliensGetOKItemTypeRegular
	case AlienTypeElite:
		return api.APIV1ChallengeFrontendIDAliensGetOKItemTypeElite
	case AlienTypeBoss:
		return api.APIV1ChallengeFrontendIDAliensGetOKItemTypeBoss
	default:
		return api.APIV1ChallengeFrontendIDAliensGetOKItemTypeRegular
	}
}

func GenerateDetailedAlien(rng *rand.Rand, memberID uuid.UUID, alienIdx int) DetailedAlien {
	hp := utils.GenerateRandomNumWithinRange(rng, ALIEN_ATK_HP_SPD_LOWER, ALIEN_ATK_HP_SPD_UPPER)
	atk := utils.GenerateRandomNumWithinRange(rng, ALIEN_ATK_HP_SPD_LOWER, ALIEN_ATK_HP_SPD_UPPER)
	spd := utils.GenerateRandomNumWithinRange(rng, ALIEN_ATK_HP_SPD_LOWER, ALIEN_ATK_HP_SPD_UPPER)

	nameIndex := rng.Intn(len(data.AlienNames))
	name := data.AlienNames[nameIndex]

	profileIndex := rng.Intn(len(data.AlienProfileURLs))
	profileURL := data.AlienProfileURLs[profileIndex]

	typeIndex := rng.Intn(len(alienTypes))
	alienType := alienTypes[typeIndex]

	alienID := generateSimpleUUID(rng, alienIdx)

	alien := CreateDetailedAlien(alienID, hp, atk, spd, name, alienType, profileURL)
	return alien
}

func generateSimpleUUID(rng *rand.Rand, alienIdx int) uuid.UUID {
	alienSeed := rng.Int63() + int64(alienIdx*1000)
	alienRNG := rand.New(rand.NewSource(alienSeed))
	var alienID uuid.UUID
	for i := range alienID {
		alienID[i] = byte(alienRNG.Intn(256))
	}
	return alienID
}
