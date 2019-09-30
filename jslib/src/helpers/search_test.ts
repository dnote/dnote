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
import { TokenKind, tokenize, parse } from './search';

describe('search.ts', () => {
  describe('tokenize', () => {
    const testCases = [
      {
        input: 'foo',
        tokens: [
          {
            kind: TokenKind.id,
            value: 'foo'
          },
          {
            kind: TokenKind.eof
          }
        ]
      },
      {
        input: '123',
        tokens: [
          {
            kind: TokenKind.id,
            value: '123'
          },
          {
            kind: TokenKind.eof
          }
        ]
      },
      {
        input: 'foo123',
        tokens: [
          {
            kind: TokenKind.id,
            value: 'foo123'
          },
          {
            kind: TokenKind.eof
          }
        ]
      },
      {
        input: 'foo\tbar',
        tokens: [
          {
            kind: TokenKind.id,
            value: 'foo'
          },
          {
            kind: TokenKind.id,
            value: 'bar'
          },
          {
            kind: TokenKind.eof
          }
        ]
      },
      {
        input: ' foo \tbar\t',
        tokens: [
          {
            kind: TokenKind.id,
            value: 'foo'
          },
          {
            kind: TokenKind.id,
            value: 'bar'
          },
          {
            kind: TokenKind.eof
          }
        ]
      },
      {
        input: 'foo:bar',
        tokens: [
          {
            kind: TokenKind.id,
            value: 'foo'
          },
          {
            kind: TokenKind.colon
          },
          {
            kind: TokenKind.id,
            value: 'bar'
          },
          {
            kind: TokenKind.eof
          }
        ]
      },
      {
        input: '"foo" bar',
        tokens: [
          {
            kind: TokenKind.id,
            value: '"foo"'
          },
          {
            kind: TokenKind.id,
            value: 'bar'
          },
          {
            kind: TokenKind.eof
          }
        ]
      },
      {
        input: '"foo:bar"',
        tokens: [
          {
            kind: TokenKind.id,
            value: '"foo'
          },
          {
            kind: TokenKind.colon
          },
          {
            kind: TokenKind.id,
            value: 'bar"'
          },
          {
            kind: TokenKind.eof
          }
        ]
      }
    ];

    for (let i = 0; i < testCases.length; i++) {
      const tc = testCases[i];

      it(`tokenizes ${tc.input}`, () => {
        const result = tokenize(tc.input);

        expect(result).to.eql(tc.tokens);
      });
    }
  });

  describe('parse', () => {
    function run(testCases) {
      for (let i = 0; i < testCases.length; i++) {
        const tc = testCases[i];

        it(`keyword [${tc.keywords}] - parses ${tc.input} `, () => {
          const result = parse(tc.input, tc.keywords);

          expect(result).to.eql(tc.result);
        });
      }
    }

    describe('text only', () => {
      const testCases = [
        {
          input: 'foo',
          keywords: [],
          result: {
            text: 'foo',
            filters: {}
          }
        },
        {
          input: '123',
          keywords: [],
          result: {
            text: '123',
            filters: {}
          }
        },
        {
          input: 'foo123',
          keywords: [],
          result: {
            text: 'foo123',
            filters: {}
          }
        },
        {
          input: '"',
          keywords: [],
          result: {
            text: '"',
            filters: {}
          }
        },
        {
          input: '""',
          keywords: [],
          result: {
            text: '""',
            filters: {}
          }
        },
        {
          input: `'`,
          keywords: [],
          result: {
            text: `'`,
            filters: {}
          }
        },
        {
          input: `''`,
          keywords: [],
          result: {
            text: `''`,
            filters: {}
          }
        },
        {
          input: `'foo:bar'`,
          keywords: [],
          result: {
            text: `'foo:bar'`,
            filters: {}
          }
        },
        {
          input: 'foo bar',
          keywords: [],
          result: {
            text: 'foo bar',
            filters: {}
          }
        },
        {
          input: ' foo \t bar ',
          keywords: [],
          result: {
            text: 'foo bar',
            filters: {}
          }
        },
        {
          input: '"foo:bar"',
          keywords: ['foo'],
          result: {
            text: '"foo:bar"',
            filters: {}
          }
        },
        {
          input: '"foo:bar""',
          keywords: ['foo'],
          result: {
            text: '"foo:bar""',
            filters: {}
          }
        },
        {
          input: '"foo:bar""""',
          keywords: ['foo'],
          result: {
            text: '"foo:bar""""',
            filters: {}
          }
        }
      ];

      run(testCases);
    });

    describe('filter only', () => {
      const testCases = [
        {
          input: 'foo:bar',
          keywords: ['foo'],
          result: {
            text: '',
            filters: {
              foo: 'bar'
            }
          }
        },
        {
          input: '123:bar',
          keywords: ['123'],
          result: {
            text: '',
            filters: {
              '123': 'bar'
            }
          }
        },
        {
          input: 'foo123:bar',
          keywords: ['foo123'],
          result: {
            text: '',
            filters: {
              foo123: 'bar'
            }
          }
        },
        {
          input: '123:456',
          keywords: ['123'],
          result: {
            text: '',
            filters: {
              '123': '456'
            }
          }
        },
        {
          input: 'foo:bar baz:quz 123:qux',
          keywords: ['foo'],
          result: {
            text: 'baz:quz 123:qux',
            filters: {
              foo: 'bar'
            }
          }
        },
        {
          input: 'foo:bar baz:quz',
          keywords: ['foo'],
          result: {
            text: 'baz:quz',
            filters: {
              foo: 'bar'
            }
          }
        },
        {
          input: 'foo:bar baz:quz',
          keywords: ['bar'],
          result: {
            text: 'foo:bar baz:quz',
            filters: {}
          }
        },
        {
          input: 'foo:bar baz:quz',
          keywords: ['foo', 'baz'],
          result: {
            text: '',
            filters: {
              foo: 'bar',
              baz: 'quz'
            }
          }
        },
        {
          input: 'foo:bar',
          keywords: [],
          result: {
            text: 'foo:bar',
            filters: {}
          }
        },
        {
          input: 'foo:bar baz:quz',
          keywords: [],
          result: {
            text: 'foo:bar baz:quz',
            filters: {}
          }
        },
        {
          input: 'foo:bar foo:baz',
          keywords: ['foo'],
          result: {
            text: '',
            filters: {
              foo: ['bar', 'baz']
            }
          }
        }
      ];

      run(testCases);
    });

    describe('text and filter', () => {
      const testCases = [
        {
          input: 'foo:bar baz',
          keywords: ['foo'],
          result: {
            text: 'baz',
            filters: {
              foo: 'bar'
            }
          }
        },
        {
          input: 'foo:bar baz quz:qux1 ',
          keywords: ['foo', 'quz'],
          result: {
            text: 'baz',
            filters: {
              foo: 'bar',
              quz: 'qux1'
            }
          }
        },
        {
          input: 'foo:bar baz quz:qux1 qux',
          keywords: ['foo', 'quz'],
          result: {
            text: 'baz qux',
            filters: {
              foo: 'bar',
              quz: 'qux1'
            }
          }
        },
        {
          input: 'foo:bar baz quz:qux1 qux "quux:fooz"',
          keywords: ['foo', 'quz'],
          result: {
            text: 'baz qux "quux:fooz"',
            filters: {
              foo: 'bar',
              quz: 'qux1'
            }
          }
        },
        {
          input: 'foo:bar baz quz:qux1 qux "quux:fooz"',
          keywords: ['foo', 'quux'],
          result: {
            text: 'baz quz:qux1 qux "quux:fooz"',
            filters: {
              foo: 'bar'
            }
          }
        }
      ];

      run(testCases);
    });
  });
});
