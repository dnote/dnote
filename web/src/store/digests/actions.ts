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
import { START_FETCHING, RECEIVE, RECEIVE_ERROR, RESET } from './type';

function receiveDigests(total, items, page) {
  return {
    type: RECEIVE,
    data: {
      total,
      items,
      page
    }
  };
}

function startFetchingDigests() {
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

export function resetDigests() {
  return {
    type: RESET
  };
}

export function getDigests(params: { page: number; status: string }) {
  return async dispatch => {
    try {
      dispatch(startFetchingDigests());

      const res = await operations.digests.fetchAll(params);

      dispatch(receiveDigests(res.total, res.items, params.page));
    } catch (err) {
      console.log('Error fetching digests', err.stack);
      dispatch(receiveError(err.message));
    }
  };
}
