package repetition

import (
	"sort"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func init() {
	testutils.InitTestDB()
}

func assertLastActive(t *testing.T, ruleUUID string, lastActive int64) {
	db := database.DBConn

	var rule database.RepetitionRule
	testutils.MustExec(t, db.Where("uuid = ?", ruleUUID).First(&rule), "finding rule1")

	assert.Equal(t, rule.LastActive, lastActive, "LastActive mismatch")
}

func assertRepetitionCount(t *testing.T, rule database.RepetitionRule, expected int) {
	db := database.DBConn

	var digestCount int
	testutils.MustExec(t, db.Model(&database.Digest{}).Where("rule_id = ? AND user_id = ?", rule.ID, rule.UserID).Count(&digestCount), "counting digest")
	assert.Equal(t, digestCount, expected, "digest count mismatch")
}

func TestDo(t *testing.T) {
	t.Run("processes the rule on time", func(t *testing.T) {
		defer testutils.ClearData()

		// Set up
		user := testutils.SetupUserData()
		r1 := database.RepetitionRule{
			Title:      "Rule 1",
			Frequency:  (time.Hour * 24 * 3).Milliseconds(), // three days
			Hour:       12,
			Minute:     2,
			LastActive: 0,
			UserID:     user.ID,
			BookDomain: database.BookDomainAll,
		}

		db := database.DBConn
		testutils.MustExec(t, db.Save(&r1), "preparing rule1")

		c := clock.NewMock()

		// Test
		c.SetNow(time.Date(2009, time.November, 10, 12, 1, 0, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(0))
		assertRepetitionCount(t, r1, 0)

		c.SetNow(time.Date(2009, time.November, 10, 12, 2, 0, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1257854520000))
		assertRepetitionCount(t, r1, 1)

		c.SetNow(time.Date(2009, time.November, 10, 12, 3, 0, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1257854520000))
		assertRepetitionCount(t, r1, 1)

		// 1 day later
		c.SetNow(time.Date(2009, time.November, 11, 12, 2, 0, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1257854520000))
		assertRepetitionCount(t, r1, 1)
		// 2 days later
		c.SetNow(time.Date(2009, time.November, 12, 12, 2, 0, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1257854520000))
		assertRepetitionCount(t, r1, 1)
		// 3 days later - should be processed
		c.SetNow(time.Date(2009, time.November, 13, 12, 2, 0, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1258113720000))
		assertRepetitionCount(t, r1, 2)
		// 4 days later
		c.SetNow(time.Date(2009, time.November, 14, 12, 2, 0, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1258113720000))
		assertRepetitionCount(t, r1, 2)
		// 5 days later
		c.SetNow(time.Date(2009, time.November, 15, 12, 2, 0, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1258113720000))
		assertRepetitionCount(t, r1, 2)
		// 6 days later - should be processed
		c.SetNow(time.Date(2009, time.November, 16, 12, 2, 0, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1258372920000))
		assertRepetitionCount(t, r1, 3)
	})
}

func TestDo_RandomStrategy(t *testing.T) {
	type testData struct {
		User  database.User
		Book1 database.Book
		Book2 database.Book
		Book3 database.Book
		Note1 database.Note
		Note2 database.Note
		Note3 database.Note
	}

	setup := func() testData {
		db := database.DBConn
		user := testutils.SetupUserData()

		b1 := database.Book{
			UserID: user.ID,
			Label:  "js",
		}
		testutils.MustExec(t, db.Save(&b1), "preparing b1")
		b2 := database.Book{
			UserID: user.ID,
			Label:  "css",
		}
		testutils.MustExec(t, db.Save(&b2), "preparing b2")
		b3 := database.Book{
			UserID: user.ID,
			Label:  "golang",
		}
		testutils.MustExec(t, db.Save(&b3), "preparing b3")

		n1 := database.Note{
			UserID:   user.ID,
			BookUUID: b1.UUID,
		}
		testutils.MustExec(t, db.Save(&n1), "preparing n1")
		n2 := database.Note{
			UserID:   user.ID,
			BookUUID: b2.UUID,
		}
		testutils.MustExec(t, db.Save(&n2), "preparing n2")
		n3 := database.Note{
			UserID:   user.ID,
			BookUUID: b3.UUID,
		}
		testutils.MustExec(t, db.Save(&n3), "preparing n3")

		return testData{
			User:  user,
			Book1: b1,
			Book2: b2,
			Book3: b3,
			Note1: n1,
			Note2: n2,
			Note3: n3,
		}
	}

	t.Run("all books", func(t *testing.T) {
		defer testutils.ClearData()

		// Set up
		dat := setup()

		db := database.DBConn
		r1 := database.RepetitionRule{
			Title:      "Rule 1",
			Frequency:  (time.Hour * 24 * 7).Milliseconds(),
			Hour:       21,
			Minute:     0,
			LastActive: 0,
			UserID:     dat.User.ID,
			BookDomain: database.BookDomainAll,
			NoteCount:  5,
		}
		testutils.MustExec(t, db.Save(&r1), "preparing rule1")

		// Execute
		c := clock.NewMock()

		c.SetNow(time.Date(2009, time.November, 10, 21, 0, 0, 0, time.UTC))
		Do(c)

		// Test
		assertLastActive(t, r1.UUID, int64(1257886800000))
		assertRepetitionCount(t, r1, 1)

		var repetition database.Digest
		testutils.MustExec(t, db.Where("rule_id = ? AND user_id = ?", r1.ID, r1.UserID).Preload("Notes").First(&repetition), "finding repetition")

		sort.SliceStable(repetition.Notes, func(i, j int) bool {
			n1 := repetition.Notes[i]
			n2 := repetition.Notes[j]

			return n1.ID < n2.ID
		})

		var n1Record, n2Record, n3Record database.Note
		testutils.MustExec(t, db.Where("uuid = ?", dat.Note1.UUID).First(&n1Record), "finding n1")
		testutils.MustExec(t, db.Where("uuid = ?", dat.Note2.UUID).First(&n2Record), "finding n2")
		testutils.MustExec(t, db.Where("uuid = ?", dat.Note3.UUID).First(&n3Record), "finding n3")
		expected := []database.Note{n1Record, n2Record, n3Record}
		assert.DeepEqual(t, repetition.Notes, expected, "result mismatch")
	})

	t.Run("excluding books", func(t *testing.T) {
		defer testutils.ClearData()

		// Set up
		dat := setup()

		db := database.DBConn
		r1 := database.RepetitionRule{
			Title:      "Rule 1",
			Frequency:  (time.Hour * 24 * 7).Milliseconds(),
			Hour:       21,
			Minute:     0,
			LastActive: 0,
			UserID:     dat.User.ID,
			BookDomain: database.BookDomainExluding,
			Books:      []database.Book{dat.Book1},
			NoteCount:  5,
		}
		testutils.MustExec(t, db.Save(&r1), "preparing rule1")

		// Execute
		c := clock.NewMock()

		c.SetNow(time.Date(2009, time.November, 10, 21, 0, 0, 0, time.UTC))
		Do(c)

		// Test
		assertLastActive(t, r1.UUID, int64(1257886800000))
		assertRepetitionCount(t, r1, 1)

		var repetition database.Digest
		testutils.MustExec(t, db.Where("rule_id = ? AND user_id = ?", r1.ID, r1.UserID).Preload("Notes").First(&repetition), "finding repetition")

		sort.SliceStable(repetition.Notes, func(i, j int) bool {
			n1 := repetition.Notes[i]
			n2 := repetition.Notes[j]

			return n1.ID < n2.ID
		})

		var n2Record, n3Record database.Note
		testutils.MustExec(t, db.Where("uuid = ?", dat.Note2.UUID).First(&n2Record), "finding n2")
		testutils.MustExec(t, db.Where("uuid = ?", dat.Note3.UUID).First(&n3Record), "finding n3")
		expected := []database.Note{n2Record, n3Record}
		assert.DeepEqual(t, repetition.Notes, expected, "result mismatch")
	})

	t.Run("including books", func(t *testing.T) {
		defer testutils.ClearData()

		// Set up
		dat := setup()

		db := database.DBConn
		r1 := database.RepetitionRule{
			Title:      "Rule 1",
			Frequency:  (time.Hour * 24 * 7).Milliseconds(),
			Hour:       21,
			Minute:     0,
			LastActive: 0,
			UserID:     dat.User.ID,
			BookDomain: database.BookDomainIncluding,
			Books:      []database.Book{dat.Book1, dat.Book2},
			NoteCount:  5,
		}
		testutils.MustExec(t, db.Save(&r1), "preparing rule1")

		// Execute
		c := clock.NewMock()

		c.SetNow(time.Date(2009, time.November, 10, 21, 0, 0, 0, time.UTC))
		Do(c)

		// Test
		assertLastActive(t, r1.UUID, int64(1257886800000))
		assertRepetitionCount(t, r1, 1)

		var repetition database.Digest
		testutils.MustExec(t, db.Where("rule_id = ? AND user_id = ?", r1.ID, r1.UserID).Preload("Notes").First(&repetition), "finding repetition")

		sort.SliceStable(repetition.Notes, func(i, j int) bool {
			n1 := repetition.Notes[i]
			n2 := repetition.Notes[j]

			return n1.ID < n2.ID
		})

		var n1Record, n2Record database.Note
		testutils.MustExec(t, db.Where("uuid = ?", dat.Note1.UUID).First(&n1Record), "finding n1")
		testutils.MustExec(t, db.Where("uuid = ?", dat.Note2.UUID).First(&n2Record), "finding n2")
		expected := []database.Note{n1Record, n2Record}
		assert.DeepEqual(t, repetition.Notes, expected, "result mismatch")
	})
}
