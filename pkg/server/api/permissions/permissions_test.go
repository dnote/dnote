package permissions

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func init() {
	testutils.InitTestDB()
}

func TestViewNote(t *testing.T) {
	user := testutils.SetupUserData()
	anotherUser := testutils.SetupUserData()

	db := database.DBConn
	defer testutils.ClearData()

	b1 := database.Book{
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b1), "preparing b1")

	privateNote := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "privateNote content",
		Deleted:  false,
		Public:   false,
	}

	publicNote := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
		Body:     "privateNote content",
		Deleted:  false,
		Public:   true,
	}
	testutils.MustExec(t, db.Save(&privateNote), "preparing privateNote")

	t.Run("owner accessing private note", func(t *testing.T) {
		result := ViewNote(&user, privateNote)
		assert.Equal(t, result, true, "result mismatch")
	})

	t.Run("owner accessing public note", func(t *testing.T) {
		result := ViewNote(&user, publicNote)
		assert.Equal(t, result, true, "result mismatch")
	})

	t.Run("non-owner accessing private note", func(t *testing.T) {
		result := ViewNote(&anotherUser, privateNote)
		assert.Equal(t, result, false, "result mismatch")
	})

	t.Run("non-owner accessing public note", func(t *testing.T) {
		result := ViewNote(&anotherUser, publicNote)
		assert.Equal(t, result, true, "result mismatch")
	})

	t.Run("guest accessing private note", func(t *testing.T) {
		result := ViewNote(nil, privateNote)
		assert.Equal(t, result, false, "result mismatch")
	})

	t.Run("guest accessing public note", func(t *testing.T) {
		result := ViewNote(nil, publicNote)
		assert.Equal(t, result, true, "result mismatch")
	})
}
