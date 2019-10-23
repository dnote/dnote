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

import { populateParams } from './paths';

describe('paths.ts', () => {
  describe('populateParams', () => {
    const testCases = [
      {
        pathDef: '/foo/:bar',
        params: {
          bar: '123'
        },
        expected: '/foo/123'
      },
      {
        pathDef: '/foo/:bar/baz',
        params: {
          bar: '123'
        },
        expected: '/foo/123/baz'
      },
      {
        pathDef: '/foo/:bar/:baz/:quz/qux',
        params: {
          bar: '123',
          baz: '456',
          quz: 'abcd'
        },
        expected: '/foo/123/456/abcd/qux'
      }
    ];

    for (let i = 0; i < testCases.length; i++) {
      const tc = testCases[i];

      const stringifiedParams = JSON.stringify(tc.params);
      it(`populates ${tc.pathDef} with params ${stringifiedParams}`, () => {
        const result = populateParams(tc.pathDef, tc.params);
        expect(result).to.equal(tc.expected);
      });
    }
  });
});
