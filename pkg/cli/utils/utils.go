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

package utils

import (
	"regexp"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// GenerateUUID returns a uuid v4 in string
func GenerateUUID() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "generating uuid")
	}

	return u.String(), nil
}

// regexNumber is a regex that matches a string that looks like an integer
var regexNumber = regexp.MustCompile(`^\d+$`)

// IsNumber checks if the given string is in the form of a number
func IsNumber(s string) bool {
	if s == "" {
		return false
	}

	return regexNumber.MatchString(s)
}
