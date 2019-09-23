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

import operations from 'web/libs/operations';
import { NoteData } from 'jslib/operations/types';
import { RECEIVE, START_FETCHING, ERROR, RESET } from './type';
import { ThunkAction } from '../types';

export function receiveNote(note) {
  return {
    type: RECEIVE,
    data: { note }
  };
}

export function resetNote() {
  return {
    type: RESET
  };
}

function startFetchingNote() {
  return {
    type: START_FETCHING
  };
}

function receiveNoteError(errorMessage) {
  return {
    type: ERROR,
    data: { errorMessage }
  };
}

interface GetNoteFacets {
  q?: string;
}

export const getNote = (
  noteUUID: string,
  params: GetNoteFacets
): ThunkAction<NoteData | void> => {
  return dispatch => {
    dispatch(startFetchingNote());

    return operations.notes
      .fetchOne(noteUUID, params)
      .then(note => {
        dispatch(receiveNote(note));

        return note;
      })
      .catch(err => {
        console.log('getNote error', err.message);
        dispatch(receiveNoteError(err.message));
      });
  };
};
