package find

import (
	"fmt"
	"testing"

	"github.com/dnote/cli/testutils"
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

			testutils.AssertEqual(t, nextIdx, tc.retIdx, "retIdx mismatch")
			testutils.AssertDeepEqual(t, tok, tc.retTok, "retTok mismatch")
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

			testutils.AssertDeepEqual(t, tokens, tc.tokens, "tokens mismatch")
		})
	}
}
