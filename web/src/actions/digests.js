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

export const START_FETCHING_DIGESTS = 'digests/START_FETCHING_DIGESTS';
export const RECEIVE_DIGESTS = 'digests/RECEIVE_DIGESTS';
export const RECEIVE_ERROR = 'digests/RECEIVE_ERROR';

function receiveDigests(digests) {
  return {
    type: RECEIVE_DIGESTS,
    data: {
      digests
    }
  };
}

function startFetchingDigests() {
  return {
    type: START_FETCHING_DIGESTS
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

export function getDigests() {
  return async dispatch => {
    try {
      dispatch(startFetchingDigests());

      const digests = await digestsService.fetchAll();

      dispatch(receiveDigests(digests));
    } catch (err) {
      console.log('Error fetching digests', err.stack);
      dispatch(receiveError(err.message));
    }
  };
}
