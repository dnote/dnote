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
	"os"
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/job/ctx"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// BuildEmailParams is the params for building an email
type BuildEmailParams struct {
	Now    time.Time
	User   database.User
	Digest database.Digest
	Rule   database.RepetitionRule
}

// BuildEmail builds an email for the spaced repetition
func BuildEmail(c ctx.Ctx, p BuildEmailParams) (string, string, error) {
	date := p.Now.Format("Jan 02 2006")
	subject := fmt.Sprintf("%s %s", p.Rule.Title, date)
	tok, err := mailer.GetToken(c.DB, p.User, database.TokenTypeRepetition)
	if err != nil {
		return "", "", errors.Wrap(err, "getting email frequency token")
	}

	t1 := p.Now.AddDate(0, 0, -3).UnixNano()
	t2 := p.Now.AddDate(0, 0, -7).UnixNano()

	noteInfos := []mailer.DigestNoteInfo{}
	for _, note := range p.Digest.Notes {
		var stage int
		if note.AddedOn > t1 {
			stage = 1
		} else if note.AddedOn > t2 && note.AddedOn < t1 {
			stage = 2
		} else if note.AddedOn < t2 {
			stage = 3
		}

		info := mailer.NewNoteInfo(note, stage)
		noteInfos = append(noteInfos, info)
	}

	bookCount := 0
	bookMap := map[string]bool{}
	for _, n := range p.Digest.Notes {
		if ok := bookMap[n.Book.Label]; !ok {
			bookCount++
			bookMap[n.Book.Label] = true
		}
	}

	tmplData := mailer.DigestTmplData{
		Subject:           subject,
		NoteInfo:          noteInfos,
		ActiveBookCount:   bookCount,
		ActiveNoteCount:   len(p.Digest.Notes),
		EmailSessionToken: tok.Value,
		RuleUUID:          p.Rule.UUID,
		RuleTitle:         p.Rule.Title,
		WebURL:            os.Getenv("WebURL"),
	}

	body, err := c.EmailTmpl.Execute(mailer.EmailTypeWeeklyDigest, tmplData)
	if err != nil {
		return "", "", errors.Wrap(err, "executing digest email template")
	}

	return subject, body, nil
}

func getEligibleRules(db *gorm.DB, now time.Time) ([]database.RepetitionRule, error) {
	hour := now.Hour()
	minute := now.Minute()

	var ret []database.RepetitionRule
	if err := db.
		Where("users.cloud AND repetition_rules.hour = ? AND repetition_rules.minute = ? AND repetition_rules.enabled", hour, minute).
		Joins("INNER JOIN users ON users.id = repetition_rules.user_id").
		Find(&ret).Error; err != nil {
		return nil, errors.Wrap(err, "querying db")
	}

	return ret, nil
}

func build(tx *gorm.DB, rule database.RepetitionRule) (database.Digest, error) {
	notes, err := getBalancedNotes(tx, rule)
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

func notify(c ctx.Ctx, now time.Time, user database.User, digest database.Digest, rule database.RepetitionRule) error {
	var account database.Account
	if err := c.DB.Where("user_id = ?", user.ID).First(&account).Error; err != nil {
		return errors.Wrap(err, "getting account")
	}

	if !account.Email.Valid || !account.EmailVerified {
		log.WithFields(log.Fields{
			"user_id": user.ID,
		}).Info("Skipping repetition delivery because email is not valid or verified")
		return nil
	}

	subject, body, err := BuildEmail(c, BuildEmailParams{
		Now:    now,
		User:   user,
		Digest: digest,
		Rule:   rule,
	})
	if err != nil {
		return errors.Wrap(err, "making email")
	}

	if err := c.EmailBackend.Queue(subject, "noreply@getdnote.com", []string{account.Email.String}, body); err != nil {
		return errors.Wrap(err, "queueing email")
	}

	notif := database.Notification{
		Type:   "email_weekly",
		UserID: user.ID,
	}
	if err := c.DB.Create(&notif).Error; err != nil {
		return errors.Wrap(err, "creating notification")
	}

	return nil
}

func checkCooldown(now time.Time, rule database.RepetitionRule) bool {
	present := now.UnixNano() / int64(time.Millisecond)

	return present >= rule.NextActive
}

func getNextActive(base int64, frequency int64, now time.Time) int64 {
	candidate := base + frequency
	if candidate >= now.UnixNano()/int64(time.Millisecond) {
		return candidate
	}

	return getNextActive(candidate, frequency, now)
}

func touchTimestamp(tx *gorm.DB, rule database.RepetitionRule, now time.Time) error {
	lastActive := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location()).UnixNano() / int64(time.Millisecond)

	rule.LastActive = lastActive
	rule.NextActive = getNextActive(rule.LastActive, rule.Frequency, now)

	if err := tx.Save(&rule).Error; err != nil {
		return errors.Wrap(err, "updating repetition rule")
	}

	return nil
}

func process(c ctx.Ctx, now time.Time, rule database.RepetitionRule) error {
	log.WithFields(log.Fields{
		"uuid": rule.UUID,
	}).Info("processing repetition")

	tx := c.DB.Begin()

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

	if err := touchTimestamp(tx, rule, now); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "touching last_active")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "committing transaction")
	}

	if err := notify(c, now, user, digest, rule); err != nil {
		return errors.Wrap(err, "notifying user")
	}

	log.WithFields(log.Fields{
		"uuid": rule.UUID,
	}).Info("finished processing repetition")

	return nil
}

// Do creates spaced repetitions and delivers the results based on the rules
func Do(c ctx.Ctx) error {
	now := c.Clock.Now().UTC()

	rules, err := getEligibleRules(c.DB, now)
	if err != nil {
		return errors.Wrap(err, "getting eligible repetition rules")
	}

	log.WithFields(log.Fields{
		"hour":      now.Hour(),
		"minute":    now.Minute(),
		"num_rules": len(rules),
	}).Info("processing rules")

	for _, rule := range rules {
		if err := process(c, now, rule); err != nil {
			log.WithFields(log.Fields{
				"rule uuid": rule.UUID,
			}).ErrorWrap(err, "Could not process the repetition rule")
			continue
		}
	}

	return nil
}
