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

// Package client provides interfaces for interacting with the Dnote server
// and the data structures for responses
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

// ErrInvalidLogin is an error for invalid credentials for login
var ErrInvalidLogin = errors.New("wrong credentials")

// ErrContentTypeMismatch is an error for invalid credentials for login
var ErrContentTypeMismatch = errors.New("content type mismatch")

// HTTPError represents an HTTP error response from the server
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf(`response %d "%s"`, e.StatusCode, e.Message)
}

// IsConflict returns true if the error is a 409 Conflict error
func (e *HTTPError) IsConflict() bool {
	return e.StatusCode == 409
}

var contentTypeApplicationJSON = "application/json"
var contentTypeNone = ""

// requestOptions contains options for requests
type requestOptions struct {
	HTTPClient *http.Client
	// ExpectedContentType is the Content-Type that the client is expecting from the server
	ExpectedContentType *string
}

const (
	// clientRateLimitPerSecond is the max requests per second the client will make
	clientRateLimitPerSecond = 50
	// clientRateLimitBurst is the burst capacity for rate limiting
	clientRateLimitBurst = 100
)

// rateLimitedTransport wraps an http.RoundTripper with rate limiting
type rateLimitedTransport struct {
	transport http.RoundTripper
	limiter   *rate.Limiter
}

func (t *rateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Wait for rate limiter to allow the request
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return t.transport.RoundTrip(req)
}

// NewRateLimitedHTTPClient creates an HTTP client with rate limiting
func NewRateLimitedHTTPClient() *http.Client {
	// Calculate interval from rate: 1 second / requests per second
	interval := time.Second / time.Duration(clientRateLimitPerSecond)

	transport := &rateLimitedTransport{
		transport: http.DefaultTransport,
		limiter:   rate.NewLimiter(rate.Every(interval), clientRateLimitBurst),
	}
	return &http.Client{
		Transport: transport,
	}
}

func getHTTPClient(ctx context.DnoteCtx, options *requestOptions) *http.Client {
	if options != nil && options.HTTPClient != nil {
		return options.HTTPClient
	}

	if ctx.HTTPClient != nil {
		return ctx.HTTPClient
	}

	return &http.Client{}
}

func getExpectedContentType(options *requestOptions) string {
	if options != nil && options.ExpectedContentType != nil {
		return *options.ExpectedContentType
	}

	return contentTypeApplicationJSON
}

func getReq(ctx context.DnoteCtx, path, method, body string) (*http.Request, error) {
	endpoint := fmt.Sprintf("%s%s", ctx.APIEndpoint, path)
	req, err := http.NewRequest(method, endpoint, strings.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "constructing http request")
	}

	req.Header.Set("CLI-Version", ctx.Version)

	if ctx.SessionKey != "" {
		credential := fmt.Sprintf("Bearer %s", ctx.SessionKey)
		req.Header.Set("Authorization", credential)
	}

	return req, nil
}

// checkRespErr checks if the given http response indicates an error. It returns a boolean indicating
// if the response is an error, and a decoded error message.
func checkRespErr(res *http.Response) error {
	if res.StatusCode < 400 {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrapf(err, "server responded with %d but client could not read the response body", res.StatusCode)
	}

	bodyStr := string(body)
	return &HTTPError{
		StatusCode: res.StatusCode,
		Message:    strings.TrimRight(bodyStr, "\n"),
	}
}

func checkContentType(res *http.Response, options *requestOptions) error {
	expected := getExpectedContentType(options)

	got := res.Header.Get("Content-Type")
	if got != expected {
		return errors.Wrapf(ErrContentTypeMismatch, "got: '%s' want: '%s'. Did you configure your endpoint correctly?", got, expected)
	}

	return nil
}

// doReq does a http request to the given path in the api endpoint
func doReq(ctx context.DnoteCtx, method, path, body string, options *requestOptions) (*http.Response, error) {
	req, err := getReq(ctx, path, method, body)
	if err != nil {
		return nil, errors.Wrap(err, "getting request")
	}

	log.Debug("HTTP %s %s\n", method, path)

	hc := getHTTPClient(ctx, options)
	res, err := hc.Do(req)
	if err != nil {
		return res, errors.Wrap(err, "making http request")
	}

	log.Debug("HTTP %d %s\n", res.StatusCode, res.Status)

	if err = checkRespErr(res); err != nil {
		return res, errors.Wrap(err, "server responded with an error")
	}

	if err = checkContentType(res, options); err != nil {
		return res, errors.Wrap(err, "unexpected Content-Type")
	}

	return res, nil
}

// doAuthorizedReq does a http request to the given path in the api endpoint as a user,
// with the appropriate headers. The given path should include the preceding slash.
func doAuthorizedReq(ctx context.DnoteCtx, method, path, body string, options *requestOptions) (*http.Response, error) {
	if ctx.SessionKey == "" {
		return nil, errors.New("no session key found")
	}

	return doReq(ctx, method, path, body, options)
}

// GetSyncStateResp is the response get sync state endpoint
type GetSyncStateResp struct {
	FullSyncBefore int   `json:"full_sync_before"`
	MaxUSN         int   `json:"max_usn"`
	CurrentTime    int64 `json:"current_time"`
}

// GetSyncState gets the sync state response from the server
func GetSyncState(ctx context.DnoteCtx) (GetSyncStateResp, error) {
	var ret GetSyncStateResp

	res, err := doAuthorizedReq(ctx, "GET", "/v3/sync/state", "", nil)
	if err != nil {
		return ret, errors.Wrap(err, "constructing http request")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ret, errors.Wrap(err, "reading the response body")
	}

	if err = json.Unmarshal(body, &ret); err != nil {
		return ret, errors.Wrap(err, "unmarshalling the payload")
	}

	return ret, nil
}

// SyncFragNote represents a note in a sync fragment and contains only the necessary information
// for the client to sync the note locally
type SyncFragNote struct {
	UUID      string    `json:"uuid"`
	BookUUID  string    `json:"book_uuid"`
	USN       int       `json:"usn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AddedOn   int64     `json:"added_on"`
	EditedOn  int64     `json:"edited_on"`
	Body      string    `json:"content"`
	Deleted   bool      `json:"deleted"`
}

// SyncFragBook represents a book in a sync fragment and contains only the necessary information
// for the client to sync the note locally
type SyncFragBook struct {
	UUID      string    `json:"uuid"`
	USN       int       `json:"usn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AddedOn   int64     `json:"added_on"`
	Label     string    `json:"label"`
	Deleted   bool      `json:"deleted"`
}

// SyncFragment contains a piece of information about the server's state.
type SyncFragment struct {
	FragMaxUSN    int            `json:"frag_max_usn"`
	UserMaxUSN    int            `json:"user_max_usn"`
	CurrentTime   int64          `json:"current_time"`
	Notes         []SyncFragNote `json:"notes"`
	Books         []SyncFragBook `json:"books"`
	ExpungedNotes []string       `json:"expunged_notes"`
	ExpungedBooks []string       `json:"expunged_books"`
}

// GetSyncFragmentResp is the response from the get sync fragment endpoint
type GetSyncFragmentResp struct {
	Fragment SyncFragment `json:"fragment"`
}

// GetSyncFragment gets a sync fragment response from the server
func GetSyncFragment(ctx context.DnoteCtx, afterUSN int) (GetSyncFragmentResp, error) {
	v := url.Values{}
	v.Set("after_usn", strconv.Itoa(afterUSN))
	queryStr := v.Encode()

	path := fmt.Sprintf("/v3/sync/fragment?%s", queryStr)
	res, err := doAuthorizedReq(ctx, "GET", path, "", nil)
	if err != nil {
		return GetSyncFragmentResp{}, errors.Wrap(err, "making the request")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return GetSyncFragmentResp{}, errors.Wrap(err, "reading the response body")
	}

	var resp GetSyncFragmentResp
	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, errors.Wrap(err, "unmarshalling the payload")
	}

	return resp, nil
}

// RespBook is the book in the response from the create book api
type RespBook struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	USN       int       `json:"usn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Label     string    `json:"label"`
}

// CreateBookPayload is a payload for creating a book
type CreateBookPayload struct {
	Name string `json:"name"`
}

// CreateBookResp is the response from create book api
type CreateBookResp struct {
	Book RespBook `json:"book"`
}

// CreateBook creates a new book in the server
func CreateBook(ctx context.DnoteCtx, label string) (CreateBookResp, error) {
	payload := CreateBookPayload{
		Name: label,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return CreateBookResp{}, errors.Wrap(err, "marshaling payload")
	}

	res, err := doAuthorizedReq(ctx, "POST", "/v3/books", string(b), nil)
	if err != nil {
		return CreateBookResp{}, errors.Wrap(err, "posting a book to the server")
	}

	var resp CreateBookResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return resp, errors.Wrap(err, "decoding response payload")
	}

	return resp, nil
}

type updateBookPayload struct {
	Name *string `json:"name"`
}

// UpdateBookResp is the response from create book api
type UpdateBookResp struct {
	Book RespBook `json:"book"`
}

// UpdateBook updates a book in the server
func UpdateBook(ctx context.DnoteCtx, label, uuid string) (UpdateBookResp, error) {
	payload := updateBookPayload{
		Name: &label,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return UpdateBookResp{}, errors.Wrap(err, "marshaling payload")
	}

	endpoint := fmt.Sprintf("/v3/books/%s", uuid)
	res, err := doAuthorizedReq(ctx, "PATCH", endpoint, string(b), nil)
	if err != nil {
		return UpdateBookResp{}, errors.Wrap(err, "posting a book to the server")
	}

	var resp UpdateBookResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return resp, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

// DeleteBookResp is the response from create book api
type DeleteBookResp struct {
	Status int      `json:"status"`
	Book   RespBook `json:"book"`
}

// DeleteBook deletes a book in the server
func DeleteBook(ctx context.DnoteCtx, uuid string) (DeleteBookResp, error) {
	endpoint := fmt.Sprintf("/v3/books/%s", uuid)
	res, err := doAuthorizedReq(ctx, "DELETE", endpoint, "", nil)
	if err != nil {
		return DeleteBookResp{}, errors.Wrap(err, "deleting a book in the server")
	}

	var resp DeleteBookResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return resp, errors.Wrap(err, "decoding the response")
	}

	return resp, nil
}

// CreateNotePayload is a payload for creating a note
type CreateNotePayload struct {
	BookUUID string `json:"book_uuid"`
	Body     string `json:"content"`
}

// CreateNoteResp is the response from create note endpoint
type CreateNoteResp struct {
	Result RespNote `json:"result"`
}

type respNoteBook struct {
	UUID  string `json:"uuid"`
	Label string `json:"label"`
}

type respNoteUser struct {
	Name string `json:"name"`
}

// RespNote is a note in the response
type RespNote struct {
	UUID      string       `json:"uuid"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Body      string       `json:"content"`
	AddedOn   int64        `json:"added_on"`
	Public    bool         `json:"public"`
	USN       int          `json:"usn"`
	Book      respNoteBook `json:"book"`
	User      respNoteUser `json:"user"`
}

// CreateNote creates a note in the server
func CreateNote(ctx context.DnoteCtx, bookUUID, content string) (CreateNoteResp, error) {
	payload := CreateNotePayload{
		BookUUID: bookUUID,
		Body:     content,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return CreateNoteResp{}, errors.Wrap(err, "marshaling payload")
	}

	res, err := doAuthorizedReq(ctx, "POST", "/v3/notes", string(b), nil)
	if err != nil {
		return CreateNoteResp{}, errors.Wrap(err, "posting a book to the server")
	}

	var resp CreateNoteResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return CreateNoteResp{}, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

type updateNotePayload struct {
	BookUUID *string `json:"book_uuid"`
	Body     *string `json:"content"`
}

// UpdateNoteResp is the response from create book api
type UpdateNoteResp struct {
	Status int      `json:"status"`
	Result RespNote `json:"result"`
}

// UpdateNote updates a note in the server
func UpdateNote(ctx context.DnoteCtx, uuid, bookUUID, content string) (UpdateNoteResp, error) {
	payload := updateNotePayload{
		BookUUID: &bookUUID,
		Body:     &content,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return UpdateNoteResp{}, errors.Wrap(err, "marshaling payload")
	}

	endpoint := fmt.Sprintf("/v3/notes/%s", uuid)
	res, err := doAuthorizedReq(ctx, "PATCH", endpoint, string(b), nil)
	if err != nil {
		return UpdateNoteResp{}, errors.Wrap(err, "patching a note to the server")
	}

	var resp UpdateNoteResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return UpdateNoteResp{}, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

// DeleteNoteResp is the response from remove note api
type DeleteNoteResp struct {
	Status int      `json:"status"`
	Result RespNote `json:"result"`
}

// DeleteNote removes a note in the server
func DeleteNote(ctx context.DnoteCtx, uuid string) (DeleteNoteResp, error) {
	endpoint := fmt.Sprintf("/v3/notes/%s", uuid)
	res, err := doAuthorizedReq(ctx, "DELETE", endpoint, "", nil)
	if err != nil {
		return DeleteNoteResp{}, errors.Wrap(err, "patching a note to the server")
	}

	var resp DeleteNoteResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return DeleteNoteResp{}, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

// GetBooksResp is a response from get books endpoint
type GetBooksResp []struct {
	UUID  string `json:"uuid"`
	Label string `json:"label"`
}

// GetBooks gets books from the server
func GetBooks(ctx context.DnoteCtx, sessionKey string) (GetBooksResp, error) {
	res, err := doAuthorizedReq(ctx, "GET", "/v3/books", "", nil)
	if err != nil {
		return GetBooksResp{}, errors.Wrap(err, "making http request")
	}

	var resp GetBooksResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return GetBooksResp{}, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

// PresigninResponse is a reponse from /v3/presignin endpoint
type PresigninResponse struct {
	Iteration int `json:"iteration"`
}

// GetPresignin gets presignin credentials
func GetPresignin(ctx context.DnoteCtx, email string) (PresigninResponse, error) {
	res, err := doReq(ctx, "GET", fmt.Sprintf("/v3/presignin?email=%s", email), "", nil)
	if err != nil {
		return PresigninResponse{}, errors.Wrap(err, "making http request")
	}

	var resp PresigninResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return PresigninResponse{}, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

// SigninPayload is a payload for /v3/signin
type SigninPayload struct {
	Email    string `json:"email"`
	Passowrd string `json:"password"`
}

// SigninResponse is a response from /v3/signin endpoint
type SigninResponse struct {
	Key       string `json:"key"`
	ExpiresAt int64  `json:"expires_at"`
}

// Signin requests a session token
func Signin(ctx context.DnoteCtx, email, password string) (SigninResponse, error) {
	payload := SigninPayload{
		Email:    email,
		Passowrd: password,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return SigninResponse{}, errors.Wrap(err, "marshaling payload")
	}
	res, err := doReq(ctx, "POST", "/v3/signin", string(b), nil)
	if err != nil {
		// Check if this is a 401 Unauthorized error
		var httpErr *HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusUnauthorized {
			return SigninResponse{}, ErrInvalidLogin
		}
		return SigninResponse{}, errors.Wrap(err, "making http request")
	}

	var resp SigninResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return SigninResponse{}, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

// Signout deletes a user session on the server side
func Signout(ctx context.DnoteCtx, sessionKey string) error {
	// Create a client that shares the transport (and thus rate limiter) from ctx.HTTPClient
	// but doesn't follow redirects
	var hc *http.Client
	if ctx.HTTPClient != nil {
		hc = &http.Client{
			Transport: ctx.HTTPClient.Transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	} else {
		log.Warnf("No HTTP client configured for signout - falling back\n")
		hc = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	opts := requestOptions{
		HTTPClient:          hc,
		ExpectedContentType: &contentTypeNone,
	}
	_, err := doAuthorizedReq(ctx, "POST", "/v3/signout", "", &opts)
	if err != nil {
		return errors.Wrap(err, "making http request")
	}

	return nil
}
