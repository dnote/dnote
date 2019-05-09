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
	"net/http"

	"github.com/dnote/dnote/server/api/helpers"
	"github.com/dnote/dnote/server/api/operations"
	"github.com/dnote/dnote/server/api/presenters"
	"github.com/dnote/dnote/server/database"
	"github.com/pkg/errors"
)

type createBookV2Payload struct {
	Name string `json:"name"`
}

// CreateBookV2Resp is the response from create book api
type CreateBookV2Resp struct {
	Book presenters.Book `json:"book"`
}

func validateCreateBookV2Payload(p createBookPayload) error {
	if p.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

// CreateBookV2 creates a new book
func (a *App) CreateBookV2(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		return
	}

	var params createBookPayload
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, errors.Wrap(err, "decoding payload").Error(), http.StatusInternalServerError)
		return
	}

	err = validateCreateBookPayload(params)
	if err != nil {
		http.Error(w, errors.Wrap(err, "validating payload").Error(), http.StatusBadRequest)
		return
	}

	db := database.DBConn

	var bookCount int
	err = db.Model(database.Book{}).
		Where("user_id = ? AND label = ?", user.ID, params.Name).
		Count(&bookCount).Error
	if err != nil {
		http.Error(w, errors.Wrap(err, "checking duplicate").Error(), http.StatusInternalServerError)
		return
	}
	if bookCount > 0 {
		http.Error(w, "duplicate book exists", http.StatusConflict)
		return
	}

	book, err := operations.CreateBook(user, a.Clock, params.Name)
	if err != nil {
		http.Error(w, errors.Wrap(err, "inserting book").Error(), http.StatusInternalServerError)
	}
	resp := CreateBookResp{
		Book: presenters.PresentBook(book),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// BooksOptionsV2 is a handler for OPTIONS endpoint for notes
func (a *App) BooksOptionsV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Version")
}
