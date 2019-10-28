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

import formatTime from './format';

describe('time/format.ts', () => {
  describe('formatTime', () => {
    const date = new Date(2017, 2, 30, 8, 30);
    const testCases = [
      {
        format: '%YYYY %MM %DD %hh:%mm',
        expected: '2017 03 30 08:30'
      },
      {
        format: '%YYYY %MMM %DD %hh:%mm %A',
        expected: '2017 Mar 30 08:30 AM'
      },
      {
        format: '%YYYY %MMM %DD %h:%mm%a',
        expected: '2017 Mar 30 8:30am'
      },
      {
        format: '%dddd %M/%D',
        expected: 'Thursday 3/30'
      },
      {
        format: '%MMMM %Do',
        expected: 'March 30th'
      }
    ];

    for (let i = 0; i < testCases.length; i++) {
      const tc = testCases[i];

      it(`converts the input ${tc.format}`, () => {
        const result = formatTime(date, tc.format);
        expect(result).to.equal(tc.expected);
      });
    }
  });
});
