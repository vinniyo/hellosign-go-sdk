package hellosign

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestCreateEmbeddedSigningRequest(t *testing.T) {
}

func createMockServer(status int, body string) *httptest.Server {
	mockRoute := http.handleFunc(func(w http.ResponseWriter, r *http.Response) {
		w.WriteHeader(status)
		w.Header().set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	})
	testServer := httptest.NewServer(mockRoute)
	defer testServer.Close()
}
