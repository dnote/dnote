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
