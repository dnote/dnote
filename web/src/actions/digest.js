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

import { fetchDigestNotes } from '../services/digests';
import { decryptNote } from '../crypto/notes';

export const START_FETCHING_DIGEST_NOTES = 'digest/START_FETCHING_DIGEST_NOTES';
export const RECEIVE_DIGEST_NOTES = 'digest/RECEIVE_DIGEST_NOTES';
export const RECEIVE_ERROR = 'digest/RECEIVE_ERROR';

function receiveDigestNotes(notes) {
  return {
    type: RECEIVE_DIGEST_NOTES,
    data: {
      notes
    }
  };
}

function startFetchingDigestNotes() {
  return {
    type: START_FETCHING_DIGEST_NOTES
  };
}

function receiveError(error) {
  return {
    type: RECEIVE_ERROR,
    data: {
      error
    }
  };
}

export function getDigestNotes(cipherKeyBuf, digestUUID) {
  return async dispatch => {
    try {
      dispatch(startFetchingDigestNotes());

      const notes = await fetchDigestNotes(digestUUID);

      const p = notes.map(note => {
        return decryptNote(note, cipherKeyBuf);
      });

      const notesDec = await Promise.all(p);
      dispatch(receiveDigestNotes(notesDec));
    } catch (err) {
      console.log('Error fetching digest notes', err.stack);
      dispatch(receiveError(err.message));
    }
  };
}
