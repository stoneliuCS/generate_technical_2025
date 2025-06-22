package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestClient struct {
	client  http.Client
	baseurl string
	logger  *slog.Logger
}

type TestVerify struct {
	res *http.Response
}

func CreateTestClient(port int, logger *slog.Logger) TestClient {
	portString := fmt.Sprintf("%d", port)
	return TestClient{
		client:  http.Client{Timeout: 30 * time.Second},
		baseurl: "http://localhost:" + portString,
		logger:  logger,
	}
}

// Blocks until the server is ready.
func (t TestClient) CheckServer(timeout time.Duration) bool {
	start := time.Now()
	for {
		t.logger.Info("Attempting to connect to test backend server...")
		res, err := t.client.Get(t.baseurl + "/healthcheck")
		if err != nil {
			if time.Since(start) > timeout {
				t.logger.Info("Attempting to connect and ran out of timeout.")
				return false
			}
			time.Sleep(100 * time.Millisecond)
			t.logger.Info("Retrying to connect...")
			continue // Skip to next iteration
		}

		defer res.Body.Close()

		if res.StatusCode == 200 {
			return true
		}

		if time.Since(start) > timeout {
			return false
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (t TestClient) GET(endpoint string) TestVerify {
	fn := func() (*http.Response, error) { return t.client.Get(t.baseurl + endpoint) }
	res := SafeCall(fn)
	return TestVerify{res: res}
}

func (v TestVerify) AssertStatusCode(statusCode int, t *testing.T) TestVerify {
	assert.Equal(t, statusCode, v.res.StatusCode)
	return v
}

func (v TestVerify) AssertBody(expected any, t *testing.T) {
	defer v.res.Body.Close()

	var actual any
	err := json.NewDecoder(v.res.Body).Decode(&actual)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
