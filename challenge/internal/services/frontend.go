package services

import (
	"generate_technical_challenge_2025/internal/api"
	"generate_technical_challenge_2025/internal/data"
	"generate_technical_challenge_2025/internal/utils"
	"math/rand"

	"fmt"

	"github.com/google/uuid"
)

type AlienType string

const (
	AlienTypeRegular AlienType = "Regular"
	AlienTypeElite   AlienType = "Elite"
	AlienTypeBoss    AlienType = "Boss"
)

type DetailedAlien struct {
	ID         string    `json:"id"`
	BaseAlien  Alien     `json:"base_alien"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Type       AlienType `json:"type"`
	Spd        int       `json:"spd"`
	ProfileURL string    `json:"profile_url"`
}

var alienProfileURLs = map[AlienType]string{
	AlienTypeRegular: "https://robohash.org/regular-alien?set=set2&size=200x200",
	AlienTypeElite:   "https://robohash.org/elite-alien?set=set3&size=200x200",
	AlienTypeBoss:    "https://robohash.org/boss-alien?set=set4&size=200x200",
}

// Creates a DetailedAlien with ID, HP, ATK, and SPD stats, with a first/last name, profile URL, and AlienType.
func CreateDetailedAlien(id string, hp int, atk int, spd int,
	firstName string, lastName string, alienType AlienType, profileURL string) DetailedAlien {
	return DetailedAlien{
		ID:         id,
		BaseAlien:  Alien{Hp: hp, Atk: atk},
		FirstName:  firstName,
		LastName:   lastName,
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

	firstNameIndex := rng.Intn(len(data.AlienFirstNames))
	firstName := data.AlienFirstNames[firstNameIndex]

	lastNameIndex := rng.Intn(len(data.AlienLastNames))
	lastName := data.AlienLastNames[lastNameIndex]

	typeIndex := rng.Intn(len(alienTypes))
	alienType := alienTypes[typeIndex]

	profileURL := alienProfileURLs[alienType]

	alienID := generateAlienID(rng, alienIdx)

	alien := CreateDetailedAlien(alienID, hp, atk, spd, firstName, lastName, alienType, profileURL)
	return alien
}

func generateAlienID(rng *rand.Rand, alienIdx int) string {
	alienSeed := rng.Int63() + int64(alienIdx*1000)
	alienRNG := rand.New(rand.NewSource(alienSeed))
	alienID := ""
	alienIDLength := 6
	for range alienIDLength {
		alienID += fmt.Sprint(alienRNG.Intn(10)) // [0, n)
	}

	return alienID
}
