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

// Package testutils provides utilities used in tests
package testutils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/dbconn"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// DB is the database connection to a test database
var DB *gorm.DB

// InitTestDB establishes connection pool with the test database specified by
// the environment variable configuration and initalizes a new schema
func InitTestDB() {
	db := dbconn.Open(dbconn.Config{
		Host:     os.Getenv("DBHost"),
		Port:     os.Getenv("DBPort"),
		Name:     os.Getenv("DBName"),
		User:     os.Getenv("DBUser"),
		Password: os.Getenv("DBPassword"),
	})
	database.InitSchema(db)

	DB = db
}

// ClearData deletes all records from the database
func ClearData() {
	if err := DB.Delete(&database.Book{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear books"))
	}
	if err := DB.Delete(&database.Note{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear notes"))
	}
	if err := DB.Delete(&database.Notification{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear notifications"))
	}
	if err := DB.Delete(&database.User{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear users"))
	}
	if err := DB.Delete(&database.Account{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear accounts"))
	}
	if err := DB.Delete(&database.Token{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear tokens"))
	}
	if err := DB.Delete(&database.EmailPreference{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear email preferences"))
	}
	if err := DB.Delete(&database.Session{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear sessions"))
	}
	if err := DB.Delete(&database.Digest{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear digests"))
	}
	if err := DB.Delete(&database.RepetitionRule{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear repetition rules"))
	}
}

// SetupUserData creates and returns a new user for testing purposes
func SetupUserData() database.User {
	user := database.User{
		APIKey: "test-api-key",
		Name:   "user-name",
		Cloud:  true,
	}

	if err := DB.Save(&user).Error; err != nil {
		panic(errors.Wrap(err, "Failed to prepare user"))
	}

	return user
}

// SetupAccountData creates and returns a new account for the user
func SetupAccountData(user database.User, email, password string) database.Account {
	account := database.Account{
		UserID: user.ID,
	}
	if email != "" {
		account.Email = database.ToNullString(email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(errors.Wrap(err, "Failed to hash password"))
	}
	account.Password = database.ToNullString(string(hashedPassword))

	if err := DB.Save(&account).Error; err != nil {
		panic(errors.Wrap(err, "Failed to prepare account"))
	}

	return account
}

// SetupClassicAccountData creates and returns a new account for the user
func SetupClassicAccountData(user database.User, email string) database.Account {
	// email: alice@example.com
	// password: pass1234
	// masterKey: WbUvagj9O6o1Z+4+7COjo7Uqm4MD2QE9EWFXne8+U+8=
	// authKey: /XCYisXJ6/o+vf6NUEtmrdYzJYPz+T9oAUCtMpOjhzc=
	account := database.Account{
		UserID:             user.ID,
		Salt:               "Et0joOigYjdgHBKMN/ijxg==",
		AuthKeyHash:        "SeN3PMz4H/7q9lINB+VPKpygexAuK68wO8pDAgQ4OOQ=",
		CipherKeyEnc:       "f7aFFCh7YS1WlHEOxAmDfs8rUQQoX5tr8AB7ZJQaTYCEM8NhAZCbQTsjFgKOf5iPQhhkm8eDAgPNTuhO",
		ClientKDFIteration: 100000,
		ServerKDFIteration: 100000,
	}
	if email != "" {
		account.Email = database.ToNullString(email)
	}

	if err := DB.Save(&account).Error; err != nil {
		panic(errors.Wrap(err, "Failed to prepare account"))
	}

	return account
}

// SetupSession creates and returns a new user session
func SetupSession(t *testing.T, user database.User) database.Session {
	session := database.Session{
		Key:       "Vvgm3eBXfXGEFWERI7faiRJ3DAzJw+7DdT9J1LEyNfI=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	if err := DB.Save(&session).Error; err != nil {
		t.Fatal(errors.Wrap(err, "Failed to prepare user"))
	}

	return session
}

// SetupEmailPreferenceData creates and returns a new email frequency for a user
func SetupEmailPreferenceData(user database.User, inactiveReminder bool) database.EmailPreference {
	frequency := database.EmailPreference{
		UserID:           user.ID,
		InactiveReminder: inactiveReminder,
	}

	if err := DB.Save(&frequency).Error; err != nil {
		panic(errors.Wrap(err, "Failed to prepare email frequency"))
	}

	return frequency
}

// HTTPDo makes an HTTP request and returns a response
func HTTPDo(t *testing.T, req *http.Request) *http.Response {
	hc := http.Client{
		// Do not follow redirects.
		// e.g. /logout redirects to a page but we'd like to test the redirect
		// itself, not what happens after the redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := hc.Do(req)
	if err != nil {
		t.Fatal(errors.Wrap(err, "performing http request"))
	}

	return res
}

// HTTPAuthDo makes an HTTP request with an appropriate authorization header for a user
func HTTPAuthDo(t *testing.T, req *http.Request, user database.User) *http.Response {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(errors.Wrap(err, "reading random bits"))
	}

	session := database.Session{
		Key:       base64.StdEncoding.EncodeToString(b),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
	}
	if err := DB.Save(&session).Error; err != nil {
		t.Fatal(errors.Wrap(err, "Failed to prepare user"))
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Key))

	return HTTPDo(t, req)

}

// MakeReq makes an HTTP request and returns a response
func MakeReq(endpoint string, method, path, data string) *http.Request {
	u := fmt.Sprintf("%s%s", endpoint, path)

	req, err := http.NewRequest(method, u, strings.NewReader(data))
	if err != nil {
		panic(errors.Wrap(err, "constructing http request"))
	}

	return req
}

// MustExec fails the test if the given database query has error
func MustExec(t *testing.T, db *gorm.DB, message string) {
	if err := db.Error; err != nil {
		t.Fatalf("%s: %s", message, err.Error())
	}
}

// GetCookieByName returns a cookie with the given name
func GetCookieByName(cookies []*http.Cookie, name string) *http.Cookie {
	var ret *http.Cookie

	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == name {
			ret = cookies[i]
			break
		}
	}

	return ret
}

// CreateMockStripeBackend returns a mock stripe backend that uses
// the given test server
func CreateMockStripeBackend(ts *httptest.Server) stripe.Backend {
	stripeMockBackend := stripe.GetBackendWithConfig(
		stripe.APIBackend,
		&stripe.BackendConfig{
			URL:        ts.URL,
			HTTPClient: ts.Client(),
		},
	)

	return stripeMockBackend
}

// MustRespondJSON responds with the JSON-encoding of the given interface. If the encoding
// fails, the test fails. It is used by test servers.
func MustRespondJSON(t *testing.T, w http.ResponseWriter, i interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(i); err != nil {
		t.Fatal(message)
	}
}

// MockEmail is a mock email data
type MockEmail struct {
	Subject string
	From    string
	To      []string
	Body    string
}

// MockEmailbackendImplementation is an email backend that simply discards the emails
type MockEmailbackendImplementation struct {
	mu     sync.RWMutex
	Emails []MockEmail
}

// Clear clears the mock email queue
func (b *MockEmailbackendImplementation) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Emails = []MockEmail{}
}

// Queue is an implementation of Backend.Queue.
func (b *MockEmailbackendImplementation) Queue(subject, from string, to []string, contentType, body string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Emails = append(b.Emails, MockEmail{
		Subject: subject,
		From:    from,
		To:      to,
		Body:    body,
	})

	return nil
}
