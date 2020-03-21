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

package remind

import (
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Context holds data that repetition job needs in order to perform
type Context struct {
	DB           *gorm.DB
	Clock        clock.Clock
	EmailTmpl    mailer.Templates
	EmailBackend mailer.Backend
	Config       config.Config
}

type inactiveUserInfo struct {
	userID         int
	email          string
	sampleNoteUUID string
}

func (c *Context) sampleUserNote(userID int) (database.Note, error) {
	var ret database.Note
	// FIXME: ordering by random() requires a sequential scan on the whole table and does not scale
	if err := c.DB.Where("user_id = ?", userID).Order("random() DESC").First(&ret).Error; err != nil {
		return ret, errors.Wrap(err, "getting a random note")
	}

	return ret, nil
}

func (c *Context) getInactiveUserInfo() ([]inactiveUserInfo, error) {
	ret := []inactiveUserInfo{}

	threshold := c.Clock.Now().AddDate(0, 0, -14).Unix()

	rows, err := c.DB.Raw(`
SELECT
	notes.user_id AS user_id,
	accounts.email,
	SUM(
		CASE
			WHEN notes.created_at > to_timestamp(?) THEN 1
			ELSE 0
		END
	) AS recent_note_count,
	COUNT(*) AS total_note_count
FROM notes
INNER JOIN accounts ON accounts.user_id = notes.user_id
WHERE accounts.email IS NOT NULL AND accounts.email_verified IS TRUE
GROUP BY notes.user_id, accounts.email`, threshold).Rows()
	if err != nil {
		return ret, errors.Wrap(err, "executing note count SQL query")
	}
	defer rows.Close()
	for rows.Next() {
		var userID, recentNoteCount, totalNoteCount int
		var email string
		if err := rows.Scan(&userID, &email, &recentNoteCount, &totalNoteCount); err != nil {
			return nil, errors.Wrap(err, "scanning a row")
		}

		if recentNoteCount == 0 && totalNoteCount > 0 {
			note, err := c.sampleUserNote(userID)
			if err != nil {
				return nil, errors.Wrap(err, "sampling user note")
			}

			ret = append(ret, inactiveUserInfo{
				userID:         userID,
				email:          email,
				sampleNoteUUID: note.UUID,
			})
		}
	}

	return ret, nil
}

func (c *Context) canNotify(info inactiveUserInfo) (bool, error) {
	var pref database.EmailPreference
	if err := c.DB.Where("user_id = ?", info.userID).First(&pref).Error; err != nil {
		return false, errors.Wrap(err, "getting email preference")
	}

	if !pref.InactiveReminder {
		return false, nil
	}

	var notif database.Notification
	conn := c.DB.Where("user_id = ? AND type = ?", info.userID, mailer.EmailTypeInactiveReminder).Order("created_at DESC").First(&notif)

	if conn.RecordNotFound() {
		return true, nil
	} else if err := conn.Error; err != nil {
		return false, errors.Wrap(err, "checking cooldown")
	}

	t := c.Clock.Now().AddDate(0, 0, -14)
	if notif.CreatedAt.Before(t) {
		return true, nil
	}

	return false, nil
}

func (c *Context) process(info inactiveUserInfo) error {
	ok, err := c.canNotify(info)
	if err != nil {
		return errors.Wrap(err, "checking if user can be notified")
	}
	if !ok {
		return nil
	}

	sender, err := app.GetSenderEmail(c.Config, "noreply@getdnote.com")
	if err != nil {
		return errors.Wrap(err, "getting sender email")
	}

	tok, err := mailer.GetToken(c.DB, info.userID, database.TokenTypeEmailPreference)
	if err != nil {
		return errors.Wrap(err, "getting email token")
	}

	tmplData := mailer.InactiveReminderTmplData{
		WebURL:         c.Config.WebURL,
		SampleNoteUUID: info.sampleNoteUUID,
		Token:          tok.Value,
	}
	body, err := c.EmailTmpl.Execute(mailer.EmailTypeInactiveReminder, mailer.EmailKindText, tmplData)
	if err != nil {
		return errors.Wrap(err, "executing inactive email template")
	}

	if err := c.EmailBackend.Queue("Your Dnote stopped growing", sender, []string{info.email}, mailer.EmailKindText, body); err != nil {
		return errors.Wrap(err, "queueing email")
	}

	if err := c.DB.Create(&database.Notification{
		Type:   mailer.EmailTypeInactiveReminder,
		UserID: info.userID,
	}).Error; err != nil {
		return errors.Wrap(err, "creating notification")
	}

	return nil
}

// Result holds the result of the job
type Result struct {
	SuccessCount  int
	FailedUserIDs []int
}

// DoInactive sends reminder for users with no recent notes
func DoInactive(c Context) (Result, error) {
	log.Info("performing reminder for no recent notes")

	result := Result{}
	items, err := c.getInactiveUserInfo()
	if err != nil {
		return result, errors.Wrap(err, "getting inactive user information")
	}

	log.WithFields(log.Fields{
		"user_count": len(items),
	}).Info("counted inactive users")

	for _, item := range items {
		err := c.process(item)

		if err == nil {
			result.SuccessCount = result.SuccessCount + 1
		} else {
			log.WithFields(log.Fields{
				"user_id": item.userID,
			}).ErrorWrap(err, "Could not process no recent notes reminder")

			result.FailedUserIDs = append(result.FailedUserIDs, item.userID)
		}
	}

	return result, nil
}
