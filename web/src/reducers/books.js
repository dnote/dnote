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
  ADD,
  REMOVE,
  START_FETCHING,
  FINISH_FETCHING
} from '../actions/books';

const initialState = {
  items: [],
  isFetching: false,
  isFetched: false
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
        items: action.data.books
      };
    }
    case REMOVE: {
      return {
        ...state,
        items: state.items.filter(item => {
          return item.uuid !== action.data.bookUUID;
        })
      };
    }
    case ADD: {
      return {
        ...state,
        items: [...state.items, action.data.book]
      };
    }
    default:
      return state;
  }
}
