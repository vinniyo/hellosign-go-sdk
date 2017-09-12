package hellosign

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func TestCreateEmbeddedSignatureRequestNotReturnNil(t *testing.T) {
	// Start our recorder
	vcr := fixture("fixtures/embedded_signature_request")
	defer vcr.Stop() // Make sure recorder is stopped once done with it

	client := createVcrClient(vcr)

	embReq := createEmbeddedRequest()
	res, err := client.CreateEmbeddedSignatureRequest(embReq)

	assert.NotNil(t, res, "Should return response")
	assert.Nil(t, err, "Should not return error")

	assert.Equal(t, "6d7ad140141a7fe6874fec55931c363e0301c353", res.SignatureRequestID)
	assert.Equal(t, "awesome", res.Subject)
	assert.Equal(t, true, res.TestMode)
	assert.Equal(t, false, res.IsComplete)
	assert.Equal(t, false, res.IsDeclined)
}

// Private Functions

func fixture(path string) *recorder.Recorder {
	vcr, err := recorder.New(path)
	if err != nil {
		log.Fatal(err)
	}
	return vcr
}

func createVcrClient(transport *recorder.Recorder) Client {
	httpClient := &http.Client{Transport: transport}

	client := Client{
		APIKey:     os.Getenv("HELLO_SIGN_API_KEY"),
		HTTPClient: httpClient,
	}
	return client
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
	fileTwo, _ := os.Open("fixtures/offer_letter.pdf")

	return EmbeddedRequest{
		TestMode: true,
		ClientID: os.Getenv("HELLOSIGN_CLIENT_ID"),
		File: []*os.File{
			fileOne,
			fileTwo,
		},
		Title:   "cool title",
		Subject: "awesome",
		Message: "cool message bro",
		// SigningRedirectURL: "example signing redirect url",
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
					Type:     "text",
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
