package hellosign

import (
	"bytes"
	"mime/multipart"
	"net/http"
)

const (
	baseURL string = "https://api.hellosign.com/v3/"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

type Signer struct {
	name  string
	email string
	order int
	pin   string
}

type EmbeddedRequest struct {
	// TODO: change arrays of maps to arrays of structs and add struct tags
	TestMode              bool   `field:"test_mode"`
	ClientId              string `field:"client_id"`
	FileURL               string `field:"file_url[]"`
	Title                 string `field="title"`
	Subject               string `field:"subject"`
	Message               string `field="message"`
	SigningRedirectURL    string `field="signing_redirect_url"`
	Signers               []Signer
	CCEmailAddress        string `field="cc_email_address"`
	UseTextTags           bool   `field="use_text_tags"`
	HideTextTags          bool   `field="hide_text_tags"`
	Metadata              []map[string]string
	FormFieldsPerDocument []map[string]string
}

func (m *Client) CreateEmbeddedSignatureRequest(
	embeddedRequest EmbeddedRequest) (*http.Response, error) {

	params, writer, err := m.marshalMultipartRequest(embeddedRequest)
	if err != nil {
		return nil, err
	}
	return m.sendEmbeddedSignatureRequest(params, *writer)
}

func (m *Client) marshalMultipartRequest(
	request EmbeddedRequest) (*bytes.Buffer, *multipart.Writer, error) {

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fileUrl, err := w.CreateFormField("file_url[]")
	if err != nil {
		return nil, nil, err
	}
	fileUrl.Write([]byte(request.FileURL))

	client_id, err := w.CreateFormField("client_id")
	if err != nil {
		return nil, nil, err
	}
	client_id.Write([]byte(request.ClientId))

	subject, err := w.CreateFormField("subject")
	if err != nil {
		return nil, nil, err
	}
	subject.Write([]byte(request.Subject))

	message, err := w.CreateFormField("message")
	if err != nil {
		return nil, nil, err
	}
	message.Write([]byte(request.Subject))

	email, err := w.CreateFormField("signers[0][email_address]")
	if err != nil {
		return nil, nil, err
	}
	email.Write([]byte(request.Signers[0].email))

	name, err := w.CreateFormField("signers[0][name]")
	if err != nil {
		return nil, nil, err
	}
	name.Write([]byte(request.Signers[0].name))

	testMode, err := w.CreateFormField("test_mode")
	if err != nil {
		return nil, nil, err
	}
	testMode.Write([]byte(m.boolToIntString(request.TestMode)))

	w.Close()
	return &b, w, nil
}

func (m *Client) sendEmbeddedSignatureRequest(
	params *bytes.Buffer, w multipart.Writer) (*http.Response, error) {
	endpoint := m.getEndpoint()
	request, _ := http.NewRequest("POST", endpoint, params)
	request.Header.Add("Content-Type", w.FormDataContentType())
	request.SetBasicAuth(m.APIKey, "")

	response, err := m.getHTTPClient().Do(request)
	defer response.Body.Close()
	return response, err
}

func (m *Client) getEndpoint() string {
	var url string
	if m.BaseURL != "" {
		url = m.BaseURL
	} else {
		url = baseURL + "signature_request/create_embedded"
	}
	return url
}

func (m *Client) getHTTPClient() *http.Client {
	var http_client *http.Client
	if m.HTTPClient != nil {
		http_client = m.HTTPClient
	} else {
		http_client = &http.Client{}
	}
	return http_client
}

func (m *Client) boolToIntString(value bool) string {
	if value == true {
		return "1"
	} else {
		return "0"
	}
}
