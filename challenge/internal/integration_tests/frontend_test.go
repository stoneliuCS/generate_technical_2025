package integrationtests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

var (
	UUID        = uuid.New()
	invalidUUID = "not-a-uuid"
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
