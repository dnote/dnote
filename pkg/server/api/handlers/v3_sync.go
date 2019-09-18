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
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/dnote/dnote/pkg/server/api/helpers"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/pkg/errors"
)

// fullSyncBefore is the system-wide timestamp that represents the point in time
// before which clients must perform a full-sync rather than incremental sync.
const fullSyncBefore = 0

// SyncFragment contains a piece of information about the server's state.
// It is used to transfer the server's state to the client gradually without having to
// transfer the whole state at once.
type SyncFragment struct {
	FragMaxUSN    int            `json:"frag_max_usn"`
	UserMaxUSN    int            `json:"user_max_usn"`
	CurrentTime   int64          `json:"current_time"`
	Notes         []SyncFragNote `json:"notes"`
	Books         []SyncFragBook `json:"books"`
	ExpungedNotes []string       `json:"expunged_notes"`
	ExpungedBooks []string       `json:"expunged_books"`
}

// SyncFragNote represents a note in a sync fragment and contains only the necessary information
// for the client to sync the note locally
type SyncFragNote struct {
	UUID      string    `json:"uuid"`
	BookUUID  string    `json:"book_uuid"`
	USN       int       `json:"usn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AddedOn   int64     `json:"added_on"`
	EditedOn  int64     `json:"edited_on"`
	Body      string    `json:"content"`
	Public    bool      `json:"public"`
	Deleted   bool      `json:"deleted"`
}

// NewFragNote presents the given note as a SyncFragNote
func NewFragNote(note database.Note) SyncFragNote {
	return SyncFragNote{
		UUID:      note.UUID,
		USN:       note.USN,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
		AddedOn:   note.AddedOn,
		EditedOn:  note.EditedOn,
		Body:      note.Body,
		Public:    note.Public,
		Deleted:   note.Deleted,
		BookUUID:  note.BookUUID,
	}
}

// SyncFragBook represents a book in a sync fragment and contains only the necessary information
// for the client to sync the note locally
type SyncFragBook struct {
	UUID      string    `json:"uuid"`
	USN       int       `json:"usn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AddedOn   int64     `json:"added_on"`
	Label     string    `json:"label"`
	Deleted   bool      `json:"deleted"`
}

// NewFragBook presents the given book as a SyncFragBook
func NewFragBook(book database.Book) SyncFragBook {
	return SyncFragBook{
		UUID:      book.UUID,
		USN:       book.USN,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
		AddedOn:   book.AddedOn,
		Label:     book.Label,
		Deleted:   book.Deleted,
	}
}

type usnItem struct {
	usn int
	val interface{}
}

type queryParamError struct {
	key     string
	value   string
	message string
}

func (e *queryParamError) Error() string {
	return fmt.Sprintf("invalid query param %s=%s. %s", e.key, e.value, e.message)
}

func (a *App) newFragment(userID, userMaxUSN, afterUSN, limit int) (SyncFragment, error) {
	db := database.DBConn

	var notes []database.Note
	if err := db.Where("user_id = ? AND usn > ? AND usn <= ?", userID, afterUSN, userMaxUSN).Order("usn ASC").Limit(limit).Find(&notes).Error; err != nil {
		return SyncFragment{}, nil
	}
	var books []database.Book
	if err := db.Where("user_id = ? AND usn > ? AND usn <= ?", userID, afterUSN, userMaxUSN).Order("usn ASC").Limit(limit).Find(&books).Error; err != nil {
		return SyncFragment{}, nil
	}

	var items []usnItem
	for _, note := range notes {
		i := usnItem{
			usn: note.USN,
			val: note,
		}
		items = append(items, i)
	}
	for _, book := range books {
		i := usnItem{
			usn: book.USN,
			val: book,
		}
		items = append(items, i)
	}

	// order by usn in ascending order
	sort.Slice(items, func(i, j int) bool {
		return items[i].usn < items[j].usn
	})

	fragNotes := []SyncFragNote{}
	fragBooks := []SyncFragBook{}
	fragExpungedNotes := []string{}
	fragExpungedBooks := []string{}

	fragMaxUSN := 0
	for i := 0; i < limit; i++ {
		if i > len(items)-1 {
			break
		}

		item := items[i]

		fragMaxUSN = item.usn

		switch v := item.val.(type) {
		case database.Note:
			note := item.val.(database.Note)

			if note.Deleted {
				fragExpungedNotes = append(fragExpungedNotes, note.UUID)
			} else {
				fragNotes = append(fragNotes, NewFragNote(note))
			}
		case database.Book:
			book := item.val.(database.Book)

			if book.Deleted {
				fragExpungedBooks = append(fragExpungedBooks, book.UUID)
			} else {
				fragBooks = append(fragBooks, NewFragBook(book))
			}
		default:
			return SyncFragment{}, errors.Errorf("unknown internal item type %s", v)
		}
	}

	ret := SyncFragment{
		FragMaxUSN:    fragMaxUSN,
		UserMaxUSN:    userMaxUSN,
		CurrentTime:   a.Clock.Now().Unix(),
		Notes:         fragNotes,
		Books:         fragBooks,
		ExpungedNotes: fragExpungedNotes,
		ExpungedBooks: fragExpungedBooks,
	}

	return ret, nil
}

func parseGetSyncFragmentQuery(q url.Values) (afterUSN, limit int, err error) {
	afterUSNStr := q.Get("after_usn")
	limitStr := q.Get("limit")

	if len(afterUSNStr) > 0 {
		afterUSN, err = strconv.Atoi(afterUSNStr)

		if err != nil {
			err = errors.Wrap(err, "invalid after_usn")
			return
		}
	} else {
		afterUSN = 0
	}

	if len(limitStr) > 0 {
		l, e := strconv.Atoi(limitStr)

		if e != nil {
			err = errors.Wrap(e, "invalid limit")
			return
		}

		if l > 100 {
			err = &queryParamError{
				key:     "limit",
				value:   limitStr,
				message: "maximum value is 100",
			}
			return
		}

		limit = l
	} else {
		limit = 100
	}

	return
}

// GetSyncFragmentResp represents a response from GetSyncFragment handler
type GetSyncFragmentResp struct {
	Fragment SyncFragment `json:"fragment"`
}

// GetSyncFragment responds with a sync fragment
func (a *App) GetSyncFragment(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	afterUSN, limit, err := parseGetSyncFragmentQuery(r.URL.Query())
	if err != nil {
		handleError(w, "parsing query params", err, http.StatusInternalServerError)
		return
	}

	fragment, err := a.newFragment(user.ID, user.MaxUSN, afterUSN, limit)
	if err != nil {
		handleError(w, "getting fragment", err, http.StatusInternalServerError)
		return
	}

	response := GetSyncFragmentResp{
		Fragment: fragment,
	}
	respondJSON(w, response)
}

// GetSyncStateResp represents a response from GetSyncFragment handler
type GetSyncStateResp struct {
	FullSyncBefore int   `json:"full_sync_before"`
	MaxUSN         int   `json:"max_usn"`
	CurrentTime    int64 `json:"current_time"`
}

// GetSyncState responds with a sync fragment
func (a *App) GetSyncState(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		handleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	response := GetSyncStateResp{
		FullSyncBefore: fullSyncBefore,
		MaxUSN:         user.MaxUSN,
		// TODO: exposing server time means we probably shouldn't seed random generator with time?
		CurrentTime: a.Clock.Now().Unix(),
	}

	log.WithFields(log.Fields{
		"user_id": user.ID,
		"resp":    response,
	}).Info("getting sync state")

	respondJSON(w, response)
}
