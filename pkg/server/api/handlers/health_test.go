package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestCheckHealth(t *testing.T) {
	defer testutils.ClearData()

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	// Execute
	req := testutils.MakeReq(server, "GET", "/health", "")
	res := testutils.HTTPDo(t, req)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "Status code mismtach")
}
