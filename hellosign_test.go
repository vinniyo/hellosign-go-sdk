package hellosign

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCreateEmbeddedSigningRequest(t *testing.T) {
	mockClient := createMockClient("1234")

	// Create new embedded request struct
	embReq := EmbeddedRequest{
		ClientId: "0987",
		FileURL:  "matrix",
		Subject:  "awesome",
		Message:  "cool message bro",
		Signers: []map[string]string{
			{
				"email": "freddy@hellosign.com",
				"name":  "Freddy Rangel",
			},
		},
		TestMode: true,
	}
	// Call #CreateEmdeddedSignatureRequest on client struct
	res, err := mockClient.CreateEmbeddedSignatureRequest(embReq)
	assert.Nil(t, err, "Should not return error")
	assert.NotNil(t, res, "Should return response")
}

func createMockClient(key string) Client {
	fake_server := createMockServer(200, "Everything is cool")
	defer fake_server.Close()

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(fake_server.URL)
		},
	}
	mockHTTPClient := &http.Client{Transport: transport}

	client := Client{
		APIKey:     key,
		BaseURL:    fake_server.URL,
		HTTPClient: mockHTTPClient,
	}
	return client
}

func createMockServer(status int, body string) *httptest.Server {
	testServer := httptest.NewServer(createMockHandler(status, body))
	return testServer
}

func createMockHandler(status int, _ string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "Meow")
	})
}
