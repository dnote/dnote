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

export const START_FETCHING = 'digests/START_FETCHING';
export const START_FETCHING_MORE = 'digests/START_FETCHING_MORE';
export const RECEIVE = 'digests/RECEIVE';
export const RECEIVE_MORE = 'digests/RECEIVE_MORE';
export const RECEIVE_ERROR = 'digests/RECEIVE_ERROR';

function receiveDigests(total, items) {
  return {
    type: RECEIVE,
    data: {
      total,
      items
    }
  };
}

function receiveMoreDigests(items) {
  return {
    type: RECEIVE_MORE,
    data: {
      items
    }
  };
}

function startFetchingDigests() {
  return {
    type: START_FETCHING
  };
}

function startFetchingMoreDigests() {
  return {
    type: START_FETCHING_MORE
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

export function getDigests(option = {}) {
  const { demo } = option;

  return async dispatch => {
    try {
      dispatch(startFetchingDigests());

      const res = await digestsService.fetchAll({ demo });

      dispatch(receiveDigests(res.total, res.digests));
    } catch (err) {
      console.log('Error fetching digests', err.stack);
      dispatch(receiveError(err.message));
    }
  };
}

export function getMoreDigests(option = {}) {
  const { demo } = option;

  return async (dispatch, getState) => {
    try {
      dispatch(startFetchingMoreDigests());

      const { digests } = getState();
      const nextPage = digests.page + 1;
      const res = await digestsService.fetchAll({ page: nextPage, demo });

      dispatch(receiveMoreDigests(res.digests));
    } catch (err) {
      console.log('Error fetching digests', err.stack);
      dispatch(receiveError(err.message));
    }
  };
}
