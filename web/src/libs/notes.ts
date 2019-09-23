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

import { nanosecToMillisec } from '../helpers/time';
import { NoteData } from '../../../jslib/src/operations/types';

export interface NotesGroupData {
  year: number;
  month: number;
  data: NoteData[];
}

function encodeGroupKey(year: number, month: number): string {
  return `${year}-${month}`;
}

function decodeGroupKey(key: string): { year: number; month: number } {
  const [yearStr, monthStr] = key.split('-');

  const year = parseInt(yearStr, 10);
  const month = parseInt(monthStr, 10);

  return { year, month };
}

function makeGroup(
  year: number,
  month: number,
  notes: NoteData[]
): NotesGroupData {
  return {
    year,
    month,
    data: notes
  };
}

// groupNotes groups the notes to note groups
export function groupNotes(notes: NoteData[]): NotesGroupData[] {
  const ret: NotesGroupData[] = [];

  const map: { [key: string]: NoteData[] } = {};

  for (let i = 0; i < notes.length; i++) {
    const note = notes[i];

    const date = new Date(nanosecToMillisec(note.added_on));
    const year = date.getUTCFullYear();
    const month = date.getUTCMonth() + 1;

    const key = encodeGroupKey(year, month);

    if (map[key]) {
      map[key].push(note);
    } else {
      map[key] = [note];
    }
  }

  const keys = Object.keys(map);
  for (let i = 0; i < keys.length; i++) {
    const key = keys[i];
    const items = map[key];

    const { year, month } = decodeGroupKey(key);

    const group = makeGroup(year, month, items);
    ret.push(group);
  }

  return ret;
}
