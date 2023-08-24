package gorp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRPClient(t *testing.T) {
	t.Parallel()
	client := NewClient("http://host.com", "prj", "uuid")

	assert.Equal(t, "prj", client.project)
	assert.Equal(t, "http://host.com", client.http.HostURL)
	assert.Equal(t, "uuid", client.http.Token)
}

func TestHandleErrors(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClient(server.URL, "prj", "uuid")
	_, err := client.GetLaunches()
	assert.Error(t, err)
}
