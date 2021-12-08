/* Copyright (C) 2019, 2020, 2021 Monomax Software Pty Ltd
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

import { getRange } from './arr';

describe('getRange', () => {
  const testCases = [
    {
      input: 1,
      expected: [1]
    },
    {
      input: 3,
      expected: [1, 2, 3]
    }
  ];

  for (let i = 0; i < testCases.length; ++i) {
    const tc = testCases[i];

    test(`generates a range for ${tc.input}`, () => {
      const result = getRange(tc.input);

      expect(result).toStrictEqual(tc.expected);
    });
  }
});
