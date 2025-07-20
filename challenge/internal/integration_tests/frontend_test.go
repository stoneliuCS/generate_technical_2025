package integrationtests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
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

func registerForChallenge(t *testing.T) {
	if testMemberAlreadyRegistered {
		return
	}
	client := CLIENT.AddBody(map[string]any{
		"email": "somebody@northeastern.edu",
		"nuid":  "123456789", // NUID is 9 characters long
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

	// Make sure that the member is successfully registere before preceding in test suite.
	testVerify.AssertStatusCode(201, t).AssertProperty("id", pred, t)
	testMemberAlreadyRegistered = true
}

func TestDefaultLimitAndOffsetFrontend(t *testing.T) {
	registerForChallenge(t)

	testVerify := CLIENT.GET(fmt.Sprintf("/api/v1/challenge/frontend/%s/aliens", memberUUID.String()))
	testVerify.AssertStatusCode(200, t).AssertArrayLength(DEFAULT_FRONTEND_ALIEN_COUNT, t)
}
