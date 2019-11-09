/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import { expect } from 'chai';
import { TokenKind, tokenize, scanToken } from './lexer';

describe('lexer.ts', () => {
  describe('scanToken', () => {
    const testCases = [
      {
        input: 'foo bar',
        idx: 1,
        retTok: { value: 'o', kind: TokenKind.char },
        retIdx: 2
      },
      {
        input: 'foo bar',
        idx: 6,
        retTok: { value: 'r', kind: TokenKind.char },
        retIdx: -1
      },
      {
        input: 'foo <bar>',
        idx: 4,
        retTok: { value: '<', kind: TokenKind.char },
        retIdx: 5
      },
      {
        input: 'foo <dnotehL>',
        idx: 4,
        retTok: { value: '<', kind: TokenKind.char },
        retIdx: 5
      },
      {
        input: 'foo <dnotehl>bar</dnotehl> foo bar',
        idx: 4,
        retTok: { kind: TokenKind.hlBegin },
        retIdx: 13
      },
      {
        input: 'foo <dnotehl>bar</dnotehl> <dnotehl>foo</dnotehl> bar',
        idx: 4,
        retTok: { kind: TokenKind.hlBegin },
        retIdx: 13
      },
      {
        input: 'foo <dnotehl>bar</dnotehl> <dnotehl>foo</dnotehl> bar',
        idx: 27,
        retTok: { kind: TokenKind.hlBegin },
        retIdx: 36
      },
      {
        input: 'foo <dnotehl>bar</dnotehl> foo bar',
        idx: 13,
        retTok: { value: 'b', kind: TokenKind.char },
        retIdx: 14
      },
      {
        input: 'foo <dnotehl>bar</dnotehl> foo bar',
        idx: 16,
        retTok: { kind: TokenKind.hlEnd },
        retIdx: 26
      },
      {
        input: '<dno<dnotehl>tehl>',
        idx: 0,
        retTok: { value: '<', kind: TokenKind.char },
        retIdx: 1
      },
      {
        input: '<dno<dnotehl>tehl>',
        idx: 4,
        retTok: { kind: TokenKind.hlBegin },
        retIdx: 13
      },
      {
        input: 'foo <dnotehl>bar</dnotehl>',
        idx: 16,
        retTok: { kind: TokenKind.hlEnd },
        retIdx: -1
      },
      // user writes reserved token
      {
        input: 'foo <dnotehl>',
        idx: 4,
        retTok: { kind: TokenKind.hlBegin },
        retIdx: -1
      }
    ];

    for (let i = 0; i < testCases.length; i++) {
      const tc = testCases[i];

      it(`scans ${tc.input}`, () => {
        const result = scanToken(tc.idx, tc.input);

        expect(result.tok).to.deep.equal(tc.retTok);
      });
    }
  });

  describe('tokenize', () => {
    const testCases = [
      {
        input: 'ab<dnotehl>c</dnotehl>',
        tokens: [
          {
            kind: TokenKind.char,
            value: 'a'
          },
          {
            kind: TokenKind.char,
            value: 'b'
          },
          {
            kind: TokenKind.hlBegin
          },
          {
            kind: TokenKind.char,
            value: 'c'
          },
          {
            kind: TokenKind.hlEnd
          },
          {
            kind: TokenKind.eol
          }
        ]
      },
      {
        input: 'ab<dnotehl>c</dnotehl>d',
        tokens: [
          {
            kind: TokenKind.char,
            value: 'a'
          },
          {
            kind: TokenKind.char,
            value: 'b'
          },
          {
            kind: TokenKind.hlBegin
          },
          {
            kind: TokenKind.char,
            value: 'c'
          },
          {
            kind: TokenKind.hlEnd
          },
          {
            kind: TokenKind.char,
            value: 'd'
          },
          {
            kind: TokenKind.eol
          }
        ]
      },
      // user writes a reserved token
      {
        input: '<dnotehl><dnotehl></dnotehl>',
        tokens: [
          {
            kind: TokenKind.hlBegin
          },
          {
            kind: TokenKind.hlBegin
          },
          {
            kind: TokenKind.hlEnd
          },
          {
            kind: TokenKind.eol
          }
        ]
      },
      {
        input: '<dnotehl></dnotehl></dnotehl>',
        tokens: [
          {
            kind: TokenKind.hlBegin
          },
          {
            kind: TokenKind.hlEnd
          },
          {
            kind: TokenKind.hlEnd
          },
          {
            kind: TokenKind.eol
          }
        ]
      }
    ];

    for (let i = 0; i < testCases.length; i++) {
      const tc = testCases[i];

      it(`tokenizes ${tc.input}`, () => {
        const result = tokenize(tc.input);

        expect(result).to.deep.equal(tc.tokens);
      });
    }
  });
});
