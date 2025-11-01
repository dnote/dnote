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
