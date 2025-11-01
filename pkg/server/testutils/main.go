/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package testutils provides utilities used in tests
package testutils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB opens a database at the given path and initializes the schema
func InitDB(dbPath string) *gorm.DB {
	db := database.Open(dbPath)
	database.InitSchema(db)
	database.Migrate(db)
	return db
}

// InitMemoryDB creates an in-memory SQLite database with the schema initialized
func InitMemoryDB(t *testing.T) *gorm.DB {
	// Use file-based in-memory database with unique UUID per test to avoid sharing
	uuid, err := helpers.GenUUID()
	if err != nil {
		t.Fatalf("failed to generate UUID for test database: %v", err)
	}
	dbName := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid)
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	database.InitSchema(db)
	database.Migrate(db)

	return db
}

// MustUUID generates a UUID and fails the test on error
func MustUUID(t *testing.T) string {
	uuid, err := helpers.GenUUID()
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to generate UUID"))
	}
	return uuid
}

// SetupUserData creates and returns a new user with email and password for testing purposes
func SetupUserData(db *gorm.DB, email, password string) database.User {
	uuid, err := helpers.GenUUID()
	if err != nil {
		panic(errors.Wrap(err, "Failed to generate UUID"))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(errors.Wrap(err, "Failed to hash password"))
	}

	user := database.User{
		UUID:     uuid,
		Email:    database.ToNullString(email),
		Password: database.ToNullString(string(hashedPassword)),
	}

	if err := db.Save(&user).Error; err != nil {
		panic(errors.Wrap(err, "Failed to prepare user"))
	}

	return user
}

// SetupSession creates and returns a new user session
func SetupSession(db *gorm.DB, user database.User) database.Session {
	session := database.Session{
		Key:       "Vvgm3eBXfXGEFWERI7faiRJ3DAzJw+7DdT9J1LEyNfI=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	if err := db.Save(&session).Error; err != nil {
		panic(errors.Wrap(err, "Failed to prepare user"))
	}

	return session
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

// SetReqAuthHeader sets the authorization header in the given request for the given user with a specific DB
func SetReqAuthHeader(t *testing.T, db *gorm.DB, req *http.Request, user database.User) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(errors.Wrap(err, "reading random bits"))
	}

	session := database.Session{
		Key:       base64.StdEncoding.EncodeToString(b),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
	}
	if err := db.Save(&session).Error; err != nil {
		t.Fatal(errors.Wrap(err, "Failed to prepare user"))
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Key))
}

// HTTPAuthDo makes an HTTP request with an appropriate authorization header for a user with a specific DB
func HTTPAuthDo(t *testing.T, db *gorm.DB, req *http.Request, user database.User) *http.Response {
	SetReqAuthHeader(t, db, req, user)

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

// MakeFormReq makes an HTTP request and returns a response
func MakeFormReq(endpoint, method, path string, data url.Values) *http.Request {
	req := MakeReq(endpoint, method, path, data.Encode())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

// EndpointType is the type of endpoint to be tested
type EndpointType int

const (
	// EndpointWeb represents a web endpoint returning HTML
	EndpointWeb EndpointType = iota
	// EndpointAPI represents an API endpoint returning JSON
	EndpointAPI
)

type endpointTest func(t *testing.T, target EndpointType)

// RunForWebAndAPI runs the given test function for web and API
func RunForWebAndAPI(t *testing.T, name string, runTest endpointTest) {
	t.Run(fmt.Sprintf("%s-web", name), func(t *testing.T) {
		runTest(t, EndpointWeb)
	})

	t.Run(fmt.Sprintf("%s-api", name), func(t *testing.T) {
		runTest(t, EndpointAPI)
	})
}

// PayloadWrapper is a wrapper for a payload that can be converted to
// either URL form values or JSON
type PayloadWrapper struct {
	Data interface{}
}

func (p PayloadWrapper) ToURLValues() url.Values {
	values := url.Values{}

	el := reflect.ValueOf(p.Data)
	if el.Kind() == reflect.Ptr {
		el = el.Elem()
	}
	iVal := el
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		fi := typ.Field(i)
		name := fi.Tag.Get("schema")
		if name == "" {
			name = fi.Name
		}

		if !iVal.Field(i).IsNil() {
			values.Set(name, fmt.Sprint(iVal.Field(i).Elem()))
		}
	}

	return values
}

func (p PayloadWrapper) ToJSON(t *testing.T) string {
	b, err := json.Marshal(p.Data)
	if err != nil {
		t.Fatal(err)
	}

	return string(b)
}

// TrueVal is a true value
var TrueVal = true

// FalseVal is a false value
var FalseVal = false
