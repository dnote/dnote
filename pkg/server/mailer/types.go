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

// DigestTmplData is a template data for digest emails
type DigestTmplData struct {
	EmailSessionToken string
	DigestUUID        string
	DigestVersion     int
	RuleUUID          string
	RuleTitle         string
	WebURL            string
}

// EmailVerificationTmplData is a template data for email verification emails
type EmailVerificationTmplData struct {
	Token  string
	WebURL string
}

// EmailResetPasswordTmplData is a template data for reset password emails
type EmailResetPasswordTmplData struct {
	AccountEmail string
	Token        string
	WebURL       string
}

// EmailResetPasswordAlertTmplData is a template data for reset password emails
type EmailResetPasswordAlertTmplData struct {
	AccountEmail string
	WebURL       string
}

// WelcomeTmplData is a template data for welcome emails
type WelcomeTmplData struct {
	AccountEmail string
	WebURL       string
}

// EmailTypeSubscriptionConfirmationTmplData is a template data for reset password emails
type EmailTypeSubscriptionConfirmationTmplData struct {
	AccountEmail string
	WebURL       string
}
