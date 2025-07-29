package integrationtests

import (
	"testing"

	"github.com/google/uuid"
)

func TestUserWithNonValidNUIDReceives400(t *testing.T) {
	client := CLIENT.AddBody(map[string]any{
		"email": "notavalidemail@gmail.com",
		"nuid":  "1231",
	}).AddHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	testVerify := client.POST("/api/v1/member/register")
	testVerify.AssertStatusCode(400, t).AssertBody(map[string]any{
		"message": "Not a valid northeastern email address.",
	}, t)
}

func TestUserWithBadNUIDReceives400(t *testing.T) {
	client := CLIENT.AddBody(map[string]any{
		"email": "somebody@northeastern.edu",
		"nuid":  "1231",
	}).AddHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	testVerify := client.POST("/api/v1/member/register")
	testVerify.AssertStatusCode(400, t).AssertBody(map[string]any{
		"message": "Not a valid NUID.",
	}, t)
}

func TestUserReceives201OnGoodRequest(t *testing.T) {
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
		_, err := uuid.Parse(s)
		return err == nil
	}
	testVerify.AssertStatusCode(201, t).AssertProperty("id", pred, t)
}

func TestMemberCannotRegisterTwice(t *testing.T) {
	client := CLIENT.AddBody(map[string]any{
		"email": "hasneverregisteredbefore@northeastern.edu",
		"nuid":  "123456789", // NUID is 9 characters long
	}).AddHeaders(map[string]string{
		"Content-Type": "application/json",
	})

	firstRegistration := client.POST("/api/v1/member/register")
	firstRegistration.AssertStatusCode(201, t)
	nextClient := CLIENT.AddBody(map[string]any{
		"email": "hasneverregisteredbefore@northeastern.edu",
		"nuid":  "123456789", // NUID is 9 characters long
	}).AddHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	testVerify := nextClient.POST("/api/v1/member/register")
	testVerify.AssertStatusCode(409, t)
}

func TestMemberGets200IfUserIsFoundInDatabase(t *testing.T) {
	client := CLIENT
	testVerifyGET := client.GET("/api/v1/member?email=somebody%40northeastern.edu&nuid=123456789")
	pred := func(prop any) bool {
		s, ok := prop.(string)
		if !ok {
			return ok
		}
		_, err := uuid.Parse(s)
		return err == nil
	}
	testVerifyGET.AssertStatusCode(200, t).AssertProperty("id", pred, t)
}

func TestMemberGets400ForMalformedNUIDOrEmail(t *testing.T) {
	client := CLIENT
	testVerifyGET := client.GET("/api/v1/member?email=somebody%40gmail.com&nuid=2134")

	testVerifyGET.AssertStatusCode(400, t)
	testVerifyGET = client.GET("/api/v1/member?email=somebody%40northeastern.com&nuid=1234")
	testVerifyGET.AssertStatusCode(400, t)
}

func TestMemberGets404IfNotFound(t *testing.T) {
	client := CLIENT
	testVerifyGET := client.GET("/api/v1/member?email=somebodyNotExist%40northeastern.edu&nuid=123456789")
	testVerifyGET.AssertStatusCode(404, t)
}
