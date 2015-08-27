package hellosign

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
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
	TestMode              bool                  `form_field:"test_mode"`
	ClientId              string                `form_field:"client_id"`
	FileURL               []string              `form_field:"file_url"`
	File                  []*os.File            `form_field:"file"`
	Title                 string                `form_field:"title"`
	Subject               string                `form_field:"subject"`
	Message               string                `form_field:"message"`
	SigningRedirectURL    string                `form_field:"signing_redirect_url"`
	Signers               []Signer              `form_field:"signers"`
	CCEmailAddresses      []string              `form_field:"cc_email_addresses"`
	UseTextTags           bool                  `form_field:"use_text_tags"`
	HideTextTags          bool                  `form_field:"hide_text_tags"`
	Metadata              map[string]string     `form_field:"metadata"`
	FormFieldsPerDocument [][]DocumentFormField `form_field:"form_fields_per_document"`
}

type Signer struct {
	Name  string `field:"name"`
	Email string `field:"email"`
	Order int    `field:"order"`
	Pin   string `field:"pin"`
}

type DocumentFormField struct {
	APIId    string `json:"api_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Required bool   `json:"required"`
	Signer   int    `json:"signer"`
}

func (m *Client) CreateEmbeddedSignatureRequest(
	embeddedRequest EmbeddedRequest) (*http.Response, error) {

	params, writer, err := m.marshalMultipartRequest(embeddedRequest)
	if err != nil {
		return nil, err
	}
	return m.sendEmbeddedSignatureRequest(params, *writer)
}

// Private Methods

func (m *Client) marshalMultipartRequest(
	embRequest EmbeddedRequest) (*bytes.Buffer, *multipart.Writer, error) {

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	structType := reflect.TypeOf(embRequest)
	val := reflect.ValueOf(embRequest)

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		f := valueField.Interface()
		val := reflect.ValueOf(f)
		field := structType.Field(i)
		fieldTag := field.Tag.Get("form_field")

		switch val.Kind() {
		case reflect.Map:
			for k, v := range embRequest.Metadata {
				formField, err := w.CreateFormField(fmt.Sprintf("metadata[%v]", k))
				if err != nil {
					return nil, nil, err
				}
				formField.Write([]byte(v))
			}
		case reflect.Slice:
			switch fieldTag {
			case "signers":
				for i, signer := range embRequest.Signers {
					email, err := w.CreateFormField(fmt.Sprintf("signers[%v][email_address]", i))
					if err != nil {
						return nil, nil, err
					}
					email.Write([]byte(signer.Email))

					name, err := w.CreateFormField(fmt.Sprintf("signers[%v][name]", i))
					if err != nil {
						return nil, nil, err
					}
					name.Write([]byte(signer.Name))

					if signer.Order != 0 {
						order, err := w.CreateFormField(fmt.Sprintf("signers[%v][order]", i))
						if err != nil {
							return nil, nil, err
						}
						order.Write([]byte(strconv.Itoa(signer.Order)))
					}

					if signer.Pin != "" {
						pin, err := w.CreateFormField(fmt.Sprintf("signers[%v][pin]", i))
						if err != nil {
							return nil, nil, err
						}
						pin.Write([]byte(signer.Pin))
					}
				}
			case "cc_email_addresses":
				for k, v := range embRequest.CCEmailAddresses {
					formField, err := w.CreateFormField(fmt.Sprintf("cc_email_addresses[%v]", k))
					if err != nil {
						return nil, nil, err
					}
					formField.Write([]byte(v))
				}
			case "form_fields_per_document":
				formField, err := w.CreateFormField(fieldTag)
				if err != nil {
					return nil, nil, err
				}
				ffpdJSON, err := json.Marshal(embRequest.FormFieldsPerDocument)
				if err != nil {
					return nil, nil, err
				}
				formField.Write([]byte(ffpdJSON))
			case "file":
				for i, file := range embRequest.File {
					formField, err := w.CreateFormFile(fmt.Sprintf("file[%v]", i), file.Name())
					if err != nil {
						return nil, nil, err
					}
					_, err = io.Copy(formField, file)
				}
			case "file_url":
				for i, fileURL := range embRequest.FileURL {
					formField, err := w.CreateFormField(fmt.Sprintf("file_url[%v]", i))
					if err != nil {
						return nil, nil, err
					}
					formField.Write([]byte(fileURL))
				}
			}
		case reflect.Bool:
			formField, err := w.CreateFormField(fieldTag)
			if err != nil {
				return nil, nil, err
			}
			formField.Write([]byte(m.boolToIntString(val.Bool())))
		default:
			if val.String() != "" {
				formField, err := w.CreateFormField(fieldTag)
				if err != nil {
					return nil, nil, err
				}
				formField.Write([]byte(val.String()))
			}
		}
	}

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
