package integrationtests

import (
	"generate_technical_challenge_2025/internal/services"
	"testing"
)

// func TestRateLimit(t *testing.T) {
// 	registerForChallenge(t)

// 	for range 10 {
// 		CLIENT.POST(fmt.Sprintf("/api/v1/challenge/backend/%s/aliens/submit", memberUUID.String()))
// 	}

// 	exceedingLimit := CLIENT.POST(fmt.Sprintf("/api/v1/challenge/backend/%s/aliens/submit", memberUUID.String()))
// 	exceedingLimit.AssertStatusCode(429, t)
// }

const (
	NORTHEASTERN_TEST_EMAIL = "johndoe@northeastern.edu"
	NORTHEASTERN_TEST_NUID  = "123456789"
	TEST_RANGE              = 50
)

// A Collection of brute force tests to ensure reliablilty
func TestAlienGenerationEndpointFullE2E(t *testing.T) {
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
	testVerify.AssertArrayLength(services.NUM_WAVES, t)
}
