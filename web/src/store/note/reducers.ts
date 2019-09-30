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
  RECEIVE,
  START_FETCHING,
  ERROR,
  RESET,
  NoteState,
  NoteActionType
} from './type';

const initialState: NoteState = {
  data: {
    uuid: '',
    created_at: '',
    updated_at: '',
    content: '',
    added_on: 0,
    public: false,
    usn: 0,
    book: {
      uuid: '',
      label: ''
    },
    user: {
      name: '',
      uuid: ''
    }
  },
  isFetching: false,
  isFetched: false,
  errorMessage: null
};

export default function(
  state = initialState,
  action: NoteActionType
): NoteState {
  switch (action.type) {
    case START_FETCHING: {
      return {
        ...state,
        isFetching: true,
        isFetched: false
      };
    }
    case ERROR: {
      return {
        ...state,
        isFetching: false,
        errorMessage: action.data.errorMessage
      };
    }
    case RECEIVE: {
      return {
        ...state,
        data: action.data.note,
        isFetching: false,
        isFetched: true,
        errorMessage: null
      };
    }
    case RESET: {
      return initialState;
    }
    default:
      return state;
  }
}
