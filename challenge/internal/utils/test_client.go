package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	headers map[string]string
	body    io.Reader
}

type TestVerify struct {
	res *http.Response
}

func CreateTestClient(port int, logger *slog.Logger) *TestClient {
	portString := fmt.Sprintf("%d", port)
	return &TestClient{
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
	fn := func(body io.Reader) (*http.Request, error) {
		return http.NewRequest("GET", t.baseurl+endpoint, body)
	}
	return t.internalWrapper(fn)
}

func (t TestClient) POST(endpoint string) TestVerify {
	fn := func(body io.Reader) (*http.Request, error) {
		return http.NewRequest("POST", t.baseurl+endpoint, body)
	}
	return t.internalWrapper(fn)
}

func (t TestClient) internalWrapper(reqSupplier func(body io.Reader) (*http.Request, error)) TestVerify {
	var req *http.Request
	if t.body != nil {
		req = FatalCall(func() (*http.Request, error) { return reqSupplier(t.body) })
	} else {
		req = FatalCall(func() (*http.Request, error) { return reqSupplier(nil) })
	}
	if t.headers != nil {
		for key, val := range t.headers {
			req.Header.Add(key, val)
		}
	}
	resSupplier := func() (*http.Response, error) {
		return t.client.Do(req)
	}
	res := FatalCall(resSupplier)
	return TestVerify{res: res}
}

func (t *TestClient) AddHeaders(headers map[string]string) {
	t.headers = headers
}

func (t *TestClient) AddBody(body map[string]any) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	bodyReader := bytes.NewReader(jsonBody)
	t.body = bodyReader
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

// Asserts for some property returned by the API, it holds.
func (v TestVerify) AssertProperty(propertyName string, pred func(s any) bool, t *testing.T) {
	defer v.res.Body.Close()

	var actual map[string]any
	err := json.NewDecoder(v.res.Body).Decode(&actual)
	require.NoError(t, err)
	val := actual[propertyName]
	assert.True(t, pred(val))
}
