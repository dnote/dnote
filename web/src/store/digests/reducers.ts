/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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
  RECEIVE_ERROR,
  RESET,
  DigestsState,
  DigestsActionType
} from './type';

const initialState: DigestsState = {
  data: [],
  total: 0,
  page: 0,
  isFetching: false,
  isFetched: false,
  errorMessage: null
};

export default function(
  state = initialState,
  action: DigestsActionType
): DigestsState {
  switch (action.type) {
    case START_FETCHING: {
      return {
        ...state,
        errorMessage: null,
        isFetching: true,
        isFetched: false
      };
    }
    case RECEIVE: {
      return {
        ...state,
        isFetching: false,
        isFetched: true,
        total: action.data.total,
        page: action.data.page,
        data: action.data.items
      };
    }
    case RECEIVE_ERROR: {
      return {
        ...state,
        errorMessage: action.data.error
      };
    }
    case RESET: {
      return initialState;
    }
    default:
      return state;
  }
}
