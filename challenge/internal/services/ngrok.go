package services

import (
	"math/rand"
	"net/http"
)

type NgrokChallenge struct {
	Requests []NgrokRequest
}

type NgrokChallengeScore struct {
	Score int
}

type NgrokRequest interface {
	Execute(client *http.Client, baseURL string) error
}

type NgrokPostRequest struct {
	Name   string
	Points int
	Path   string
	Body   []NgrokAlien
}

type NgrokGetRequest struct {
	Name           string
	Points         int
	Path           string
	ExpectedCount  int
	ExpectedFilter func(NgrokAlien) bool
}

type NgrokAlien struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Stats Stats  `json:"stats"`
}

type Stats struct {
	Atk int `json:"atk"`
	HP  int `json:"hp"`
}

var alienNames = []string{
	"Zorg", "Blorg", "Klaatu", "Xenu", "Nebula", "Vortex", "Plasma", "Quantum",
	"Shadow", "Void", "Crystal", "Ancient", "Cosmic", "Stellar", "Galactic",
	"Destroyer", "Warrior", "Hunter", "Guardian", "Commander", "Overlord",
	"Scout", "Fighter", "Assassin", "Spitter", "Grunt", "Horror",
}

const (
	REGULAR_ALIEN = "regular"
	ELITE_ALIEN   = "elite"
	BOSS_ALIEN    = "boss"
)

func (t NgrokPostRequest) Execute(client *http.Client, baseURL string) error {
	return nil
}
func (t NgrokGetRequest) Execute(client *http.Client, baseURL string) error {
	return nil
}

func generateRandomFilterTests(rng *rand.Rand, aliens []NgrokAlien) []NgrokRequest {
	panic("Not implemented.")
}

func generateRandomAliens(rng *rand.Rand, count int) []NgrokAlien {
	aliens := make([]NgrokAlien, count)

	for i := 0; i < count; i++ {
		nameIdx := rng.Intn(len(alienNames))
		name := alienNames[nameIdx]

		var alienType string
		roll := rng.Float32()
		if roll < 0.6 {
			alienType = REGULAR_ALIEN
		} else if roll < 0.85 {
			alienType = ELITE_ALIEN
		} else {
			alienType = BOSS_ALIEN
		}

		alienHPVal := rng.Intn(ALIEN_ATK_HP_UPPER-ALIEN_ATK_HP_LOWER) + ALIEN_ATK_HP_LOWER
		alienAtkVal := rng.Intn(ALIEN_ATK_HP_UPPER-ALIEN_ATK_HP_LOWER) + ALIEN_ATK_HP_LOWER

		aliens[i] = NgrokAlien{
			Name: name,
			Type: alienType,
			Stats: Stats{
				Atk: alienAtkVal,
				HP:  alienHPVal,
			},
		}
	}

	return aliens
}
