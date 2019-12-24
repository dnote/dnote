package operations

import (
	// "fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestCreateDigest(t *testing.T) {
	t.Run("no previous digest", func(t *testing.T) {
		defer testutils.ClearData()

		db := testutils.DB

		user := testutils.SetupUserData()
		rule := database.RepetitionRule{UserID: user.ID}
		testutils.MustExec(t, testutils.DB.Save(&rule), "preparing rule")

		result, err := CreateDigest(db, rule, nil)
		if err != nil {
			t.Fatal(errors.Wrap(err, "performing"))
		}

		assert.Equal(t, result.Version, 1, "Version mismatch")
	})

	t.Run("with previous digest", func(t *testing.T) {
		defer testutils.ClearData()

		db := testutils.DB

		user := testutils.SetupUserData()
		rule := database.RepetitionRule{UserID: user.ID}
		testutils.MustExec(t, testutils.DB.Save(&rule), "preparing rule")

		d := database.Digest{UserID: user.ID, RuleID: rule.ID, Version: 8}
		testutils.MustExec(t, testutils.DB.Save(&d), "preparing digest")

		result, err := CreateDigest(db, rule, nil)
		if err != nil {
			t.Fatal(errors.Wrap(err, "performing"))
		}

		assert.Equal(t, result.Version, 9, "Version mismatch")
	})
}
