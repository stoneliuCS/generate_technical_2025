package utils

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestClient struct {
	client  http.Client
	baseurl string
}

type TestVerify struct {
	res *http.Response
}

func CreateTestClient(port string) TestClient {
	return TestClient{
		client:  http.Client{Timeout: 30 * time.Second},
		baseurl: "http://localhost:" + port,
	}
}

// Blocks until the server is ready.
func (t TestClient) CheckServer(timeout time.Duration) bool {
	start := time.Now()
	for {
		res, err := t.client.Get(t.baseurl + "/healthcheck")
		// Handle error case
		if err != nil {
			// Check timeout before sleeping
			if time.Since(start) > timeout {
				return false
			}
			time.Sleep(100 * time.Millisecond)
			continue // Skip to next iteration
		}

		// Always close response body to prevent resource leaks
		defer res.Body.Close()

		// Check if successful
		if res.StatusCode == 200 {
			return true
		}

		// Check timeout
		if time.Since(start) > timeout {
			return false
		}

		// Sleep before next attempt
		time.Sleep(100 * time.Millisecond)
	}
}

func (t TestClient) GET(endpoint string) TestVerify {
	fn := func() (*http.Response, error) { return t.client.Get(t.baseurl + endpoint) }
	res := SafeCall(fn)
	return TestVerify{res: res}
}

func (v TestVerify) AssertStatusCode(statusCode int, t *testing.T) TestVerify {
	assert.Equal(t, 200, statusCode)
	return v
}

func (v TestVerify) AssertBody(expected any, t *testing.T) {
	defer v.res.Body.Close()

	var actual any
	err := json.NewDecoder(v.res.Body).Decode(&actual)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
