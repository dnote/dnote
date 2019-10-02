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

import { getEditorSessionkey } from './editor';

describe('editor.ts', () => {
  describe('getEditorSessionkey', () => {
    const testCases = [
      {
        noteUUID: null,
        expected: 'new'
      },
      {
        noteUUID: '0ad88090-ab44-4432-be80-09c033f4c582',
        expected: '0ad88090-ab44-4432-be80-09c033f4c582'
      },
      {
        noteUUID: '6c20d136-8a15-443b-bd58-d2d963d38938',
        expected: '6c20d136-8a15-443b-bd58-d2d963d38938'
      }
    ];

    for (let i = 0; i < testCases.length; i++) {
      const tc = testCases[i];

      it(`generates a session key for input: ${tc.noteUUID}`, () => {
        const result = getEditorSessionkey(tc.noteUUID);
        expect(result).to.equal;
      });
    }
  });
});
