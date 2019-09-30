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

package find

import (
	"regexp"
)

var newLineReg = regexp.MustCompile(`\r?\n`)

const (
	// tokenKindChar represents utf-8 character
	tokenKindChar = iota
	// tokenKindHLBegin represents a beginning of a highlighted section
	tokenKindHLBegin
	// tokenKindHLEnd represents an end of a highlighted section
	tokenKindHLEnd
	// tokenKindEOL represents an end of line
	tokenKindEOL
)

type token struct {
	Value byte
	Kind  int
}

// getNextIdx validates that the given index is within the range of the given string.
// If so, it returns the given index. Otherwise it returns -1.
func getNextIdx(candidate int, s string) int {
	if candidate <= len(s)-1 {
		return candidate
	}

	return -1
}

// scanToken scans the given string for a token at the given index. It returns
// a token and the next index to look for a token. If the given string is exhausted,
// the next index will be -1.
func scanToken(idx int, s string) (token, int) {
	if s[idx] == '<' {
		if len(s)-idx >= 9 {
			lookahead := 9
			candidate := s[idx : idx+lookahead]

			if candidate == "<dnotehl>" {
				nextIdx := getNextIdx(idx+lookahead, s)
				return token{Kind: tokenKindHLBegin}, nextIdx
			}
		}

		if len(s)-idx >= 10 {
			lookahead := 10
			candidate := s[idx : idx+lookahead]

			if candidate == "</dnotehl>" {
				nextIdx := getNextIdx(idx+lookahead, s)
				return token{Kind: tokenKindHLEnd}, nextIdx
			}
		}
	}

	nextIdx := getNextIdx(idx+1, s)

	return token{Value: s[idx], Kind: tokenKindChar}, nextIdx
}

// tokenize lexically analyzes the given matched snippet from a full text search
// and builds a slice of tokens
func tokenize(s string) []token {
	var ret []token

	idx := 0
	for idx != -1 {
		var tok token
		tok, idx = scanToken(idx, s)

		ret = append(ret, tok)
	}

	ret = append(ret, token{Kind: tokenKindEOL})

	return ret
}
