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

// Package diff provides line-by-line diff feature by wrapping
// a package github.com/sergi/go-diff/diffmatchpatch
package diff

import (
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	// DiffEqual represents an equal diff
	DiffEqual = diffmatchpatch.DiffEqual
	// DiffInsert represents an insert diff
	DiffInsert = diffmatchpatch.DiffInsert
	// DiffDelete represents a delete diff
	DiffDelete = diffmatchpatch.DiffDelete
)

// Do computes line-by-line diff between two strings
func Do(s1, s2 string) (diffs []diffmatchpatch.Diff) {
	dmp := diffmatchpatch.New()
	dmp.DiffTimeout = time.Hour

	s1Chars, s2Chars, arr := dmp.DiffLinesToRunes(s1, s2)
	diffs = dmp.DiffMainRunes(s1Chars, s2Chars, false)
	diffs = dmp.DiffCharsToLines(diffs, arr)

	return diffs
}
