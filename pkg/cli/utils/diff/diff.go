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
