package hellosign

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	baseURL string = "https://api.hellosign.com/v3/"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

type EmbeddedRequest struct {
	FileUrl string
	Subject string
	Message string
	// TODO: change this to a struct rather than map
	Signers  []map[string]string
	TestMode int
}

func (m *Client) CreateEmbeddedSignatureRequest(request EmbeddedRequest, clientId string) (*http.Response, error) {
	endpoint := m.GetEndpoint()

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

	client := m.GetHTTPClient()
	apiCall, _ := http.NewRequest("POST", endpoint, &b)
	apiCall.Header.Add("Content-Type", w.FormDataContentType())
	if err != nil {
		return nil, err
	}
	apiCall.SetBasicAuth(m.APIKey, "")

	response, err := client.Do(apiCall)
	defer response.Body.Close()
	return response, err
}

func (m *Client) GetEndpoint() string {
	var url string
	if m.BaseURL != "" {
		url = m.BaseURL
	} else {
		url = baseURL + "signature_request/create_embedded"
	}
	return url
}

func (m *Client) GetHTTPClient() *http.Client {
	var http_client *http.Client
	if m.HTTPClient != nil {
		http_client = m.HTTPClient
	} else {
		http_client = &http.Client{}
	}
	return http_client
}
