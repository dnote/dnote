package remind

import (
	"os"
	"sort"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/pkg/errors"
)

func getTestContext(c clock.Clock, be *testutils.MockEmailbackendImplementation) Context {
	emailTmplDir := os.Getenv("DNOTE_TEST_EMAIL_TEMPLATE_DIR")

	con := Context{
		DB:           testutils.DB,
		Clock:        c,
		EmailTmpl:    mailer.NewTemplates(&emailTmplDir),
		EmailBackend: be,
	}

	return con
}

func TestDoInactive(t *testing.T) {
	defer testutils.ClearData()

	t1 := time.Now()

	// u1 is an active user
	u1 := testutils.SetupUserData()
	a1 := testutils.SetupAccountData(u1, "alice@example.com", "pass1234")
	testutils.MustExec(t, testutils.DB.Model(&a1).Update("email_verified", true), "setting email verified")
	testutils.MustExec(t, testutils.DB.Save(&database.EmailPreference{UserID: u1.ID, InactiveReminder: true}), "preparing email preference")

	b1 := database.Book{
		UserID: u1.ID,
		Label:  "js",
	}
	testutils.MustExec(t, testutils.DB.Save(&b1), "preparing b1")
	n1 := database.Note{
		BookUUID: b1.UUID,
		UserID:   u1.ID,
	}
	testutils.MustExec(t, testutils.DB.Save(&n1), "preparing n1")

	// u2 is an inactive user
	u2 := testutils.SetupUserData()
	a2 := testutils.SetupAccountData(u2, "bob@example.com", "pass1234")
	testutils.MustExec(t, testutils.DB.Model(&a2).Update("email_verified", true), "setting email verified")
	testutils.MustExec(t, testutils.DB.Save(&database.EmailPreference{UserID: u2.ID, InactiveReminder: true}), "preparing email preference")

	b2 := database.Book{
		UserID: u2.ID,
		Label:  "css",
	}
	testutils.MustExec(t, testutils.DB.Save(&b2), "preparing b2")
	n2 := database.Note{
		UserID:   u2.ID,
		BookUUID: b2.UUID,
	}
	testutils.MustExec(t, testutils.DB.Save(&n2), "preparing n2")
	testutils.MustExec(t, testutils.DB.Model(&n2).Update("created_at", t1.AddDate(0, 0, -15)), "preparing n2")

	// u3 is an inactive user with inactive alert email preference disabled
	u3 := testutils.SetupUserData()
	a3 := testutils.SetupAccountData(u3, "alice@example.com", "pass1234")
	testutils.MustExec(t, testutils.DB.Model(&a3).Update("email_verified", true), "setting email verified")
	testutils.MustExec(t, testutils.DB.Save(&database.EmailPreference{UserID: u3.ID, InactiveReminder: false}), "preparing email preference")

	b3 := database.Book{
		UserID: u3.ID,
		Label:  "js",
	}
	testutils.MustExec(t, testutils.DB.Save(&b3), "preparing b3")
	n3 := database.Note{
		BookUUID: b3.UUID,
		UserID:   u3.ID,
	}
	testutils.MustExec(t, testutils.DB.Save(&n3), "preparing n3")
	testutils.MustExec(t, testutils.DB.Model(&n3).Update("created_at", t1.AddDate(0, 0, -15)), "preparing n3")

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
	defer testutils.ClearData()

	// setup sets up an inactive user
	setup := func(t *testing.T, now time.Time, email string) database.User {
		u := testutils.SetupUserData()
		a := testutils.SetupAccountData(u, email, "pass1234")
		testutils.MustExec(t, testutils.DB.Model(&a).Update("email_verified", true), "setting email verified")
		testutils.MustExec(t, testutils.DB.Save(&database.EmailPreference{UserID: u.ID, InactiveReminder: true}), "preparing email preference")

		b := database.Book{
			UserID: u.ID,
			Label:  "css",
		}
		testutils.MustExec(t, testutils.DB.Save(&b), "preparing book")
		n := database.Note{
			UserID:   u.ID,
			BookUUID: b.UUID,
		}
		testutils.MustExec(t, testutils.DB.Save(&n), "preparing note")
		testutils.MustExec(t, testutils.DB.Model(&n).Update("created_at", now.AddDate(0, 0, -15)), "preparing note")

		return u
	}

	// Set up
	now := time.Now()

	setup(t, now, "alice@example.com")

	bob := setup(t, now, "bob@example.com")
	bobNotif := database.Notification{Type: "inactive", UserID: bob.ID}
	testutils.MustExec(t, testutils.DB.Create(&bobNotif), "preparing inactive notification for bob")
	testutils.MustExec(t, testutils.DB.Model(&bobNotif).Update("created_at", now.AddDate(0, 0, -7)), "preparing created_at for inactive notification for bob")

	chuck := setup(t, now, "chuck@example.com")
	chuckNotif := database.Notification{Type: "inactive", UserID: chuck.ID}
	testutils.MustExec(t, testutils.DB.Create(&chuckNotif), "preparing inactive notification for chuck")
	testutils.MustExec(t, testutils.DB.Model(&chuckNotif).Update("created_at", now.AddDate(0, 0, -15)), "preparing created_at for inactive notification for chuck")

	dan := setup(t, now, "dan@example.com")
	danNotif1 := database.Notification{Type: "inactive", UserID: dan.ID}
	testutils.MustExec(t, testutils.DB.Create(&danNotif1), "preparing inactive notification 1 for dan")
	testutils.MustExec(t, testutils.DB.Model(&danNotif1).Update("created_at", now.AddDate(0, 0, -10)), "preparing created_at for inactive notification for dan")
	danNotif2 := database.Notification{Type: "inactive", UserID: dan.ID}
	testutils.MustExec(t, testutils.DB.Create(&danNotif2), "preparing inactive notification 2 for dan")
	testutils.MustExec(t, testutils.DB.Model(&danNotif2).Update("created_at", now.AddDate(0, 0, -15)), "preparing created_at for inactive notification for dan")

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
