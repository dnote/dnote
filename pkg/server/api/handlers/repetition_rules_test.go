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

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/api/presenters"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func init() {
	testutils.InitTestDB()
}

func TestGetRepetitionRule(t *testing.T) {
	defer testutils.ClearData()
	db := database.DBConn

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	user := testutils.SetupUserData()

	b1 := database.Book{
		USN:   11,
		Label: "js",
	}
	testutils.MustExec(t, db.Save(&b1), "preparing book1")

	r1 := database.RepetitionRule{
		Title:      "Rule 1",
		Frequency:  (time.Hour * 24 * 7).Milliseconds(),
		Hour:       21,
		Minute:     0,
		LastActive: 0,
		UserID:     user.ID,
		BookDomain: database.BookDomainExluding,
		Books:      []database.Book{b1},
		NoteCount:  5,
	}
	testutils.MustExec(t, db.Save(&r1), "preparing rule1")

	// Execute
	req := testutils.MakeReq(server, "GET", fmt.Sprintf("/repetition_rules/%s", r1.UUID), "")
	res := testutils.HTTPAuthDo(t, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var payload presenters.RepetitionRule
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload"))
	}

	var r1Record database.RepetitionRule
	testutils.MustExec(t, db.Where("uuid = ?", r1.UUID).First(&r1Record), "finding r1Record")
	var b1Record database.Book
	testutils.MustExec(t, db.Where("uuid = ?", b1.UUID).First(&b1Record), "finding b1Record")

	expected := presenters.RepetitionRule{
		UUID:       r1Record.UUID,
		Title:      r1Record.Title,
		Enabled:    r1Record.Enabled,
		Hour:       r1Record.Hour,
		Minute:     r1Record.Minute,
		Frequency:  r1Record.Frequency,
		BookDomain: r1Record.BookDomain,
		NoteCount:  r1Record.NoteCount,
		LastActive: r1Record.LastActive,
		Books: []presenters.Book{
			{
				UUID:      b1Record.UUID,
				USN:       b1Record.USN,
				Label:     b1Record.Label,
				CreatedAt: presenters.FormatTS(b1Record.CreatedAt),
				UpdatedAt: presenters.FormatTS(b1Record.UpdatedAt),
			},
		},
		CreatedAt: presenters.FormatTS(r1Record.CreatedAt),
		UpdatedAt: presenters.FormatTS(r1Record.UpdatedAt),
	}

	assert.DeepEqual(t, payload, expected, "payload mismatch")
}

func TestGetRepetitionRules(t *testing.T) {
	defer testutils.ClearData()
	db := database.DBConn

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	user := testutils.SetupUserData()

	b1 := database.Book{
		USN:   11,
		Label: "js",
	}
	testutils.MustExec(t, db.Save(&b1), "preparing book1")

	r1 := database.RepetitionRule{
		Title:      "Rule 1",
		Frequency:  (time.Hour * 24 * 7).Milliseconds(),
		Hour:       21,
		Minute:     0,
		LastActive: 0,
		UserID:     user.ID,
		BookDomain: database.BookDomainExluding,
		Books:      []database.Book{b1},
		NoteCount:  5,
	}
	testutils.MustExec(t, db.Save(&r1), "preparing rule1")

	r2 := database.RepetitionRule{
		Title:      "Rule 2",
		Frequency:  (time.Hour * 24 * 7 * 2).Milliseconds(),
		Hour:       2,
		Minute:     0,
		LastActive: 0,
		UserID:     user.ID,
		BookDomain: database.BookDomainExluding,
		Books:      []database.Book{},
		NoteCount:  5,
	}
	testutils.MustExec(t, db.Save(&r2), "preparing rule2")

	// Execute
	req := testutils.MakeReq(server, "GET", "/repetition_rules", "")
	res := testutils.HTTPAuthDo(t, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var payload []presenters.RepetitionRule
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatal(errors.Wrap(err, "decoding payload"))
	}

	var r1Record, r2Record database.RepetitionRule
	testutils.MustExec(t, db.Where("uuid = ?", r1.UUID).First(&r1Record), "finding r1Record")
	testutils.MustExec(t, db.Where("uuid = ?", r2.UUID).First(&r2Record), "finding r2Record")
	var b1Record database.Book
	testutils.MustExec(t, db.Where("uuid = ?", b1.UUID).First(&b1Record), "finding b1Record")

	expected := []presenters.RepetitionRule{
		{
			UUID:       r1Record.UUID,
			Title:      r1Record.Title,
			Enabled:    r1Record.Enabled,
			Hour:       r1Record.Hour,
			Minute:     r1Record.Minute,
			Frequency:  r1Record.Frequency,
			BookDomain: r1Record.BookDomain,
			NoteCount:  r1Record.NoteCount,
			LastActive: r1Record.LastActive,
			Books: []presenters.Book{
				{
					UUID:      b1Record.UUID,
					USN:       b1Record.USN,
					Label:     b1Record.Label,
					CreatedAt: presenters.FormatTS(b1Record.CreatedAt),
					UpdatedAt: presenters.FormatTS(b1Record.UpdatedAt),
				},
			},
			CreatedAt: presenters.FormatTS(r1Record.CreatedAt),
			UpdatedAt: presenters.FormatTS(r1Record.UpdatedAt),
		},
		{
			UUID:       r2Record.UUID,
			Title:      r2Record.Title,
			Enabled:    r2Record.Enabled,
			Hour:       r2Record.Hour,
			Minute:     r2Record.Minute,
			Frequency:  r2Record.Frequency,
			BookDomain: r2Record.BookDomain,
			NoteCount:  r2Record.NoteCount,
			LastActive: r2Record.LastActive,
			Books:      []presenters.Book{},
			CreatedAt:  presenters.FormatTS(r2Record.CreatedAt),
			UpdatedAt:  presenters.FormatTS(r2Record.UpdatedAt),
		},
	}

	assert.DeepEqual(t, payload, expected, "payload mismatch")
}

func TestCreateRepetitionRules(t *testing.T) {
	t.Run("all books", func(t *testing.T) {
		defer testutils.ClearData()
		db := database.DBConn

		// Setup
		server := httptest.NewServer(NewRouter(&App{
			Clock: clock.NewMock(),
		}))
		defer server.Close()

		user := testutils.SetupUserData()

		// Execute
		dat := `{
	"title": "Rule 1",
	"enabled": true,
	"hour": 8,
	"minute": 30,
	"frequency": 6048000000,
	"book_domain": "all",
	"book_uuids": [],
	"note_count": 20
}`
		req := testutils.MakeReq(server, "POST", "/repetition_rules", dat)
		res := testutils.HTTPAuthDo(t, req, user)

		// Test
		assert.StatusCodeEquals(t, res, http.StatusCreated, "")

		var ruleCount int
		testutils.MustExec(t, db.Model(&database.RepetitionRule{}).Count(&ruleCount), "counting rules")
		assert.Equalf(t, ruleCount, 1, "reperition rule count mismatch")

		var rule database.RepetitionRule
		testutils.MustExec(t, db.Preload("Books").First(&rule), "finding b1Record")

		assert.NotEqual(t, rule.UUID, "", "rule UUID mismatch")
		assert.Equal(t, rule.Title, "Rule 1", "rule Title mismatch")
		assert.Equal(t, rule.Enabled, true, "rule Enabled mismatch")
		assert.Equal(t, rule.Hour, 8, "rule HourTitle mismatch")
		assert.Equal(t, rule.Minute, 30, "rule Minute mismatch")
		assert.Equal(t, rule.Frequency, int64(6048000000), "rule Frequency mismatch")
		assert.Equal(t, rule.BookDomain, "all", "rule BookDomain mismatch")
		assert.DeepEqual(t, rule.Books, []database.Book{}, "rule Books mismatch")
		assert.Equal(t, rule.NoteCount, 20, "rule NoteCount mismatch")
	})

	bookDomainTestCases := []string{
		"including",
		"excluding",
	}
	for _, tc := range bookDomainTestCases {
		t.Run(tc, func(t *testing.T) {
			defer testutils.ClearData()
			db := database.DBConn

			// Setup
			server := httptest.NewServer(NewRouter(&App{
				Clock: clock.NewMock(),
			}))
			defer server.Close()

			user := testutils.SetupUserData()

			b1 := database.Book{
				UserID: user.ID,
				Label:  "css",
			}
			testutils.MustExec(t, db.Save(&b1), "preparing b1")

			// Execute
			dat := fmt.Sprintf(`{
	"title": "Rule 1",
	"enabled": true,
	"hour": 8,
	"minute": 30,
	"frequency": 6048000000,
	"book_domain": "%s",
	"book_uuids": ["%s"],
	"note_count": 20
}`, tc, b1.UUID)
			req := testutils.MakeReq(server, "POST", "/repetition_rules", dat)
			res := testutils.HTTPAuthDo(t, req, user)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusCreated, "")

			var ruleCount int
			testutils.MustExec(t, db.Model(&database.RepetitionRule{}).Count(&ruleCount), "counting rules")
			assert.Equalf(t, ruleCount, 1, "reperition rule count mismatch")

			var rule database.RepetitionRule
			testutils.MustExec(t, db.Preload("Books").First(&rule), "finding b1Record")

			var b1Record database.Book
			testutils.MustExec(t, db.Where("id = ?", b1.ID).First(&b1Record), "finding b1Record")

			assert.NotEqual(t, rule.UUID, "", "rule UUID mismatch")
			assert.Equal(t, rule.Title, "Rule 1", "rule Title mismatch")
			assert.Equal(t, rule.Enabled, true, "rule Enabled mismatch")
			assert.Equal(t, rule.Hour, 8, "rule HourTitle mismatch")
			assert.Equal(t, rule.Minute, 30, "rule Minute mismatch")
			assert.Equal(t, rule.Frequency, int64(6048000000), "rule Frequency mismatch")
			assert.Equal(t, rule.BookDomain, tc, "rule BookDomain mismatch")
			assert.DeepEqual(t, rule.Books, []database.Book{b1Record}, "rule Books mismatch")
			assert.Equal(t, rule.NoteCount, 20, "rule NoteCount mismatch")
		})
	}
}

func TestUpdateRepetitionRules(t *testing.T) {
	defer testutils.ClearData()
	db := database.DBConn

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	user := testutils.SetupUserData()

	// Execute
	r1 := database.RepetitionRule{
		Title:      "Rule 1",
		UserID:     user.ID,
		Enabled:    false,
		Hour:       8,
		Minute:     30,
		Frequency:  6048000000,
		BookDomain: "all",
		Books:      []database.Book{},
		NoteCount:  20,
	}
	testutils.MustExec(t, db.Save(&r1), "preparing r1")
	b1 := database.Book{
		UserID: user.ID,
		USN:    11,
		Label:  "js",
	}
	testutils.MustExec(t, db.Save(&b1), "preparing book1")

	dat := fmt.Sprintf(`{
	"title": "Rule 1 - edited",
	"enabled": true,
	"hour": 18,
	"minute": 40,
	"frequency": 259200000,
	"book_domain": "including",
	"book_uuids": ["%s"],
	"note_count": 30
}`, b1.UUID)
	endpoint := fmt.Sprintf("/repetition_rules/%s", r1.UUID)
	req := testutils.MakeReq(server, "PATCH", endpoint, dat)
	res := testutils.HTTPAuthDo(t, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var totalRuleCount int
	testutils.MustExec(t, db.Model(&database.RepetitionRule{}).Count(&totalRuleCount), "counting rules")
	assert.Equalf(t, totalRuleCount, 1, "reperition rule count mismatch")

	var rule database.RepetitionRule
	testutils.MustExec(t, db.Preload("Books").First(&rule), "finding b1Record")

	var b1Record database.Book
	testutils.MustExec(t, db.Where("id = ?", b1.ID).First(&b1Record), "finding b1Record")

	assert.NotEqual(t, rule.UUID, "", "rule UUID mismatch")
	assert.Equal(t, rule.Title, "Rule 1 - edited", "rule Title mismatch")
	assert.Equal(t, rule.Enabled, true, "rule Enabled mismatch")
	assert.Equal(t, rule.Hour, 18, "rule HourTitle mismatch")
	assert.Equal(t, rule.Minute, 40, "rule Minute mismatch")
	assert.Equal(t, rule.Frequency, int64(259200000), "rule Frequency mismatch")
	assert.Equal(t, rule.BookDomain, "including", "rule BookDomain mismatch")
	assert.DeepEqual(t, rule.Books, []database.Book{b1Record}, "rule Books mismatch")
	assert.Equal(t, rule.NoteCount, 30, "rule NoteCount mismatch")
}

func TestDeleteRepetitionRules(t *testing.T) {
	defer testutils.ClearData()
	db := database.DBConn

	// Setup
	server := httptest.NewServer(NewRouter(&App{
		Clock: clock.NewMock(),
	}))
	defer server.Close()

	user := testutils.SetupUserData()

	// Execute
	r1 := database.RepetitionRule{
		Title:      "Rule 1",
		UserID:     user.ID,
		Enabled:    true,
		Hour:       8,
		Minute:     30,
		Frequency:  6048000000,
		BookDomain: "all",
		Books:      []database.Book{},
		NoteCount:  20,
	}
	testutils.MustExec(t, db.Save(&r1), "preparing r1")

	r2 := database.RepetitionRule{
		Title:      "Rule 1",
		UserID:     user.ID,
		Enabled:    true,
		Hour:       8,
		Minute:     30,
		Frequency:  6048000000,
		BookDomain: "all",
		Books:      []database.Book{},
		NoteCount:  20,
	}
	testutils.MustExec(t, db.Save(&r2), "preparing r2")

	endpoint := fmt.Sprintf("/repetition_rules/%s", r1.UUID)
	req := testutils.MakeReq(server, "DELETE", endpoint, "")
	res := testutils.HTTPAuthDo(t, req, user)

	// Test
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var totalRuleCount int
	testutils.MustExec(t, db.Model(&database.RepetitionRule{}).Count(&totalRuleCount), "counting rules")
	assert.Equalf(t, totalRuleCount, 1, "reperition rule count mismatch")

	var r2Count int
	testutils.MustExec(t, db.Model(&database.RepetitionRule{}).Where("id = ?", r2.ID).Count(&r2Count), "counting r2")
	assert.Equalf(t, r2Count, 1, "r2 count mismatch")
}

func TestCreateUpdateRepetitionRules_BadRequest(t *testing.T) {
	testCases := []string{
		// empty title
		`{
			"title": "",
			"enabled": true,
			"hour": 8,
			"minute": 30,
			"frequency": 6048000000,
			"book_domain": "all",
			"book_uuids": [],
			"note_count": 20
		}`,
		// empty frequency
		`{
			"title": "Rule 1",
			"enabled": true,
			"hour": 8,
			"minute": 30,
			"frequency": 0,
			"book_domain": "some_invalid_book_domain",
			"book_uuids": [],
			"note_count": 20
		}`,
		// empty note count
		`{
			"title": "Rule 1",
			"enabled": true,
			"hour": 8,
			"minute": 30,
			"frequency": 6048000000,
			"book_domain": "all",
			"book_uuids": [],
			"note_count": 0
		}`,
		// invalid book doamin
		`{
			"title": "Rule 1",
			"enabled": true,
			"hour": 8,
			"minute": 30,
			"frequency": 6048000000,
			"book_domain": "some_invalid_book_domain",
			"book_uuids": [],
			"note_count": 20
		}`,
		// invalid combination of book domain and book_uuids
		`{
			"title": "Rule 1",
			"enabled": true,
			"hour": 8,
			"minute": 30,
			"frequency": 6048000000,
			"book_domain": "excluding",
			"book_uuids": [],
			"note_count": 20
		}`,
		`{
			"title": "Rule 1",
			"enabled": true,
			"hour": 8,
			"minute": 30,
			"frequency": 6048000000,
			"book_domain": "including",
			"book_uuids": [],
			"note_count": 20
		}`,
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case - create %d", idx), func(t *testing.T) {
			defer testutils.ClearData()
			db := database.DBConn

			// Setup
			server := httptest.NewServer(NewRouter(&App{
				Clock: clock.NewMock(),
			}))
			defer server.Close()

			user := testutils.SetupUserData()

			// Execute
			req := testutils.MakeReq(server, "POST", "/repetition_rules", tc)
			res := testutils.HTTPAuthDo(t, req, user)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusBadRequest, "")

			var ruleCount int
			testutils.MustExec(t, db.Model(&database.RepetitionRule{}).Count(&ruleCount), "counting rules")
			assert.Equalf(t, ruleCount, 0, "reperition rule count mismatch")
		})

		t.Run(fmt.Sprintf("test case %d - update", idx), func(t *testing.T) {
			defer testutils.ClearData()
			db := database.DBConn

			// Setup
			user := testutils.SetupUserData()
			r1 := database.RepetitionRule{
				Title:      "Rule 1",
				UserID:     user.ID,
				Enabled:    false,
				Hour:       8,
				Minute:     30,
				Frequency:  6048000000,
				BookDomain: "all",
				Books:      []database.Book{},
				NoteCount:  20,
			}
			testutils.MustExec(t, db.Save(&r1), "preparing r1")
			b1 := database.Book{
				UserID: user.ID,
				USN:    11,
				Label:  "js",
			}
			testutils.MustExec(t, db.Save(&b1), "preparing book1")

			server := httptest.NewServer(NewRouter(&App{
				Clock: clock.NewMock(),
			}))
			defer server.Close()

			// Execute
			req := testutils.MakeReq(server, "PATCH", fmt.Sprintf("/repetition_rules/%s", r1.UUID), tc)
			res := testutils.HTTPAuthDo(t, req, user)

			// Test
			assert.StatusCodeEquals(t, res, http.StatusBadRequest, "")

			var ruleCount int
			testutils.MustExec(t, db.Model(&database.RepetitionRule{}).Count(&ruleCount), "counting rules")
			assert.Equalf(t, ruleCount, 1, "reperition rule count mismatch")
		})
	}
}
