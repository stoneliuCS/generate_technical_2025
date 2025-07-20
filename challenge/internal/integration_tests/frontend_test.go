package integrationtests

import (
	"fmt"
	"generate_technical_challenge_2025/internal/services"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	UUID                        = uuid.New()
	invalidUUID                 = "not-a-uuid"
	testMemberAlreadyRegistered = false
	memberUUID                  = uuid.Nil
)

const (
	DEFAULT_FRONTEND_ALIEN_COUNT = 10
)

func TestMemberUUIDInvalidReceives400(t *testing.T) {
	client := CLIENT.AddHeaders(map[string]string{
		"Content-Type": "application/json",
	})

	expectedErrorMessage :=
		fmt.Sprintf(
			`operation APIV1ChallengeFrontendIDAliensGet: decode params: path: "id": invalid UUID length: %d`,
			len(invalidUUID))

	testVerify := client.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens", invalidUUID))
	testVerify.AssertStatusCode(400, t).AssertBody(map[string]any{
		"error_message": expectedErrorMessage,
	}, t)
}

func TestMemberUUIDNotFoundReceives404(t *testing.T) {
	client := CLIENT.AddHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	testVerify := client.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens", UUID.String()))
	testVerify.AssertStatusCode(404, t).AssertBody(map[string]any{
		"message": "Unable to find member id.",
	}, t)
}

// Idempotent function to register--call it at the beginning of any testing function in this file
// that intends to use a protected endpoint.
func registerForChallenge(t *testing.T) {
	if testMemberAlreadyRegistered {
		return
	}
	client := CLIENT.AddBody(map[string]any{
		"email": "somefrontendperson@northeastern.edu",
		"nuid":  "123456789",
	}).AddHeaders(map[string]string{
		"Content-Type": "application/json",
	})

	testVerify := client.POST("/api/v1/member/register")
	pred := func(prop any) bool {
		s, ok := prop.(string)
		if !ok {
			return ok
		}

		// Assign the global memberUUID for future use in testing.
		var err error
		memberUUID, err = uuid.Parse(s)
		return err == nil
	}

	// Make sure that the member is successfully registered before preceding in test suite.
	testVerify.AssertStatusCode(201, t).AssertProperty("id", pred, t)
	testMemberAlreadyRegistered = true
}

func TestLimitInvalidReceives400(t *testing.T) {
	registerForChallenge(t)

	invalidLimit := -3
	testVerify := CLIENT.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens?limit=%d", memberUUID.String(), invalidLimit))
	testVerify.AssertStatusCode(400, t)
}

func TestOffsetInvalidReceives400(t *testing.T) {
	registerForChallenge(t)

	invalidOffset := -3
	testVerify := CLIENT.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens?offset=%d", memberUUID.String(), invalidOffset))
	testVerify.AssertStatusCode(400, t)
}

func TestDefaultLimitAndOffsetFrontend(t *testing.T) {
	registerForChallenge(t)

	testVerify := CLIENT.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens", memberUUID.String()))
	testVerify.AssertStatusCode(200, t).AssertArrayLength(DEFAULT_FRONTEND_ALIEN_COUNT, t)
}

func TestCustomLimitDefaultOffsetFrontend(t *testing.T) {
	registerForChallenge(t)

	validLimit := 25
	testVerify := CLIENT.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens?limit=%d", memberUUID.String(), validLimit))
	testVerify.AssertStatusCode(200, t).AssertArrayLengthBetween(0, validLimit, t)
}

func TestDefaultLimitCustomOffsetFrontend(t *testing.T) {
	registerForChallenge(t)

	offset := 3

	noOffset := CLIENT.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens", memberUUID.String()))
	noOffset.AssertStatusCode(200, t).AssertArrayLengthBetween(
		0, services.UPPER_ALIEN_AMOUNT, t)

	withOffset := CLIENT.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens?offset=%d", memberUUID.String(), offset))
	withOffset.AssertStatusCode(200, t).AssertArrayLengthBetween(
		0, services.UPPER_ALIEN_AMOUNT-offset, t)
}

func TestCustomLimitCustomOffsetFrontend(t *testing.T) {
	registerForChallenge(t)

	offset := 3
	limit := 8

	noOffset := CLIENT.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens?limit=%d", memberUUID.String(), limit))
	noOffset.AssertStatusCode(200, t).AssertArrayLengthBetween(
		0, limit, t)

	withOffset := CLIENT.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens?limit=%d&offset=%d", memberUUID.String(), limit, offset))
	withOffset.AssertStatusCode(200, t).AssertArrayLengthBetween(
		0, limit, t)

	limit = 2

	offsetBiggerThanLimit := CLIENT.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens?limit=%d&offset=%d", memberUUID.String(), limit, offset))
	offsetBiggerThanLimit.AssertStatusCode(200, t).AssertArrayLength(limit, t)
}

func TestOffsetShiftsAliens(t *testing.T) {
	registerForChallenge(t)

	firstBatch := CLIENT.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens", memberUUID.String()))
	firstAliens := firstBatch.GetBodyAsArray(t)

	offsetBatch := CLIENT.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens?offset=2", memberUUID.String()))
	offsetAliens := offsetBatch.GetBodyAsArray(t)

	// Alien #3 from first batch should be alien #1 in offset batch.
	if len(firstAliens) > 2 && len(offsetAliens) > 0 {
		alien3 := firstAliens[2].(map[string]any)["id"]
		alien1Offset := offsetAliens[0].(map[string]any)["id"]
		assert.Equal(t, alien3, alien1Offset)
	}
}
