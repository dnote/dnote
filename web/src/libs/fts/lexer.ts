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

export enum TokenKind {
  char,
  hlBegin,
  hlEnd,
  eol
}

export interface Token {
  kind: TokenKind;
  value?: string;
}

interface ScanTokenResult {
  tok: Token;
  nextIdx: number;
}

// getNextIdx validates that the given index is within the range of the given string.
// If so, it returns the given index. Otherwise it returns -1.
function getNextIdx(candidate: number, s: string): number {
  if (candidate <= s.length - 1) {
    return candidate;
  }

  return -1;
}

// scanToken scans the given string for a token at the given index. It returns
// a token and the next index to look for a token. If the given string is exhausted,
// the next index will be -1.
export function scanToken(idx: number, s: string): ScanTokenResult {
  if (s[idx] === '<') {
    if (s.length - idx >= 9) {
      const lookahead = 9;
      const candidate = s.substring(idx, idx + lookahead);

      if (candidate === '<dnotehl>') {
        return {
          tok: {
            kind: TokenKind.hlBegin
          },
          nextIdx: getNextIdx(idx + lookahead, s)
        };
      }
    }

    if (s.length - idx >= 10) {
      const lookahead = 10;
      const candidate = s.substring(idx, idx + lookahead);

      if (candidate === '</dnotehl>') {
        return {
          tok: {
            kind: TokenKind.hlEnd
          },
          nextIdx: getNextIdx(idx + lookahead, s)
        };
      }
    }
  }

  const nextIdx = getNextIdx(idx + 1, s);
  return {
    tok: {
      value: s[idx],
      kind: TokenKind.char
    },
    nextIdx
  };
}

// tokenize lexically analyzes the given matched snippet from a full text search
// and builds a slice of tokens
export function tokenize(s: string): Token[] {
  const ret: Token[] = [];

  let idx = 0;
  while (idx !== -1) {
    const { tok, nextIdx } = scanToken(idx, s);

    idx = nextIdx;
    ret.push(tok);
  }

  ret.push({ kind: TokenKind.eol });

  return ret;
}
