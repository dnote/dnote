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

package models

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/testutils"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// TestDB is the database connection to a test database
var TestDB *gorm.DB

// InitTestDB establishes connection pool with the test database specified by
// the environment variable configuration and initalizes a new schema
func InitTestDB() {
	c := config.Load()
	fmt.Println(c.DB.GetConnectionStr())
	db := Open(c)

	InitSchema(db)

	TestDB = db
}

// ClearTestData deletes all records from the database
func ClearTestData(db *gorm.DB) {
	if err := db.Delete(&Book{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear books"))
	}
	if err := db.Delete(&Note{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear notes"))
	}
	if err := db.Delete(&Notification{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear notifications"))
	}
	if err := db.Delete(&User{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear users"))
	}
	if err := db.Delete(&Account{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear accounts"))
	}
	if err := db.Delete(&Token{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear tokens"))
	}
	if err := db.Delete(&EmailPreference{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear email preferences"))
	}
	if err := db.Delete(&Session{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear sessions"))
	}
}

// SetUpUserData creates and returns a new user for testing purposes
func SetUpUserData() User {
	user := User{
		Cloud: true,
	}

	if err := TestDB.Save(&user).Error; err != nil {
		panic(errors.Wrap(err, "Failed to prepare user"))
	}

	return user
}

// SetUpAccountData creates and returns a new account for the user
func SetUpAccountData(user User, email, password string) Account {
	account := Account{
		UserID: user.ID,
	}
	if email != "" {
		account.Email = ToNullString(email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(errors.Wrap(err, "Failed to hash password"))
	}
	account.Password = ToNullString(string(hashedPassword))

	if err := TestDB.Save(&account).Error; err != nil {
		panic(errors.Wrap(err, "Failed to prepare account"))
	}

	return account
}

// SetupSession creates and returns a new user session
func SetupSession(t *testing.T, user User) Session {
	session := Session{
		Key:       "Vvgm3eBXfXGEFWERI7faiRJ3DAzJw+7DdT9J1LEyNfI=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	if err := TestDB.Save(&session).Error; err != nil {
		t.Fatal(errors.Wrap(err, "Failed to prepare user"))
	}

	return session
}

// SetupEmailPreferenceData creates and returns a new email frequency for a user
func SetupEmailPreferenceData(user User, inactiveReminder bool) EmailPreference {
	frequency := EmailPreference{
		UserID:           user.ID,
		InactiveReminder: inactiveReminder,
	}

	if err := TestDB.Save(&frequency).Error; err != nil {
		panic(errors.Wrap(err, "Failed to prepare email frequency"))
	}

	return frequency
}

// SetReqAuthHeader sets the authorization header in the given request for the given user
func SetReqAuthHeader(t *testing.T, req *http.Request, user User) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(errors.Wrap(err, "reading random bits"))
	}

	session := Session{
		Key:       base64.StdEncoding.EncodeToString(b),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
	}
	if err := TestDB.Save(&session).Error; err != nil {
		t.Fatal(errors.Wrap(err, "Failed to prepare user"))
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Key))
}

// HTTPAuthDo makes an HTTP request with an appropriate authorization header for a user
func HTTPAuthDo(t *testing.T, req *http.Request, user User) *http.Response {
	SetReqAuthHeader(t, req, user)

	return testutils.HTTPDo(t, req)
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
