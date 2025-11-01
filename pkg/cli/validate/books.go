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

package validate

import (
	"strings"

	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
)

var reservedBookNames = []string{"trash", "conflicts"}

// ErrBookNameReserved is an error incidating that the specified book name is reserved
var ErrBookNameReserved = errors.New("The book name is reserved")

// ErrBookNameNumeric is an error for a book name that only contains numbers
var ErrBookNameNumeric = errors.New("The book name cannot contain only numbers")

// ErrBookNameHasSpace is an error for a book name that has any space
var ErrBookNameHasSpace = errors.New("The book name cannot contain spaces")

// ErrBookNameEmpty is an error for an empty book name
var ErrBookNameEmpty = errors.New("The book name is empty")

// ErrBookNameMultiline is an error for a book name that has linebreaks
var ErrBookNameMultiline = errors.New("The book name contains multiple lines")

func isReservedName(name string) bool {
	for _, n := range reservedBookNames {
		if name == n {
			return true
		}
	}

	return false
}

// BookName validates a book name
func BookName(name string) error {
	if name == "" {
		return ErrBookNameEmpty
	}

	if isReservedName(name) {
		return ErrBookNameReserved
	}

	if utils.IsNumber(name) {
		return ErrBookNameNumeric
	}

	if strings.Contains(name, " ") {
		return ErrBookNameHasSpace
	}

	if strings.Contains(name, "\n") || strings.Contains(name, "\r\n") {
		return ErrBookNameMultiline
	}

	return nil
}
