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
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/middleware"
	"github.com/dnote/dnote/pkg/server/operations"
	"github.com/pkg/errors"
)

var newlineRegexp = regexp.MustCompile(`\r?\n`)

// tmplData is the data to be passed to the app shell template
type tmplData struct {
	Title    string
	MetaTags template.HTML
}

type noteMetaTagsData struct {
	Title       string
	Description string
}

type notePage struct {
	Note database.Note
	T    *template.Template
}

func (a AppShell) newNotePage(r *http.Request, noteUUID string) (notePage, error) {
	user, _, err := middleware.AuthWithSession(a.DB, r)
	if err != nil {
		return notePage{}, errors.Wrap(err, "authenticating with session")
	}

	note, ok, err := operations.GetNote(a.DB, noteUUID, &user)

	if !ok {
		return notePage{}, ErrNotFound
	}
	if err != nil {
		return notePage{}, errors.Wrap(err, "getting note")
	}

	return notePage{note, a.T}, nil
}

func (p notePage) getTitle() string {
	note := p.Note
	date := time.Unix(0, note.AddedOn).Format("Jan 2 2006")

	return fmt.Sprintf("Note: %s (%s)", note.Book.Label, date)
}

func excerpt(s string, maxLen int) string {
	if len(s) > maxLen {

		var lastIdx int
		if maxLen > 3 {
			lastIdx = maxLen - 3
		} else {
			lastIdx = maxLen
		}

		return s[:lastIdx] + "..."
	}

	return s
}

func formatMetaDescContent(s string) string {
	desc := excerpt(s, 200)
	desc = strings.Trim(desc, " ")

	return newlineRegexp.ReplaceAllString(desc, " ")
}

func (p notePage) getMetaTags() (template.HTML, error) {
	title := p.getTitle()
	desc := formatMetaDescContent(p.Note.Body)

	data := noteMetaTagsData{
		Title:       title,
		Description: desc,
	}

	var buf bytes.Buffer
	if err := p.T.ExecuteTemplate(&buf, templateNoteMetaTags, data); err != nil {
		return "", errors.Wrap(err, "executing template")
	}

	return template.HTML(buf.String()), nil
}

func (p notePage) getData() (tmplData, error) {
	mt, err := p.getMetaTags()
	if err != nil {
		return tmplData{}, errors.Wrap(err, "getting meta tags")
	}

	dat := tmplData{
		Title:    p.getTitle(),
		MetaTags: mt,
	}

	return dat, nil
}

type defaultPage struct {
}

func (p defaultPage) getData() tmplData {
	return tmplData{
		Title:    "Dnote",
		MetaTags: "",
	}
}
