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

import { daysToSec } from './time';

describe('time.ts', () => {
  describe('daysToSec', () => {
    const testCases = [
      {
        input: 1,
        expected: 86400
      },

      {
        input: 14,
        expected: 1209600
      }
    ];

    for (let i = 0; i < testCases.length; i++) {
      const tc = testCases[i];

      it(`converts the input ${tc.input}`, () => {
        const result = daysToSec(tc.input);
        expect(result).to.equal(tc.expected);
      });
    }
  });
});
