package token

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		kind string
	}{
		{
			kind: database.TokenTypeEmailPreference,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("token type %s", tc.kind), func(t *testing.T) {
			defer testutils.ClearData()

			// Set up
			u := testutils.SetupUserData()

			// Execute
			tok, err := Create(testutils.DB, u.ID, tc.kind)
			if err != nil {
				t.Fatal(errors.Wrap(err, "performing"))
			}

			// Test
			var count int
			testutils.MustExec(t, testutils.DB.Model(&database.Token{}).Count(&count), "counting token")
			assert.Equalf(t, count, 1, "error mismatch")

			var tokenRecord database.Token
			testutils.MustExec(t, testutils.DB.First(&tokenRecord), "finding token")
			assert.Equalf(t, tokenRecord.UserID, tok.UserID, "UserID mismatch")
			assert.Equalf(t, tokenRecord.Value, tok.Value, "Value mismatch")
			assert.Equalf(t, tokenRecord.Type, tok.Type, "Type mismatch")
		})
	}
}
