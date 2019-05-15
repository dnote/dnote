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
  START_FETCHING,
  RECEIVE,
  RECEIVE_MORE,
  START_FETCHING_MORE,
  RECEIVE_ERROR
} from '../actions/digests';

const initialState = {
  items: [],
  total: 0,
  page: 0,
  isFetching: false,
  isFetched: false,
  error: null
};

export default function(state = initialState, action) {
  switch (action.type) {
    case START_FETCHING: {
      return {
        ...state,
        isFetching: true,
        isFetched: false
      };
    }
    case START_FETCHING_MORE: {
      return {
        ...state,
        isFetchingMore: true,
        hasFetchedMore: false
      };
    }
    case RECEIVE: {
      return {
        ...state,
        isFetching: false,
        isFetched: true,
        page: 1,
        total: action.data.total,
        items: action.data.items
      };
    }
    case RECEIVE_MORE: {
      return {
        ...state,
        page: state.page + 1,
        isFetchingMore: false,
        hasFetchedMore: true,
        items: [...state.items, ...action.data.items]
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
