/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
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

package tmpl

import (
	"html/template"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func TestDefaultPageGetData(t *testing.T) {
	p := defaultPage{}

	result := p.getData()

	assert.Equal(t, result.MetaTags, template.HTML(""), "MetaTags mismatch")
	assert.Equal(t, result.Title, "Dnote", "Title mismatch")
}

func TestNotePageGetData(t *testing.T) {
	// Set time.Local to UTC for deterministic test
	time.Local = time.UTC

	db := testutils.InitMemoryDB(t)
	a, err := NewAppShell(db, nil)
	if err != nil {
		t.Fatal(errors.Wrap(err, "preparing app shell"))
	}

	p := notePage{
		Note: database.Note{
			Book: database.Book{
				Label: "vocabulary",
			},
			AddedOn: time.Date(2019, time.January, 2, 0, 0, 0, 0, time.UTC).UnixNano(),
		},
		T: a.T,
	}

	result, err := p.getData()
	if err != nil {
		t.Fatal(errors.Wrap(err, "executing"))
	}

	assert.NotEqual(t, result.MetaTags, template.HTML(""), "MetaTags should not be empty")
	assert.Equal(t, result.Title, "Note: vocabulary (Jan 2 2019)", "Title mismatch")
}
