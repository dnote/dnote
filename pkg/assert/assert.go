/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package assert provides functions to assert a condition in tests
package assert

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(errors.Wrap(err, "reading body"))
		}

		t.Errorf("status code mismatch. %s: got %v want %v. Message was: '%s'", message, res.StatusCode, expected, string(body))
	}
}
