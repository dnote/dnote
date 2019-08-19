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
  ADD,
  REFRESH,
  RECEIVE,
  START_FETCHING,
  RECEIVE_ERROR,
  RESET,
  REMOVE,
  NotesActionType,
  NotesState
} from './type';

const initialState: NotesState = {
  data: [],
  total: 0,
  isFetching: false,
  isFetched: false,
  errorMessage: null
};

export default function(
  state = initialState,
  action: NotesActionType
): NotesState {
  switch (action.type) {
    case START_FETCHING: {
      return {
        ...state,
        isFetched: false,
        isFetching: true,
        errorMessage: ''
      };
    }
    case ADD: {
      return state;
    }
    case REFRESH: {
      return state;
    }
    case REMOVE: {
      return state;
    }
    case RECEIVE: {
      const { notes, total } = action.data;

      return {
        ...state,
        data: notes,
        total,
        isFetched: true,
        isFetching: false
      };
    }
    case RECEIVE_ERROR: {
      return {
        ...state,
        errorMessage: action.data.error,
        isFetching: false,
        isFetched: false
      };
    }
    case RESET: {
      return initialState;
    }
    default:
      return state;
  }
}
