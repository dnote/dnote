/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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
	"fmt"
	"net/http"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestCreateNoteReview(t *testing.T) {
	defer testutils.ClearData()

	// Setup
	server := MustNewServer(t, &app.App{
		Clock: clock.NewMock(),
	})
	defer server.Close()

	user := testutils.SetupUserData()
	b1 := database.Book{
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")
	n1 := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
	}
	testutils.MustExec(t, testutils.DB.Save(&n1), "preparing n1")
	d1 := database.Digest{
		UserID: user.ID,
	}
	testutils.MustExec(t, testutils.DB.Save(&d1), "preparing d1")

	// multiple requests should create at most one receipt
	for i := 0; i < 3; i++ {
		dat := fmt.Sprintf(`{"note_uuid": "%s", "digest_uuid": "%s"}`, n1.UUID, d1.UUID)
		req := testutils.MakeReq(server.URL, http.MethodPost, "/note_review", dat)
		res := testutils.HTTPAuthDo(t, req, user)
		assert.StatusCodeEquals(t, res, http.StatusOK, "")

		var noteReviewCount int
		testutils.MustExec(t, testutils.DB.Model(&database.NoteReview{}).Count(&noteReviewCount), "counting note_reviews")
		assert.Equalf(t, noteReviewCount, 1, "counting note_review")

		var noteReviewRecord database.NoteReview
		testutils.MustExec(t, testutils.DB.Where("user_id = ? AND note_id = ? AND digest_id = ?", user.ID, n1.ID, d1.ID).First(&noteReviewRecord), "finding note_review record")
	}
}

func TestDeleteNoteReview(t *testing.T) {
	defer testutils.ClearData()

	// Setup
	server := MustNewServer(t, &app.App{
		Clock: clock.NewMock(),
	})
	defer server.Close()

	user := testutils.SetupUserData()
	b1 := database.Book{
		UserID: user.ID,
		Label:  "js",
	}
	testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")
	n1 := database.Note{
		UserID:   user.ID,
		BookUUID: b1.UUID,
	}
	testutils.MustExec(t, testutils.DB.Save(&n1), "preparing n1")
	d1 := database.Digest{
		UserID: user.ID,
	}
	testutils.MustExec(t, testutils.DB.Save(&d1), "preparing d1")
	nr1 := database.NoteReview{
		UserID:   user.ID,
		NoteID:   n1.ID,
		DigestID: d1.ID,
	}
	testutils.MustExec(t, testutils.DB.Save(&nr1), "preparing nr1")

	dat := fmt.Sprintf(`{"note_uuid": "%s", "digest_uuid": "%s"}`, n1.UUID, d1.UUID)
	req := testutils.MakeReq(server.URL, http.MethodDelete, "/note_review", dat)
	res := testutils.HTTPAuthDo(t, req, user)
	assert.StatusCodeEquals(t, res, http.StatusOK, "")

	var noteReviewCount int
	testutils.MustExec(t, testutils.DB.Model(&database.NoteReview{}).Count(&noteReviewCount), "counting note_reviews")
	assert.Equal(t, noteReviewCount, 0, "counting note_review")
}
