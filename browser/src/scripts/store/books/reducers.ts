/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import {
  START_FETCHING,
  RECEIVE,
  RECEIVE_ERROR,
  BooksState,
  BooksActionType
} from './types';

const initialState = {
  items: [],
  isFetching: false,
  isFetched: false,
  error: null
};

export default function (
  state = initialState,
  action: BooksActionType
): BooksState {
  switch (action.type) {
    case START_FETCHING:
      return {
        ...state,
        isFetching: true,
        isFetched: false
      };
    case RECEIVE: {
      const { books } = action.data;

      // get uuids of deleted books and that of a currently selected book
      return {
        ...state,
        isFetching: false,
        isFetched: true,
        items: [...state.items, ...books]
      };
    }
    case RECEIVE_ERROR:
      return {
        ...state,
        isFetching: false,
        isFetched: true,
        error: action.data.error
      };
    default:
      return state;
  }
}
