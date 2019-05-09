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

import {
  START_FETCHING_DIGEST_NOTES,
  RECEIVE_DIGEST_NOTES,
  RECEIVE_ERROR
} from '../actions/digest';

const initialState = {
  notes: [],
  isFetching: false,
  isFetched: false,
  error: null
};

export default function(state = initialState, action) {
  switch (action.type) {
    case START_FETCHING_DIGEST_NOTES: {
      return {
        ...state,
        isFetching: true,
        isFetched: false
      };
    }
    case RECEIVE_DIGEST_NOTES: {
      return {
        ...state,
        isFetching: false,
        isFetched: true,
        notes: action.data.notes
      };
    }
    case RECEIVE_ERROR: {
      return {
        ...state,
        error: action.data.error
      };
    }
    default:
      return state;
  }
}
