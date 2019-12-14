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
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/presenters"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (a *API) getRepetitionRule(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	repetitionRuleUUID := vars["repetitionRuleUUID"]

	if ok := helpers.ValidateUUID(repetitionRuleUUID); !ok {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	var repetitionRule database.RepetitionRule
	if err := a.App.DB.Where("user_id = ? AND uuid = ?", user.ID, repetitionRuleUUID).Preload("Books").Find(&repetitionRule).Error; err != nil {
		HandleError(w, "getting repetition rules", err, http.StatusInternalServerError)
		return
	}

	resp := presenters.PresentRepetitionRule(repetitionRule)
	respondJSON(w, http.StatusOK, resp)
}

func (a *API) getRepetitionRules(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	var repetitionRules []database.RepetitionRule
	if err := a.App.DB.Where("user_id = ?", user.ID).Preload("Books").Order("last_active DESC").Find(&repetitionRules).Error; err != nil {
		HandleError(w, "getting repetition rules", err, http.StatusInternalServerError)
		return
	}

	resp := presenters.PresentRepetitionRules(repetitionRules)
	respondJSON(w, http.StatusOK, resp)
}

func validateBookDomain(val string) error {
	if val == database.BookDomainAll || val == database.BookDomainIncluding || val == database.BookDomainExluding {
		return nil
	}

	return errors.Errorf("invalid book_domain %s", val)
}

type repetitionRuleParams struct {
	Title      *string   `json:"title"`
	Enabled    *bool     `json:"enabled"`
	Hour       *int      `json:"hour"`
	Minute     *int      `json:"minute"`
	Frequency  *int64    `json:"frequency"`
	BookDomain *string   `json:"book_domain"`
	BookUUIDs  *[]string `json:"book_uuids"`
	NoteCount  *int      `json:"note_count"`
}

func (r repetitionRuleParams) GetEnabled() bool {
	if r.Enabled == nil {
		return false
	}

	return *r.Enabled
}

func (r repetitionRuleParams) GetFrequency() int64 {
	if r.Frequency == nil {
		return 0
	}

	return *r.Frequency
}

func (r repetitionRuleParams) GetTitle() string {
	if r.Title == nil {
		return ""
	}

	return *r.Title
}

func (r repetitionRuleParams) GetNoteCount() int {
	if r.NoteCount == nil {
		return 0
	}

	return *r.NoteCount
}

func (r repetitionRuleParams) GetBookDomain() string {
	if r.BookDomain == nil {
		return ""
	}

	return *r.BookDomain
}

func (r repetitionRuleParams) GetBookUUIDs() []string {
	if r.BookUUIDs == nil {
		return []string{}
	}

	return *r.BookUUIDs
}

func (r repetitionRuleParams) GetHour() int {
	if r.Hour == nil {
		return 0
	}

	return *r.Hour
}

func (r repetitionRuleParams) GetMinute() int {
	if r.Minute == nil {
		return 0
	}

	return *r.Minute
}

func validateRepetitionRuleParams(p repetitionRuleParams) error {
	if p.Frequency != nil && p.GetFrequency() == 0 {
		return errors.New("frequency is required")
	}

	if p.Title != nil {
		title := p.GetTitle()

		if len(title) == 0 {
			return errors.New("Title is required")
		}
		if len(title) > 50 {
			return errors.New("Title is too long")
		}
	}

	if p.NoteCount != nil && p.GetNoteCount() == 0 {
		return errors.New("note count has to be greater than 0")
	}

	if p.BookDomain != nil {
		bookDomain := p.GetBookDomain()
		if err := validateBookDomain(bookDomain); err != nil {
			return err
		}

		bookUUIDs := p.GetBookUUIDs()
		if bookDomain == database.BookDomainAll {
			if len(bookUUIDs) > 0 {
				return errors.New("a global repetition should not specify book_uuids")
			}
		} else {
			if len(bookUUIDs) == 0 {
				return errors.New("book_uuids is required")
			}
		}
	}

	if p.Hour != nil {
		hour := p.GetHour()

		if hour < 0 && hour > 23 {
			return errors.New("invalid hour")
		}
	}

	if p.Minute != nil {
		minute := p.GetMinute()

		if minute < 0 && minute > 60 {
			return errors.New("invalid minute")
		}
	}

	return nil
}

func validateCreateRepetitionRuleParams(p repetitionRuleParams) error {
	if p.Title == nil {
		return errors.New("title is required")
	}
	if p.Frequency == nil {
		return errors.New("frequency is required")
	}
	if p.NoteCount == nil {
		return errors.New("note_count is required")
	}
	if p.BookDomain == nil {
		return errors.New("book_domain is required")
	}
	if p.Hour == nil {
		return errors.New("hour is required")
	}
	if p.Minute == nil {
		return errors.New("minute is required")
	}
	if p.Enabled == nil {
		return errors.New("enabled is required")
	}

	return nil
}

func parseCreateRepetitionRuleParams(r *http.Request) (repetitionRuleParams, error) {
	var ret repetitionRuleParams

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	if err := d.Decode(&ret); err != nil {
		return ret, errors.Wrap(err, "decoding json")
	}

	if err := validateCreateRepetitionRuleParams(ret); err != nil {
		return ret, errors.Wrap(err, "validating params")
	}

	if err := validateRepetitionRuleParams(ret); err != nil {
		return ret, errors.Wrap(err, "validating params")
	}

	return ret, nil
}

type calcNextActiveParams struct {
	Hour      int
	Minute    int
	Frequency int64
}

// calcNextActive calculates the NextActive value for a repetition rule by adding the given
// frequency to the given present date time at the given hour and minute.
func calcNextActive(now time.Time, p calcNextActiveParams) int64 {
	t0 := time.Date(now.Year(), now.Month(), now.Day(), p.Hour, p.Minute, 0, 0, now.Location()).UnixNano() / int64(time.Millisecond)

	return t0 + p.Frequency
}

func (a *API) createRepetitionRule(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	params, err := parseCreateRepetitionRuleParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var books []database.Book
	if err := a.App.DB.Where("user_id = ? AND uuid IN (?)", user.ID, params.GetBookUUIDs()).Find(&books).Error; err != nil {
		HandleError(w, "finding books", nil, http.StatusInternalServerError)
		return
	}

	nextActive := calcNextActive(a.App.Clock.Now(), calcNextActiveParams{
		Hour:      params.GetHour(),
		Minute:    params.GetMinute(),
		Frequency: params.GetFrequency(),
	})

	record := database.RepetitionRule{
		UserID:     user.ID,
		Title:      params.GetTitle(),
		Hour:       params.GetHour(),
		Minute:     params.GetMinute(),
		Frequency:  params.GetFrequency(),
		BookDomain: params.GetBookDomain(),
		NextActive: nextActive,
		Books:      books,
		NoteCount:  params.GetNoteCount(),
		Enabled:    params.GetEnabled(),
	}
	if err := a.App.DB.Create(&record).Error; err != nil {
		HandleError(w, "creating a repetition rule", err, http.StatusInternalServerError)
		return
	}

	resp := presenters.PresentRepetitionRule(record)
	respondJSON(w, http.StatusCreated, resp)
}

func parseUpdateDigestParams(r *http.Request) (repetitionRuleParams, error) {
	var ret repetitionRuleParams

	if err := json.NewDecoder(r.Body).Decode(&ret); err != nil {
		return ret, errors.Wrap(err, "decoding json")
	}

	if err := validateRepetitionRuleParams(ret); err != nil {
		return ret, errors.Wrap(err, "validating params")
	}

	return ret, nil
}

func (a *API) deleteRepetitionRule(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	repetitionRuleUUID := vars["repetitionRuleUUID"]

	var rule database.RepetitionRule
	conn := a.App.DB.Where("uuid = ? AND user_id = ?", repetitionRuleUUID, user.ID).First(&rule)

	if conn.RecordNotFound() {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	} else if err := conn.Error; err != nil {
		HandleError(w, "finding the repetition rule", err, http.StatusInternalServerError)
		return
	}

	if err := a.App.DB.Exec("DELETE from repetition_rules WHERE uuid = ?", rule.UUID).Error; err != nil {
		HandleError(w, "deleting the repetition rule", err, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) updateRepetitionRule(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		HandleError(w, "No authenticated user found", nil, http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	repetitionRuleUUID := vars["repetitionRuleUUID"]

	params, err := parseUpdateDigestParams(r)
	if err != nil {
		http.Error(w, "parsing params", http.StatusBadRequest)
		return
	}

	tx := a.App.DB.Begin()

	var repetitionRule database.RepetitionRule
	if err := tx.Where("user_id = ? AND uuid = ?", user.ID, repetitionRuleUUID).Preload("Books").First(&repetitionRule).Error; err != nil {
		HandleError(w, "finding record", nil, http.StatusInternalServerError)
		return
	}

	if params.Title != nil {
		repetitionRule.Title = params.GetTitle()
	}
	if params.Enabled != nil {
		enabled := params.GetEnabled()
		repetitionRule.Enabled = enabled

		if enabled && !repetitionRule.Enabled {
			repetitionRule.NextActive = calcNextActive(a.App.Clock.Now(), calcNextActiveParams{
				Hour:      repetitionRule.Hour,
				Minute:    repetitionRule.Minute,
				Frequency: repetitionRule.Frequency,
			})
		} else if !enabled && repetitionRule.Enabled {
			repetitionRule.NextActive = 0
		}
	}
	if params.Hour != nil {
		repetitionRule.Hour = params.GetHour()
	}
	if params.Minute != nil {
		repetitionRule.Minute = params.GetMinute()
	}
	if params.Frequency != nil {
		frequency := params.GetFrequency()

		repetitionRule.Frequency = frequency
		repetitionRule.NextActive = calcNextActive(a.App.Clock.Now(), calcNextActiveParams{
			Hour:      repetitionRule.Hour,
			Minute:    repetitionRule.Minute,
			Frequency: frequency,
		})
	}
	if params.NoteCount != nil {
		repetitionRule.NoteCount = params.GetNoteCount()
	}
	if params.BookDomain != nil {
		repetitionRule.BookDomain = params.GetBookDomain()
	}
	if params.BookUUIDs != nil {
		var books []database.Book
		if err := tx.Where("user_id = ? AND uuid IN (?)", user.ID, *params.BookUUIDs).Find(&books).Error; err != nil {
			HandleError(w, "finding books", err, http.StatusInternalServerError)
			return
		}

		if err := tx.Model(&repetitionRule).Association("Books").Replace(books).Error; err != nil {
			tx.Rollback()
			HandleError(w, "updating books association for a repetitionRule", err, http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Save(&repetitionRule).Error; err != nil {
		tx.Rollback()
		HandleError(w, "creating a repetition rule", err, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		HandleError(w, "committing a transaction", err, http.StatusInternalServerError)
	}

	resp := presenters.PresentRepetitionRule(repetitionRule)
	respondJSON(w, http.StatusOK, resp)
}
