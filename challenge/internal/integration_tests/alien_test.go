package integrationtests

import (
	"generate_technical_challenge_2025/internal/services"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

const (
	NORTHEASTERN_TEST_EMAIL = "johndoe@northeastern.edu"
	NORTHEASTERN_TEST_NUID  = "123456789"
)

// A Collection of brute force tests to ensure reliablilty
func TestBackendAlienChallengeFullIntegration(t *testing.T) {
	// User registers with their NUID and Northeastern Email
	client := CLIENT.AddBody(map[string]any{
		"email": NORTHEASTERN_TEST_EMAIL,
		"nuid":  NORTHEASTERN_TEST_NUID,
	}).AddHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	testVerify := client.POST("/api/v1/member/register")
	testVerify.AssertStatusCode(201, t)
	var res map[string]string
	testVerify.GetBody(&res, t)
	testVerify = client.GET("/api/v1/challenge/backend/" + res["id"] + "/aliens")
	testVerify.AssertStatusCode(200, t)
	// Deserialize alien invasion data
	type AlienData struct {
		HP  int `json:"hp"`
		ATK int `json:"atk"`
	}

	type AlienWaveData struct {
		ChallengeID string      `json:"challengeID"`
		Aliens      []AlienData `json:"aliens"`
		HP          int         `json:"hp"`
	}
	target := []AlienWaveData{}
	testVerify.GetBody(&target, t)
	assert.Len(t, target, services.NUM_WAVES)
	// Build invasion states
	deserializedInvasionStates := lo.SliceToMap(target, func(item AlienWaveData) (uuid.UUID, services.InvasionState) {
		aliens := lo.Map(item.Aliens, func(alienData AlienData, _ int) services.Alien {
			return services.CreateAlien(alienData.HP, alienData.ATK)
		})
		return uuid.MustParse(item.ChallengeID), services.CreateInvasionState(aliens, item.HP)
	})
	// CASE 1: PERFECT USER SUBMISSION
	perfectDeserializedInvasionStates := lo.MapEntries(deserializedInvasionStates, func(challengeIDs uuid.UUID, initialState services.InvasionState) (uuid.UUID, services.InvasionState) {
		perfectSolution := services.OracleSolution(initialState)
		return challengeIDs, perfectSolution
	})
	// serialize into response expected by server
	serializedAnswers := lo.MapToSlice(perfectDeserializedInvasionStates, func(key uuid.UUID, val services.InvasionState) map[string]any {
		return map[string]any{
			"challengeID": key.String(),
			"state": map[string]any{
				"remainingHP":     val.GetHpLeft(),
				"remainingAliens": val.GetAliensLeft(),
				"commands":        val.GetCommandsUsed(),
			},
		}
	})
	testVerify = CLIENT.AddBody(serializedAnswers).AddHeaders(map[string]string{
		"Content-Type": "application/json",
	}).POST("/api/v1/challenge/backend/" + res["id"] + "/aliens/submit")
	testVerify.AssertStatusCode(200, t)
	response := map[string]any{}
	testVerify.GetBody(&response, t)
	assert.True(t, response["valid"].(bool))
	assert.Equal(t, 0.0, response["score"].(float64))
	// CASE 2: USER DOES NOT GIVE ANYTHING
	serializedAnswers = []map[string]any{}
	testVerify = CLIENT.AddBody(serializedAnswers).AddHeaders(map[string]string{
		"Content-Type": "application/json",
	}).POST("/api/v1/challenge/backend/" + res["id"] + "/aliens/submit")
	testVerify.AssertStatusCode(200, t)
	testVerify.GetBody(&response, t)
	assert.False(t, response["valid"].(bool))
	assert.Equal(t, "Challenge IDs do not match.", response["message"].(string))
}
