package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestGetDigest(t *testing.T) {
	defer testutils.ClearData()

	// Setup
	server := MustNewServer(t, nil)
	defer server.Close()

	owner := testutils.SetupUserData()
	nonOwner := testutils.SetupUserData()
	digest := database.Digest{
		UserID: owner.ID,
	}
	testutils.MustExec(t, testutils.DB.Save(&digest), "preparing digest")

	t.Run("owner", func(t *testing.T) {
		// Execute
		req := testutils.MakeReq(server.URL, "GET", fmt.Sprintf("/digests/%s", digest.UUID), "")
		res := testutils.HTTPAuthDo(t, req, owner)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")
	})

	t.Run("non owner", func(t *testing.T) {
		// Execute
		req := testutils.MakeReq(server.URL, "GET", fmt.Sprintf("/digests/%s", digest.UUID), "")
		res := testutils.HTTPAuthDo(t, req, nonOwner)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusNotFound, "")
	})

	t.Run("guest", func(t *testing.T) {
		// Execute
		req := testutils.MakeReq(server.URL, "GET", fmt.Sprintf("/digests/%s", digest.UUID), "")
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")
	})
}

func TestGetDigests(t *testing.T) {
	defer testutils.ClearData()

	// Setup
	server := MustNewServer(t, nil)
	defer server.Close()

	user := testutils.SetupUserData()
	digest := database.Digest{
		UserID: user.ID,
	}
	testutils.MustExec(t, testutils.DB.Save(&digest), "preparing digest")

	t.Run("user", func(t *testing.T) {
		// Execute
		req := testutils.MakeReq(server.URL, "GET", "/digests", "")
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusOK, "")
	})

	t.Run("guest", func(t *testing.T) {
		// Execute
		req := testutils.MakeReq(server.URL, "GET", fmt.Sprintf("/digests/%s", digest.UUID), "")
		res := testutils.HTTPDo(t, req)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusUnauthorized, "")
	})
}
