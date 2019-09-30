/* Copyright (C) 2019 Monomax Software Pty Ltd
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
