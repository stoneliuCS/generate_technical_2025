package services

import (
	"generate_technical_challenge_2025/internal/api"

	"github.com/google/uuid"
)

type AlienType string

const (
	AlienTypeRegular AlienType = "Regular"
	AlienTypeElite   AlienType = "Elite"
	AlienTypeBoss    AlienType = "Boss"
)

type DetailedAlien struct {
	ID        uuid.UUID
	BaseAlien Alien
	Name      string
	Type      AlienType
	Spd       int
}

// Creates a DetailedAlien with HP, ATK, and SPD stats, with a name and AlienType.
func CreateDetailedAlien(id uuid.UUID, hp int, atk int, spd int, name string, alienType AlienType) DetailedAlien {
	return DetailedAlien{ID: id, BaseAlien: Alien{Hp: hp, Atk: atk}, Name: name, Type: alienType, Spd: spd}
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
