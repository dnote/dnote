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

import * as noteService from '../services/notes';
import { decryptNote } from '../crypto/notes';

export const RECEIVE = 'note/RECEIVE';
export const START_FETCHING = 'note/START_FETCHING';
export const ERROR = 'note/ERROR';
export const RESET = 'note/RESET';

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

export function getNote(cipherKeyBuf, noteUUID, demo = false) {
  return dispatch => {
    dispatch(startFetchingNote());

    return noteService
      .fetchOne(noteUUID, { demo })
      .then(note => {
        return decryptNote(note, cipherKeyBuf)
          .then(noteDec => {
            dispatch(receiveNote(noteDec));

            return noteDec;
          })
          .catch(err => {
            dispatch(receiveNoteError(err.message));
            return note;
          });
      })
      .catch(err => {
        console.log('getNote error', err.message);
        dispatch(receiveNoteError(err.message));
      });
  };
}
