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

// Package assert provides functions to assert a condition in tests
package assert

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime/debug"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

func getErrorMessage(m string, a, b interface{}) string {
	return fmt.Sprintf(`%s.
Actual:
========================
%+v
========================

Expected:
========================
%+v
========================

%s`, m, a, b, string(debug.Stack()))
}

func checkEqual(a, b interface{}, message string) (bool, string) {
	if a == b {
		return true, ""
	}

	var m string
	if len(message) == 0 {
		m = fmt.Sprintf("%v != %v", a, b)
	} else {
		m = message
	}
	errorMessage := getErrorMessage(m, a, b)

	return false, errorMessage
}

// Equal errors a test if the actual does not match the expected
func Equal(t *testing.T, a, b interface{}, message string) {
	ok, m := checkEqual(a, b, message)
	if !ok {
		t.Error(m)
	}
}

// Equalf fails a test if the actual does not match the expected
func Equalf(t *testing.T, a, b interface{}, message string) {
	ok, m := checkEqual(a, b, message)
	if !ok {
		t.Fatal(m)
	}
}

// NotEqual fails a test if the actual matches the expected
func NotEqual(t *testing.T, a, b interface{}, message string) {
	ok, m := checkEqual(a, b, message)
	if ok {
		t.Error(m)
	}
}

// NotEqualf fails a test if the actual matches the expected
func NotEqualf(t *testing.T, a, b interface{}, message string) {
	ok, m := checkEqual(a, b, message)
	if ok {
		t.Fatal(m)
	}
}

// DeepEqual fails a test if the actual does not deeply equal the expected
func DeepEqual(t *testing.T, a, b interface{}, message string) {
	if cmp.Equal(a, b) {
		return
	}

	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}

	errorMessage := getErrorMessage(message, a, b)
	errorMessage = fmt.Sprintf("%v\n%v", errorMessage, cmp.Diff(a, b))
	t.Error(errorMessage)
}

// EqualJSON asserts that two JSON strings are equal
func EqualJSON(t *testing.T, a, b, message string) {
	var o1 interface{}
	var o2 interface{}

	err := json.Unmarshal([]byte(a), &o1)
	if err != nil {
		panic(fmt.Errorf("Error mashalling string 1 :: %s", err.Error()))
	}
	err = json.Unmarshal([]byte(b), &o2)
	if err != nil {
		panic(fmt.Errorf("Error mashalling string 2 :: %s", err.Error()))
	}

	if reflect.DeepEqual(o1, o2) {
		return
	}

	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Errorf("%s.\nActual:   %+v.\nExpected: %+v.", message, a, b)
}

// StatusCodeEquals asserts that the reponse's status code is equal to the
// expected
func StatusCodeEquals(t *testing.T, res *http.Response, expected int, message string) {
	if res.StatusCode != expected {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(errors.Wrap(err, "reading body"))
		}

		t.Errorf("status code mismatch. %s: got %v want %v. Message was: '%s'", message, res.StatusCode, expected, string(body))
	}
}
