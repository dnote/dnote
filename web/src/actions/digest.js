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

import * as digestsService from '../services/digests';
import { decryptDigest } from '../crypto/digest';

export const START_FETCHING = 'digest/START_FETCHING';
export const RECEIVE = 'digest/RECEIVE';
export const RECEIVE_ERROR = 'digest/RECEIVE_ERROR';

function receiveDigest(item) {
  return {
    type: RECEIVE,
    data: {
      item
    }
  };
}

function startFetchingDigestNotes() {
  return {
    type: START_FETCHING
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

export function getDigest(cipherKeyBuf, digestUUID, demo) {
  return async dispatch => {
    try {
      dispatch(startFetchingDigestNotes());

      const digest = await digestsService.fetch(digestUUID, { demo });
      const digestDec = await decryptDigest(digest, cipherKeyBuf);
      dispatch(receiveDigest(digestDec));
    } catch (err) {
      console.log('Error fetching digest', err.stack);
      dispatch(receiveError(err.message));
    }
  };
}
