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
	Body   []DetailedAlien
}

type NgrokGetRequest struct {
	Name           string
	Points         int
	Path           string
	ExpectedCount  int
	ExpectedFilter func(DetailedAlien) bool
}

type Stats struct {
	Atk int `json:"atk"`
	HP  int `json:"hp"`
}

func (t NgrokPostRequest) Execute(client *http.Client, baseURL string) error {
	return nil
}
func (t NgrokGetRequest) Execute(client *http.Client, baseURL string) error {
	return nil
}

func generateRandomFilterTests(rng *rand.Rand, aliens []DetailedAlien) []NgrokRequest {
	panic("Not implemented.")
}
