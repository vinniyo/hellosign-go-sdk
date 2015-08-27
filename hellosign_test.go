package hellosign

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestCreateEmbeddedSignatureRequestNotReturnNil(t *testing.T) {
	mockClient, mockServer := createMockClient("1234")
	defer mockServer.Close()
	embReq := createEmbeddedRequest()
	res, err := mockClient.CreateEmbeddedSignatureRequest(embReq)
	assert.Nil(t, err, "Should not return error")
	assert.NotNil(t, res, "Should return response")
}

// Private Functions

func createMockClient(key string) (Client, *httptest.Server) {
	mockServer := createMockServer(201, "Everything is cool")

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(mockServer.URL)
		},
	}
	mockHTTPClient := &http.Client{Transport: transport}

	client := Client{
		APIKey:     key,
		BaseURL:    mockServer.URL,
		HTTPClient: mockHTTPClient,
	}
	return client, mockServer
}

func createMockServer(status int, body string) *httptest.Server {
	testServer := httptest.NewServer(createMockHandler(status, body))
	return testServer
}

func createMockHandler(status int, body string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, body)
	})
}

func createEmbeddedRequest() EmbeddedRequest {
	fileOne, _ := os.Open("fixtures/offer_letter.pdf")
	fileOne.Close()
	fileTwo, _ := os.Open("fixtures/offer_letter.pdf")
	fileTwo.Close()
	return EmbeddedRequest{
		TestMode: true,
		ClientId: "0987",
		File: []*os.File{
			fileOne,
			fileTwo,
		},
		Title:              "cool title",
		Subject:            "awesome",
		Message:            "cool message bro",
		SigningRedirectURL: "example signing redirect url",
		Signers: []Signer{
			Signer{
				Email: "freddy@hellosign.com",
				Name:  "Freddy Rangel",
			},
			Signer{
				Email: "frederick.rangel@gmail.com",
				Name:  "Frederick Rangel",
			},
		},
		CCEmailAddresses: []string{
			"no@cats.com",
			"no@dogs.com",
		},
		UseTextTags:  false,
		HideTextTags: true,
		Metadata: map[string]string{
			"no":   "cats",
			"more": "dogs",
		},
		FormFieldsPerDocument: [][]DocumentFormField{
			[]DocumentFormField{
				DocumentFormField{
					APIId:    "api_id",
					Name:     "display name",
					Type:     "text",
					X:        123,
					Y:        456,
					Width:    678,
					Required: true,
					Signer:   0,
				},
			},
			[]DocumentFormField{
				DocumentFormField{
					APIId:    "api_id_2",
					Name:     "display name 2",
					Type:     "text 2",
					X:        123,
					Y:        456,
					Width:    678,
					Required: true,
					Signer:   1,
				},
			},
		},
	}
}
