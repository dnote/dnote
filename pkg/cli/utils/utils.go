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

package utils

import (
	"regexp"

	"github.com/satori/go.uuid"
)

// GenerateUUID returns a uid
func GenerateUUID() string {
	return uuid.NewV4().String()
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
