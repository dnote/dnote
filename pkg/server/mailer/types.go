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

package mailer

import (
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/justincampbell/timeago"
)

// DigestNoteInfo contains note information for digest emails
type DigestNoteInfo struct {
	UUID      string
	Content   string
	BookLabel string
	TimeAgo   string
	Stage     int
}

// DigestTmplData is a template data for digest emails
type DigestTmplData struct {
	Subject           string
	NoteInfo          []DigestNoteInfo
	ActiveBookCount   int
	ActiveNoteCount   int
	EmailSessionToken string
}

// NewNoteInfo returns a new NoteInfo
func NewNoteInfo(note database.Note, stage int) DigestNoteInfo {
	tm := time.Unix(0, int64(note.AddedOn))

	return DigestNoteInfo{
		UUID:      note.UUID,
		Content:   note.Body,
		BookLabel: note.Book.Label,
		TimeAgo:   timeago.FromTime(tm),
		Stage:     stage,
	}
}
