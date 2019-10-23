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

package repetition

import (
	"fmt"
	"time"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// BuildEmail builds an email for the spaced repetition
func BuildEmail(now time.Time, user database.User, emailAddr string, digest database.Digest, rule database.RepetitionRule) (*mailer.Email, error) {
	date := now.Format("Jan 02 2006")
	subject := fmt.Sprintf("%s %s", rule.Title, date)
	tok, err := mailer.GetEmailPreferenceToken(user)
	if err != nil {
		return nil, errors.Wrap(err, "getting email frequency token")
	}

	threshold1 := now.AddDate(0, 0, -1).UnixNano()
	threshold2 := now.AddDate(0, 0, -3).UnixNano()
	threshold3 := now.AddDate(0, 0, -7).UnixNano()

	noteInfos := []mailer.DigestNoteInfo{}
	for _, note := range digest.Notes {
		var stage int
		if note.AddedOn > threshold2 && note.AddedOn < threshold1 {
			stage = 1
		} else if note.AddedOn > threshold3 && note.AddedOn < threshold2 {
			stage = 2
		} else if note.AddedOn > threshold3 {
			stage = 3
		}

		info := mailer.NewNoteInfo(note, stage)
		noteInfos = append(noteInfos, info)
	}

	bookCount := 0
	bookMap := map[string]bool{}
	for _, n := range digest.Notes {
		if ok := bookMap[n.Book.Label]; !ok {
			bookCount++
			bookMap[n.Book.Label] = true
		}
	}

	tmplData := mailer.DigestTmplData{
		Subject:           subject,
		NoteInfo:          noteInfos,
		ActiveBookCount:   bookCount,
		ActiveNoteCount:   len(digest.Notes),
		EmailSessionToken: tok.Value,
	}

	email := mailer.NewEmail("noreply@getdnote.com", []string{emailAddr}, subject)
	if err := email.ParseTemplate(mailer.EmailTypeWeeklyDigest, tmplData); err != nil {
		return nil, err
	}

	return email, nil
}

func getEligibleRules(now time.Time) ([]database.RepetitionRule, error) {
	hour := now.Hour()
	minute := now.Minute()

	var ret []database.RepetitionRule
	db := database.DBConn
	if err := db.Where("hour = ? AND minute = ?", hour, minute).Find(&ret).Error; err != nil {
		return nil, errors.Wrap(err, "querying db")
	}

	return ret, nil
}

func build(tx *gorm.DB, rule database.RepetitionRule) (database.Digest, error) {
	notes, err := getRandomNotes(tx, rule)
	if err != nil {
		return database.Digest{}, errors.Wrap(err, "getting notes")
	}

	digest := database.Digest{
		RuleID: rule.ID,
		UserID: rule.UserID,
		Notes:  notes,
	}
	if err := tx.Save(&digest).Error; err != nil {
		return database.Digest{}, errors.Wrap(err, "saving digest")
	}

	return digest, nil
}

func notify(now time.Time, user database.User, digest database.Digest, rule database.RepetitionRule) error {
	db := database.DBConn

	var account database.Account
	if err := db.Where("user_id = ?", user.ID).First(&account).Error; err != nil {
		return errors.Wrap(err, "getting account")
	}

	if !account.Email.Valid || !account.EmailVerified {
		log.WithFields(log.Fields{
			"user_id": user.ID,
		}).Info("Skipping repetition delivery because email is not valid or verified")
		return nil
	}

	email, err := BuildEmail(now, user, account.Email.String, digest, rule)
	if err != nil {
		return errors.Wrap(err, "making email")
	}

	err = email.Send()
	if err != nil {
		return errors.Wrap(err, "sending email")
	}

	notif := database.Notification{
		Type:   "email_weekly",
		UserID: user.ID,
	}

	if err := db.Create(&notif).Error; err != nil {
		return errors.Wrap(err, "creating notification")
	}

	return nil
}

func checkCooldown(now time.Time, rule database.RepetitionRule) bool {
	present := now.UnixNano() / int64(time.Millisecond)

	return present >= int64(rule.LastActive+rule.Frequency)
}

func touchLastActive(tx *gorm.DB, rule database.RepetitionRule, now time.Time) error {
	lastActive := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	rule.LastActive = lastActive.UnixNano() / int64(time.Millisecond)

	if err := tx.Save(&rule).Error; err != nil {
		return errors.Wrap(err, "updating repetition rule")
	}

	return nil
}

func process(now time.Time, rule database.RepetitionRule) error {
	log.WithFields(log.Fields{
		"uuid": rule.UUID,
	}).Info("processing repetition")

	db := database.DBConn
	tx := db.Begin()

	if !checkCooldown(now, rule) {
		return nil
	}

	var user database.User
	if err := tx.Where("id = ?", rule.UserID).First(&user).Error; err != nil {
		return errors.Wrap(err, "getting user")
	}
	if !user.Cloud {
		log.WithFields(log.Fields{
			"user_id": user.ID,
		}).Info("Skipping repetition due to lack of subscription")
		return nil
	}

	digest, err := build(tx, rule)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "building repetition")
	}

	if err := touchLastActive(tx, rule, now); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "touching last_active")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "committing transaction")
	}

	if err := notify(now, user, digest, rule); err != nil {
		return errors.Wrap(err, "notifying user")
	}

	log.WithFields(log.Fields{
		"uuid": rule.UUID,
	}).Info("finished processing repetition")

	return nil
}

// Do creates spaced repetitions and delivers the results based on the rules
func Do(c clock.Clock) error {
	now := c.Now().UTC()

	rules, err := getEligibleRules(now)
	if err != nil {
		return errors.Wrap(err, "getting eligible repetition rules")
	}

	log.WithFields(log.Fields{
		"hour":      now.Hour(),
		"minute":    now.Minute(),
		"num_rules": len(rules),
	}).Info("processing rules")

	for _, rule := range rules {
		if err := process(now, rule); err != nil {
			log.WithFields(log.Fields{
				"rule uuid": rule.UUID,
			}).ErrorWrap(err, "Could not process the repetition rule")
			continue
		}
	}

	return nil
}
