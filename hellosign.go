package hellosign

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	baseUrl string = "https://api.hellosign.com/v3"
)

type HelloSign struct {
	APIKey string
}

type EmbeddedRequest struct {
	FileUrl  string
	Subject  string
	Message  string
	Signers  []map[string]string
	TestMode int
}

func New(apiKey string) *HelloSign {
	return &HelloSign{APIKey: apiKey}
}

func (hs *HelloSign) CreateEmbeddedSignatureRequest(request EmbeddedRequest, clientId string) (*http.Response, error) {
	endpoint := "https://api.hellosign.com/v3/signature_request/create_embedded"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fileUrl, err := w.CreateFormField("file_url[]")
	if err != nil {
		return nil, err
	}
	fileUrl.Write([]byte(request.FileUrl))

	client_id, err := w.CreateFormField("client_id")
	if err != nil {
		return nil, err
	}
	client_id.Write([]byte(clientId))

	subject, err := w.CreateFormField("subject")
	if err != nil {
		return nil, err
	}
	subject.Write([]byte(request.Subject))

	message, err := w.CreateFormField("message")
	if err != nil {
		return nil, err
	}
	message.Write([]byte(request.Subject))

	email, err := w.CreateFormField("signers[0][email_address]")
	if err != nil {
		return nil, err
	}
	email.Write([]byte(request.Signers[0]["email"]))

	name, err := w.CreateFormField("signers[0][name]")
	if err != nil {
		return nil, err
	}
	name.Write([]byte(request.Signers[0]["name"]))

	testMode, err := w.CreateFormField("test_mode")
	if err != nil {
		return nil, err
	}
	testMode.Write([]byte(strconv.Itoa(request.TestMode)))

	w.Close()

	client := &http.Client{}
	apiCall, _ := http.NewRequest("POST", endpoint, &b)
	apiCall.Header.Add("Content-Type", w.FormDataContentType())
	if err != nil {
		return nil, err
	}
	apiCall.SetBasicAuth(hs.APIKey, "")

	response, err := client.Do(apiCall)
	defer response.Body.Close()
	return response, err
}
