package handlers

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pkg/errors"
)

// mustNewServer is a test utility function to initialize a new server
// with the given app paratmers
func mustNewServer(t *testing.T, app *App) *httptest.Server {
	app.WebURL = os.Getenv("WebURL")

	r, err := NewRouter(app)
	if err != nil {
		t.Fatal(errors.Wrap(err, "initializing server"))
	}

	server := httptest.NewServer(r)

	return server
}
