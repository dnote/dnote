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
	"os"
	"sort"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/models"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func getTestContext(c clock.Clock, be *testutils.MockEmailbackendImplementation) Context {
	emailTmplDir := os.Getenv("DNOTE_TEST_EMAIL_TEMPLATE_DIR")

	con := Context{
		DB:           models.TestDB,
		Clock:        c,
		EmailTmpl:    mailer.NewTemplates(&emailTmplDir),
		EmailBackend: be,
	}

	return con
}

func TestDoInactive(t *testing.T) {
	defer models.ClearTestData(models.TestDB)

	t1 := time.Now()

	// u1 is an active user
	u1 := models.SetUpUserData()
	a1 := models.SetUpAccountData(u1, "alice@example.com", "pass1234")
	models.MustExec(t, models.TestDB.Model(&a1).Update("email_verified", true), "setting email verified")
	models.MustExec(t, models.TestDB.Save(&models.EmailPreference{UserID: u1.ID, InactiveReminder: true}), "preparing email preference")

	b1 := models.Book{
		UserID: u1.ID,
		Label:  "js",
	}
	models.MustExec(t, models.TestDB.Save(&b1), "preparing b1")
	n1 := models.Note{
		BookUUID: b1.UUID,
		UserID:   u1.ID,
	}
	models.MustExec(t, models.TestDB.Save(&n1), "preparing n1")

	// u2 is an inactive user
	u2 := models.SetUpUserData()
	a2 := models.SetUpAccountData(u2, "bob@example.com", "pass1234")
	models.MustExec(t, models.TestDB.Model(&a2).Update("email_verified", true), "setting email verified")
	models.MustExec(t, models.TestDB.Save(&models.EmailPreference{UserID: u2.ID, InactiveReminder: true}), "preparing email preference")

	b2 := models.Book{
		UserID: u2.ID,
		Label:  "css",
	}
	models.MustExec(t, models.TestDB.Save(&b2), "preparing b2")
	n2 := models.Note{
		UserID:   u2.ID,
		BookUUID: b2.UUID,
	}
	models.MustExec(t, models.TestDB.Save(&n2), "preparing n2")
	models.MustExec(t, models.TestDB.Model(&n2).Update("created_at", t1.AddDate(0, 0, -15)), "preparing n2")

	// u3 is an inactive user with inactive alert email preference disabled
	u3 := models.SetUpUserData()
	a3 := models.SetUpAccountData(u3, "alice@example.com", "pass1234")
	models.MustExec(t, models.TestDB.Model(&a3).Update("email_verified", true), "setting email verified")
	emailPref3 := models.EmailPreference{UserID: u3.ID}
	models.MustExec(t, models.TestDB.Save(&emailPref3), "preparing email preference")
	models.MustExec(t, models.TestDB.Model(&emailPref3).Update(map[string]interface{}{"inactive_reminder": false}), "updating email preference")

	b3 := models.Book{
		UserID: u3.ID,
		Label:  "js",
	}
	models.MustExec(t, models.TestDB.Save(&b3), "preparing b3")
	n3 := models.Note{
		BookUUID: b3.UUID,
		UserID:   u3.ID,
	}
	models.MustExec(t, models.TestDB.Save(&n3), "preparing n3")
	models.MustExec(t, models.TestDB.Model(&n3).Update("created_at", t1.AddDate(0, 0, -15)), "preparing n3")

	c := clock.NewMock()
	c.SetNow(t1)
	be := &testutils.MockEmailbackendImplementation{}

	con := getTestContext(c, be)
	if _, err := DoInactive(con); err != nil {
		t.Fatal(errors.Wrap(err, "performing"))
	}

	assert.Equalf(t, len(be.Emails), 1, "email queue count mismatch")
	assert.DeepEqual(t, be.Emails[0].To, []string{a2.Email.String}, "email address mismatch")
}

func TestDoInactive_Cooldown(t *testing.T) {
	defer models.ClearTestData(models.TestDB)

	// setup sets up an inactive user
	setup := func(t *testing.T, now time.Time, email string) models.User {
		u := models.SetUpUserData()
		a := models.SetUpAccountData(u, email, "pass1234")
		models.MustExec(t, models.TestDB.Model(&a).Update("email_verified", true), "setting email verified")
		models.MustExec(t, models.TestDB.Save(&models.EmailPreference{UserID: u.ID, InactiveReminder: true}), "preparing email preference")

		b := models.Book{
			UserID: u.ID,
			Label:  "css",
		}
		models.MustExec(t, models.TestDB.Save(&b), "preparing book")
		n := models.Note{
			UserID:   u.ID,
			BookUUID: b.UUID,
		}
		models.MustExec(t, models.TestDB.Save(&n), "preparing note")
		models.MustExec(t, models.TestDB.Model(&n).Update("created_at", now.AddDate(0, 0, -15)), "preparing note")

		return u
	}

	// Set up
	now := time.Now()

	setup(t, now, "alice@example.com")

	bob := setup(t, now, "bob@example.com")
	bobNotif := models.Notification{Type: mailer.EmailTypeInactiveReminder, UserID: bob.ID}
	models.MustExec(t, models.TestDB.Create(&bobNotif), "preparing inactive notification for bob")
	models.MustExec(t, models.TestDB.Model(&bobNotif).Update("created_at", now.AddDate(0, 0, -7)), "preparing created_at for inactive notification for bob")

	chuck := setup(t, now, "chuck@example.com")
	chuckNotif := models.Notification{Type: mailer.EmailTypeInactiveReminder, UserID: chuck.ID}
	models.MustExec(t, models.TestDB.Create(&chuckNotif), "preparing inactive notification for chuck")
	models.MustExec(t, models.TestDB.Model(&chuckNotif).Update("created_at", now.AddDate(0, 0, -15)), "preparing created_at for inactive notification for chuck")

	dan := setup(t, now, "dan@example.com")
	danNotif1 := models.Notification{Type: mailer.EmailTypeInactiveReminder, UserID: dan.ID}
	models.MustExec(t, models.TestDB.Create(&danNotif1), "preparing inactive notification 1 for dan")
	models.MustExec(t, models.TestDB.Model(&danNotif1).Update("created_at", now.AddDate(0, 0, -10)), "preparing created_at for inactive notification for dan")
	danNotif2 := models.Notification{Type: mailer.EmailTypeInactiveReminder, UserID: dan.ID}
	models.MustExec(t, models.TestDB.Create(&danNotif2), "preparing inactive notification 2 for dan")
	models.MustExec(t, models.TestDB.Model(&danNotif2).Update("created_at", now.AddDate(0, 0, -15)), "preparing created_at for inactive notification for dan")

	c := clock.NewMock()
	c.SetNow(now)
	be := &testutils.MockEmailbackendImplementation{}

	// Execute
	con := getTestContext(c, be)
	if _, err := DoInactive(con); err != nil {
		t.Fatal(errors.Wrap(err, "performing"))
	}

	// Test
	assert.Equalf(t, len(be.Emails), 2, "email queue count mismatch")

	var recipients []string
	for _, email := range be.Emails {
		recipients = append(recipients, email.To[0])
	}
	sort.SliceStable(recipients, func(i, j int) bool {
		r1 := recipients[i]
		r2 := recipients[j]

		return r1 < r2
	})

	assert.DeepEqual(t, recipients, []string{"alice@example.com", "chuck@example.com"}, "email address mismatch")
}
