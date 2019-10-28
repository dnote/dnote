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
import { BookData } from 'jslib/operations/types';
import { RECEIVE, ADD, REMOVE, START_FETCHING, FINISH_FETCHING } from './type';
import { ThunkAction } from '../types';

function receiveBooks(books: BookData[]) {
  return {
    type: RECEIVE,
    data: { books }
  };
}

function startFetchingBooks() {
  return {
    type: START_FETCHING
  };
}

function finishFetchingBooks() {
  return {
    type: FINISH_FETCHING
  };
}

export const getBooks = (): ThunkAction<void> => {
  return dispatch => {
    dispatch(startFetchingBooks());

    return operations.books
      .fetch()
      .then(books => {
        dispatch(receiveBooks(books));
        dispatch(finishFetchingBooks());
      })
      .catch(err => {
        console.log('getBooks error', err);
        // todo: handle error
      });
  };
};

export function addBook(book: BookData) {
  return {
    type: ADD,
    data: { book }
  };
}

export const createBook = (name: string): ThunkAction<BookData> => {
  return dispatch => {
    return operations.books.create({ name }).then(book => {
      dispatch(addBook(book));

      return book;
    });
  };
};

export function removeBook(bookUUID: string) {
  return {
    type: REMOVE,
    data: { bookUUID }
  };
}
