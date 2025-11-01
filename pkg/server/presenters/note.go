/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package presenters

import (
	"time"

	"github.com/dnote/dnote/pkg/server/database"
)

// Note is a result of PresentNote
type Note struct {
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"content"`
	AddedOn   int64     `json:"added_on"`
	USN       int       `json:"usn"`
	Book      NoteBook  `json:"book"`
	User      NoteUser  `json:"user"`
}

// NoteBook is a nested book for PresentNotesResult
type NoteBook struct {
	UUID  string `json:"uuid"`
	Label string `json:"label"`
}

// NoteUser is a nested book for PresentNotesResult
type NoteUser struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

// PresentNote presents note
func PresentNote(note database.Note) Note {
	ret := Note{
		UUID:      note.UUID,
		CreatedAt: FormatTS(note.CreatedAt),
		UpdatedAt: FormatTS(note.UpdatedAt),
		Body:      note.Body,
		AddedOn:   note.AddedOn,
		USN:       note.USN,
		Book: NoteBook{
			UUID:  note.Book.UUID,
			Label: note.Book.Label,
		},
		User: NoteUser{
			UUID: note.User.UUID,
		},
	}

	return ret
}

// PresentNotes presents notes
func PresentNotes(notes []database.Note) []Note {
	ret := []Note{}

	for _, note := range notes {
		p := PresentNote(note)
		ret = append(ret, p)
	}

	return ret
}
