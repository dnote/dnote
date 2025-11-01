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
	"fmt"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

func TestScanToken(t *testing.T) {
	testCases := []struct {
		input  string
		idx    int
		retTok token
		retIdx int
	}{
		{
			input:  "foo bar",
			idx:    1,
			retTok: token{Value: 'o', Kind: tokenKindChar},
			retIdx: 2,
		},
		{
			input:  "foo bar",
			idx:    6,
			retTok: token{Value: 'r', Kind: tokenKindChar},
			retIdx: -1,
		},
		{
			input:  "foo <bar>",
			idx:    4,
			retTok: token{Value: '<', Kind: tokenKindChar},
			retIdx: 5,
		},
		{
			input:  "foo <dnotehL>",
			idx:    4,
			retTok: token{Value: '<', Kind: tokenKindChar},
			retIdx: 5,
		},
		{
			input:  "foo <dnotehl>bar</dnotehl> foo bar",
			idx:    4,
			retTok: token{Kind: tokenKindHLBegin},
			retIdx: 13,
		},
		{
			input:  "foo <dnotehl>bar</dnotehl> <dnotehl>foo</dnotehl> bar",
			idx:    4,
			retTok: token{Kind: tokenKindHLBegin},
			retIdx: 13,
		},
		{
			input:  "foo <dnotehl>bar</dnotehl> <dnotehl>foo</dnotehl> bar",
			idx:    27,
			retTok: token{Kind: tokenKindHLBegin},
			retIdx: 36,
		},
		{
			input:  "foo <dnotehl>bar</dnotehl> foo bar",
			idx:    13,
			retTok: token{Value: 'b', Kind: tokenKindChar},
			retIdx: 14,
		},
		{
			input:  "foo <dnotehl>bar</dnotehl> foo bar",
			idx:    16,
			retTok: token{Kind: tokenKindHLEnd},
			retIdx: 26,
		},
		{
			input:  "<dno<dnotehl>tehl>",
			idx:    0,
			retTok: token{Value: '<', Kind: tokenKindChar},
			retIdx: 1,
		},
		{
			input:  "<dno<dnotehl>tehl>",
			idx:    4,
			retTok: token{Kind: tokenKindHLBegin},
			retIdx: 13,
		},
		{
			input:  "foo <dnotehl>bar</dnotehl>",
			idx:    16,
			retTok: token{Kind: tokenKindHLEnd},
			retIdx: -1,
		},
		// user writes reserved token
		{
			input:  "foo <dnotehl>",
			idx:    4,
			retTok: token{Kind: tokenKindHLBegin},
			retIdx: -1,
		},
	}

	for tcIdx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", tcIdx), func(t *testing.T) {
			tok, nextIdx := scanToken(tc.idx, tc.input)

			assert.Equal(t, nextIdx, tc.retIdx, "retIdx mismatch")
			assert.DeepEqual(t, tok, tc.retTok, "retTok mismatch")
		})
	}
}

func TestTokenize(t *testing.T) {
	testCases := []struct {
		input  string
		tokens []token
	}{
		{
			input: "ab<dnotehl>c</dnotehl>",
			tokens: []token{
				token{
					Kind:  tokenKindChar,
					Value: 'a',
				},
				token{
					Kind:  tokenKindChar,
					Value: 'b',
				},
				token{
					Kind: tokenKindHLBegin,
				},
				token{
					Kind:  tokenKindChar,
					Value: 'c',
				},
				token{
					Kind: tokenKindHLEnd,
				},
				token{
					Kind: tokenKindEOL,
				},
			},
		},
		{
			input: "ab<dnotehl>c</dnotehl>d",
			tokens: []token{
				token{
					Kind:  tokenKindChar,
					Value: 'a',
				},
				token{
					Kind:  tokenKindChar,
					Value: 'b',
				},
				token{
					Kind: tokenKindHLBegin,
				},
				token{
					Kind:  tokenKindChar,
					Value: 'c',
				},
				token{
					Kind: tokenKindHLEnd,
				},
				token{
					Kind:  tokenKindChar,
					Value: 'd',
				},
				token{
					Kind: tokenKindEOL,
				},
			},
		},
		// user writes a reserved token
		{
			input: "<dnotehl><dnotehl></dnotehl>",
			tokens: []token{
				token{
					Kind: tokenKindHLBegin,
				},
				token{
					Kind: tokenKindHLBegin,
				},
				token{
					Kind: tokenKindHLEnd,
				},
				token{
					Kind: tokenKindEOL,
				},
			},
		},
		{
			input: "<dnotehl></dnotehl></dnotehl>",
			tokens: []token{
				token{
					Kind: tokenKindHLBegin,
				},
				token{
					Kind: tokenKindHLEnd,
				},
				token{
					Kind: tokenKindHLEnd,
				},
				token{
					Kind: tokenKindEOL,
				},
			},
		},
	}

	for tcIdx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", tcIdx), func(t *testing.T) {
			tokens := tokenize(tc.input)

			assert.DeepEqual(t, tokens, tc.tokens, "tokens mismatch")
		})
	}
}
