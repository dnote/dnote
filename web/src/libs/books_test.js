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
import {
  validateBookName,
  checkDuplicate,
  errBookNameNumeric,
  errBookNameHasSpace,
  errBookNameReserved
} from './books';

describe('validateBookName', () => {
  const testCases = [
    {
      input: 'javascript',
      expectedErr: null
    },
    {
      input: 'node.js',
      expectedErr: null
    },
    {
      input: 'foo bar',
      expectedErr: errBookNameHasSpace
    },
    {
      input: '123',
      expectedErr: errBookNameNumeric
    },
    {
      input: '+123',
      expectedErr: null
    },
    {
      input: '-123',
      expectedErr: null
    },
    {
      input: '+javascript',
      expectedErr: null
    },
    {
      input: '0',
      expectedErr: errBookNameNumeric
    },
    {
      input: '0333',
      expectedErr: errBookNameNumeric
    },
    {
      input: ' javascript',
      expectedErr: errBookNameHasSpace
    },
    {
      input: 'java script',
      expectedErr: errBookNameHasSpace
    },
    {
      input: 'javascript (1)',
      expectedErr: errBookNameHasSpace
    },
    {
      input: 'javascript ',
      expectedErr: errBookNameHasSpace
    },
    {
      input: 'javascript (1) (2) (3)',
      expectedErr: errBookNameHasSpace
    },

    // reserved book names
    {
      input: 'trash',
      expectedErr: errBookNameReserved
    },
    {
      input: 'conflicts',
      expectedErr: errBookNameReserved
    }
  ];

  for (let i = 0; i < testCases.length; ++i) {
    const tc = testCases[i];

    it(`validates ${tc.input}`, () => {
      const base = expect(() => validateBookName(tc.input));

      if (tc.expectedErr) {
        base.to.throw(tc.expectedErr);
      } else {
        base.to.not.throw();
      }
    });
  }
});

describe('checkDuplicate', () => {
  const testCases = [
    {
      books: [],
      bookName: 'javascript',
      expected: false
    },
    {
      books: [{ label: 'golang' }, { label: 'cool_ideas' }, { label: 'foo' }],
      bookName: 'bar1',
      expected: false
    },
    {
      books: [{ label: 'golang' }, { label: 'foo' }, { label: 'bar' }],
      bookName: 'bar',
      expected: true
    }
  ];

  for (let i = 0; i < testCases.length; ++i) {
    const tc = testCases[i];

    it(`checks duplicate for the test case ${i}`, () => {
      const result = checkDuplicate(tc.books, tc.bookName);
      expect(result).to.equal(tc.expected);
    });
  }
});
