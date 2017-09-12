package hellosign

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	ClientID              string                `form_field:"client_id"`
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

type SignatureRequestResponse struct {
	SignatureRequest *SignatureRequest `json:"signature_request"`
}

type SignatureRequest struct {
	TestMode              bool            `json:"test_mode"`               // Whether this is a test signature request. Test requests have no legal value. Defaults to 0.
	SignatureRequestID    string          `json:"signature_request_id"`    // The id of the SignatureRequest.
	RequesterEmailAddress string          `json:"requester_email_address"` // The email address of the initiator of the SignatureRequest.
	Title                 string          `json:"title"`                   // The title the specified Account uses for the SignatureRequest.
	Subject               string          `json:"subject"`                 // The subject in the email that was initially sent to the signers.
	Message               string          `json:"message"`                 // The custom message in the email that was initially sent to the signers.
	IsComplete            bool            `json:"is_complete"`             // Whether or not the SignatureRequest has been fully executed by all signers.
	IsDeclined            bool            `json:"is_declined"`             // Whether or not the SignatureRequest has been declined by a signer.
	HasError              bool            `json:"has_error"`               // Whether or not an error occurred (either during the creation of the SignatureRequest or during one of the signings).
	FilesURL              string          `json:"files_url"`               // The URL where a copy of the request's documents can be downloaded.
	SigningURL            string          `json:"signing_url"`             // The URL where a signer, after authenticating, can sign the documents. This should only be used by users with existing HelloSign accounts as they will be required to log in before signing.
	DetailsURL            string          `json:"details_url"`             // The URL where the requester and the signers can view the current status of the SignatureRequest.
	CCEmailAddress        []*string       `json:"cc_email_addresses"`      // A list of email addresses that were CCed on the SignatureRequest. They will receive a copy of the final PDF once all the signers have signed.
	SigningRedirectURL    string          `json:"signing_redirect_url"`    // The URL you want the signer redirected to after they successfully sign.
	CustomFields          []*CustomField  `json:"custom_fields"`           // An array of Custom Field objects containing the name and type of each custom field.
	ResponseDsata         []*ResponseData `json:"response_data"`           // An array of form field objects containing the name, value, and type of each textbox or checkmark field filled in by the signers.
	Signatures            []*Signature    `json:"signatures"`              // An array of signature objects, 1 for each signer.
	Warnings              []*Warning      `json:"warnings"`                // An array of warning objects.
}

type CustomField struct {
	Name     string `json:"name"`     // The name of the Custom Field.
	Type     string `json:"type"`     // The type of this Custom Field. Only 'text' and 'checkbox' are currently supported.
	Value    string `json:"value"`    // A text string for text fields or true/false for checkbox fields
	Required bool   `json:"required"` // A boolean value denoting if this field is required.
	Editor   string `json:"editor"`   // The name of the Role that is able to edit this field.
}

type ResponseData struct {
	ApiID       string `json:"api_id"`       // The unique ID for this field.
	SignatureID string `json:"signature_id"` // The ID of the signature to which this response is linked.
	Name        string `json:"name"`         // The name of the form field.
	Value       string `json:"value"`        // The value of the form field.
	Required    bool   `json:"required"`     // A boolean value denoting if this field is required.
	Type        string `json:"type"`         // The type of this form field. See field types
}

type Signature struct {
	SignatureID        string `json:"signature_id"`         // Signature identifier.
	SignerEmailAddress string `json:"signer_email_address"` // The email address of the signer.
	SignerName         string `json:"signer_name"`          // The name of the signer.
	Order              int    `json:"order"`                // If signer order is assigned this is the 0-based index for this signer.
	StatusCode         string `json:"status_code"`          // The current status of the signature. eg: awaiting_signature, signed, declined
	DeclineReason      string `json:"decline_reason"`       // The reason provided by the signer for declining the request.
	SignatedAt         string `json:"signed_at"`            // Time that the document was signed or null.
	LastViewedAt       string `json:"last_viewed_at"`       //The time that the document was last viewed by this signer or null.
	LastRemindedAt     string `json:"last_reminded_at"`     //The time the last reminder email was sent to the signer or null.
	HasPin             bool   `json:"has_pin"`              // Boolean to indicate whether this signature requires a PIN to access.
}

type Warning struct {
	Message string `json:"warning_msg"`
	Name    string `json:"warning_name"`
}

func (m *Client) CreateEmbeddedSignatureRequest(
	embeddedRequest EmbeddedRequest) (*SignatureRequest, error) {

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
					if err != nil {
						log.Fatal(err)
					}
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
	params *bytes.Buffer, w multipart.Writer) (*SignatureRequest, error) {
	endpoint := fmt.Sprintf("%s%s", m.getEndpoint(), "signature_request/create_embedded")
	log.Println(endpoint)
	request, _ := http.NewRequest("POST", endpoint, params)
	request.Header.Add("Content-Type", w.FormDataContentType())
	request.SetBasicAuth(m.APIKey, "")

	response, err := m.getHTTPClient().Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	sigRequestResponse := &SignatureRequestResponse{}

	err = json.NewDecoder(response.Body).Decode(sigRequestResponse)

	sigRequest := sigRequestResponse.SignatureRequest

	return sigRequest, err
}

func (m *Client) getEndpoint() string {
	var url string
	if m.BaseURL != "" {
		url = m.BaseURL
	} else {
		url = baseURL
	}
	return url
}

func (m *Client) getHTTPClient() *http.Client {
	var httpClient *http.Client
	if m.HTTPClient != nil {
		httpClient = m.HTTPClient
	} else {
		httpClient = &http.Client{}
	}
	return httpClient
}

func (m *Client) boolToIntString(value bool) string {
	if value == true {
		return "1"
	}
	return "0"
}
