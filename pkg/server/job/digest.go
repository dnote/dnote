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

package job

import (
	"log"
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/pkg/errors"
)

// MakeDigest builds a weekly digest email
func MakeDigest(user database.User, emailAddr string) (*mailer.Email, error) {
	log.Printf("Sending for %s", emailAddr)
	db := database.DBConn

	subject := "Weekly Digest"
	tok, err := mailer.GetEmailPreferenceToken(user)
	if err != nil {
		return nil, errors.Wrap(err, "getting email frequency token")
	}

	now := time.Now()
	threshold1 := int(now.AddDate(0, 0, -1).UnixNano())
	threshold2 := int(now.AddDate(0, 0, -3).UnixNano())
	threshold3 := int(now.AddDate(0, 0, -7).UnixNano())

	var stage1 []database.Note
	var stage2 []database.Note
	var stage3 []database.Note

	// TODO: ordering by random() does not scale if table grows large
	if err := db.Where("user_id = ? AND added_on > ? AND added_on < ?", user.ID, threshold2, threshold1).Order("random()").Limit(4).Preload("Book").Find(&stage1).Error; err != nil {
		return nil, errors.Wrap(err, "Failed to get notes with threshold 1")
	}
	if err := db.Where("user_id = ? AND added_on > ? AND added_on < ?", user.ID, threshold3, threshold2).Order("random()").Limit(4).Preload("Book").Find(&stage2).Error; err != nil {
		return nil, errors.Wrap(err, "Failed to get notes with threshold 2")
	}
	if err := db.Where("user_id = ? AND added_on < ?", user.ID, threshold3).Order("random()").Limit(4).Preload("Book").Find(&stage3).Error; err != nil {
		return nil, errors.Wrap(err, "Failed to get notes with threshold 3")
	}

	noteInfos := []mailer.DigestNoteInfo{}
	for _, note := range stage1 {
		info := mailer.NewNoteInfo(note, 1)
		noteInfos = append(noteInfos, info)
	}
	for _, note := range stage2 {
		info := mailer.NewNoteInfo(note, 2)
		noteInfos = append(noteInfos, info)
	}
	for _, note := range stage3 {
		info := mailer.NewNoteInfo(note, 3)
		noteInfos = append(noteInfos, info)
	}

	notes := append(append(stage1, stage2...), stage3...)
	digest := database.Digest{
		UserID: user.ID,
		Notes:  notes,
	}
	if err := db.Save(&digest).Error; err != nil {
		return nil, errors.Wrap(err, "saving digest")
	}

	bookCount := 0
	bookMap := map[string]bool{}
	for _, n := range notes {
		if ok := bookMap[n.Book.Label]; !ok {
			bookCount++
			bookMap[n.Book.Label] = true
		}
	}

	tmplData := mailer.DigestTmplData{
		Subject:           subject,
		NoteInfo:          noteInfos,
		ActiveBookCount:   bookCount,
		ActiveNoteCount:   len(notes),
		EmailSessionToken: tok.Value,
	}

	email := mailer.NewEmail("notebot@getdnote.com", []string{emailAddr}, subject)
	if err := email.ParseTemplate(mailer.EmailTypeWeeklyDigest, tmplData); err != nil {
		return nil, err
	}

	return email, nil
}

// sendDigest sends the weekly digests to users
func sendDigest() error {
	db := database.DBConn

	var users []database.User
	if err := db.
		Preload("Account").
		Where("cloud = ?", true).
		Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		account := user.Account

		if !account.Email.Valid || !account.EmailVerified {
			continue
		}

		email, err := MakeDigest(user, account.Email.String)
		if err != nil {
			log.Printf("Error occurred while sending to %s: %s", account.Email.String, err.Error())
			continue
		}

		err = email.Send()
		if err != nil {
			log.Printf("Error occurred while sending to %s: %s", account.Email.String, err.Error())
			continue
		}

		notif := database.Notification{
			Type:   "email_weekly",
			UserID: user.ID,
		}

		if err := db.Create(&notif).Error; err != nil {
			log.Printf("Error occurred while creating notification for %s: %s", account.Email.String, err.Error())
		}
	}

	return nil
}
