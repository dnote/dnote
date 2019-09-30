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
  BooksState,
  BooksActionType,
  RECEIVE,
  ADD,
  REMOVE,
  START_FETCHING,
  FINISH_FETCHING
} from './type';

const initialState: BooksState = {
  data: [],
  isFetching: false,
  isFetched: false,
  errorMessage: ''
};

export default function(
  state = initialState,
  action: BooksActionType
): BooksState {
  switch (action.type) {
    case START_FETCHING: {
      return {
        ...state,
        isFetching: true,
        isFetched: false
      };
    }
    case FINISH_FETCHING: {
      return {
        ...state,
        isFetching: false,
        isFetched: true
      };
    }
    case RECEIVE: {
      return {
        ...state,
        data: action.data.books
      };
    }
    case REMOVE: {
      return {
        ...state,
        data: state.data.filter(item => {
          return item.uuid !== action.data.bookUUID;
        })
      };
    }
    case ADD: {
      const data = [...state.data, action.data.book];
      const sorted = data.sort((a, b) => {
        const aLabel = a.label.toLowerCase();
        const bLabel = b.label.toLowerCase();

        return aLabel.localeCompare(bLabel);
      });

      return {
        ...state,
        data: sorted
      };
    }
    default:
      return state;
  }
}
