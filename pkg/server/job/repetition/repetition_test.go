/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

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

func assertDigestCount(t *testing.T, rule database.RepetitionRule, expected int) {
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
		t0 := time.Date(2009, time.November, 1, 0, 0, 0, 0, time.UTC)
		t1 := time.Date(2009, time.November, 4, 12, 2, 0, 0, time.UTC)
		r1 := database.RepetitionRule{
			Title:      "Rule 1",
			Frequency:  (time.Hour * 24 * 3).Milliseconds(), // three days
			Hour:       12,
			Minute:     2,
			Enabled:    true,
			LastActive: 0,
			NextActive: t1.UnixNano() / int64(time.Millisecond),
			UserID:     user.ID,
			BookDomain: database.BookDomainAll,
			Model: database.Model{
				CreatedAt: t0,
				UpdatedAt: t0,
			},
		}

		db := database.DBConn
		testutils.MustExec(t, db.Save(&r1), "preparing rule1")

		c := clock.NewMock()

		// Test
		// 1 day later
		c.SetNow(time.Date(2009, time.November, 2, 12, 2, 1, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(0))
		assertDigestCount(t, r1, 0)

		// 2 days later
		c.SetNow(time.Date(2009, time.November, 3, 12, 2, 1, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(0))
		assertDigestCount(t, r1, 0)

		// 3 days later - should be processed
		c.SetNow(time.Date(2009, time.November, 4, 12, 1, 1, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(0))
		assertDigestCount(t, r1, 0)

		c.SetNow(time.Date(2009, time.November, 4, 12, 2, 1, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1257336120000))
		assertDigestCount(t, r1, 1)

		c.SetNow(time.Date(2009, time.November, 4, 12, 3, 1, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1257336120000))
		assertDigestCount(t, r1, 1)

		// 4 day later
		c.SetNow(time.Date(2009, time.November, 5, 12, 2, 1, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1257336120000))
		assertDigestCount(t, r1, 1)
		// 5 days later
		c.SetNow(time.Date(2009, time.November, 6, 12, 2, 1, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1257336120000))
		assertDigestCount(t, r1, 1)
		// 6 days later - should be processed
		c.SetNow(time.Date(2009, time.November, 7, 12, 2, 1, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1257595320000))
		assertDigestCount(t, r1, 2)
		// 7 days later
		c.SetNow(time.Date(2009, time.November, 8, 12, 2, 1, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1257595320000))
		assertDigestCount(t, r1, 2)
		// 8 days later
		c.SetNow(time.Date(2009, time.November, 9, 12, 2, 1, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1257595320000))
		assertDigestCount(t, r1, 2)
		// 9 days later - should be processed
		c.SetNow(time.Date(2009, time.November, 10, 12, 2, 1, 0, time.UTC))
		Do(c)
		assertLastActive(t, r1.UUID, int64(1257854520000))
		assertDigestCount(t, r1, 3)
	})

	t.Run("recovers correct next_active value if missed processing in the past", func(t *testing.T) {
		defer testutils.ClearData()

		// Set up
		user := testutils.SetupUserData()
		t0 := time.Date(2009, time.November, 1, 12, 2, 0, 0, time.UTC)
		t1 := time.Date(2009, time.November, 4, 12, 2, 0, 0, time.UTC)
		r1 := database.RepetitionRule{
			Title:      "Rule 1",
			Frequency:  (time.Hour * 24 * 3).Milliseconds(), // three days
			Hour:       12,
			Minute:     2,
			Enabled:    true,
			LastActive: t0.UnixNano() / int64(time.Millisecond),
			NextActive: t1.UnixNano() / int64(time.Millisecond),
			UserID:     user.ID,
			BookDomain: database.BookDomainAll,
			Model: database.Model{
				CreatedAt: t0,
				UpdatedAt: t0,
			},
		}

		db := database.DBConn
		testutils.MustExec(t, db.Save(&r1), "preparing rule1")

		c := clock.NewMock()
		c.SetNow(time.Date(2009, time.November, 10, 12, 2, 1, 0, time.UTC))
		Do(c)

		var rule database.RepetitionRule
		testutils.MustExec(t, db.Where("uuid = ?", r1.UUID).First(&rule), "finding rule1")

		assert.Equal(t, rule.LastActive, time.Date(2009, time.November, 10, 12, 2, 0, 0, time.UTC).UnixNano()/int64(time.Millisecond), "LastActive mismsatch")
		assert.Equal(t, rule.NextActive, time.Date(2009, time.November, 13, 12, 2, 0, 0, time.UTC).UnixNano()/int64(time.Millisecond), "NextActive mismsatch")
		assertDigestCount(t, r1, 1)
	})
}

func TestDo_Disabled(t *testing.T) {
	defer testutils.ClearData()

	// Set up
	user := testutils.SetupUserData()
	t0 := time.Date(2009, time.November, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2009, time.November, 4, 12, 2, 0, 0, time.UTC)
	r1 := database.RepetitionRule{
		Title:      "Rule 1",
		Frequency:  (time.Hour * 24 * 3).Milliseconds(), // three days
		Hour:       12,
		Minute:     2,
		LastActive: 0,
		NextActive: t1.UnixNano() / int64(time.Millisecond),
		UserID:     user.ID,
		Enabled:    false,
		BookDomain: database.BookDomainAll,
		Model: database.Model{
			CreatedAt: t0,
			UpdatedAt: t0,
		},
	}

	db := database.DBConn
	testutils.MustExec(t, db.Save(&r1), "preparing rule1")

	// Execute
	c := clock.NewMock()
	c.SetNow(time.Date(2009, time.November, 4, 12, 2, 0, 0, time.UTC))
	Do(c)

	// Test
	assertLastActive(t, r1.UUID, int64(0))
	assertDigestCount(t, r1, 0)
}

func TestDo_BalancedStrategy(t *testing.T) {
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
		t0 := time.Date(2009, time.November, 1, 12, 0, 0, 0, time.UTC)
		t1 := time.Date(2009, time.November, 8, 21, 0, 0, 0, time.UTC)
		r1 := database.RepetitionRule{
			Title:      "Rule 1",
			Frequency:  (time.Hour * 24 * 7).Milliseconds(),
			Hour:       21,
			Minute:     0,
			LastActive: 0,
			NextActive: t1.UnixNano() / int64(time.Millisecond),
			Enabled:    true,
			UserID:     dat.User.ID,
			BookDomain: database.BookDomainAll,
			NoteCount:  5,
			Model: database.Model{
				CreatedAt: t0,
				UpdatedAt: t0,
			},
		}
		testutils.MustExec(t, db.Save(&r1), "preparing rule1")

		// Execute
		c := clock.NewMock()

		c.SetNow(time.Date(2009, time.November, 8, 21, 0, 0, 0, time.UTC))
		Do(c)

		// Test
		assertLastActive(t, r1.UUID, int64(1257714000000))
		assertDigestCount(t, r1, 1)

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
		t0 := time.Date(2009, time.November, 1, 12, 0, 0, 0, time.UTC)
		t1 := time.Date(2009, time.November, 8, 21, 0, 0, 0, time.UTC)
		r1 := database.RepetitionRule{
			Title:      "Rule 1",
			Frequency:  (time.Hour * 24 * 7).Milliseconds(),
			Hour:       21,
			Enabled:    true,
			Minute:     0,
			LastActive: 0,
			NextActive: t1.UnixNano() / int64(time.Millisecond),
			UserID:     dat.User.ID,
			BookDomain: database.BookDomainExluding,
			Books:      []database.Book{dat.Book1},
			NoteCount:  5,
			Model: database.Model{
				CreatedAt: t0,
				UpdatedAt: t0,
			},
		}
		testutils.MustExec(t, db.Save(&r1), "preparing rule1")

		// Execute
		c := clock.NewMock()

		c.SetNow(time.Date(2009, time.November, 8, 21, 0, 1, 0, time.UTC))
		Do(c)

		// Test
		assertLastActive(t, r1.UUID, int64(1257714000000))
		assertDigestCount(t, r1, 1)

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
		t0 := time.Date(2009, time.November, 1, 12, 0, 0, 0, time.UTC)
		t1 := time.Date(2009, time.November, 8, 21, 0, 0, 0, time.UTC)
		r1 := database.RepetitionRule{
			Title:      "Rule 1",
			Frequency:  (time.Hour * 24 * 7).Milliseconds(),
			Hour:       21,
			Enabled:    true,
			Minute:     0,
			LastActive: 0,
			NextActive: t1.UnixNano() / int64(time.Millisecond),
			UserID:     dat.User.ID,
			BookDomain: database.BookDomainIncluding,
			Books:      []database.Book{dat.Book1, dat.Book2},
			NoteCount:  5,
			Model: database.Model{
				CreatedAt: t0,
				UpdatedAt: t0,
			},
		}
		testutils.MustExec(t, db.Save(&r1), "preparing rule1")

		// Execute
		c := clock.NewMock()

		c.SetNow(time.Date(2009, time.November, 8, 21, 0, 0, 0, time.UTC))
		Do(c)

		// Test
		assertLastActive(t, r1.UUID, int64(1257714000000))
		assertDigestCount(t, r1, 1)

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
