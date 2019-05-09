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

import * as booksService from '../services/books';
import { decryptBook } from '../crypto/books';

export const RECEIVE = 'books/RECEIVE';
export const ADD = 'books/ADD';
export const REMOVE = 'books/REMOVE';
export const START_FETCHING = 'books/START_FETCHING';
export const FINISH_FETCHING = 'books/FINISH_FETCHING';

function receiveBooks(books) {
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

export function getBooks(cipherKeyBuf, demo) {
  return dispatch => {
    dispatch(startFetchingBooks());

    return booksService
      .fetch({}, { demo })
      .then(books => {
        const p = books.map(book => {
          return decryptBook(book, cipherKeyBuf);
        });

        return Promise.all(p).then(booksDec => {
          dispatch(receiveBooks(booksDec));
          dispatch(finishFetchingBooks());
        });
      })
      .catch(err => {
        console.log('getBooks error', err);
        // todo: handle error
      });
  };
}

export function addBook(book) {
  return {
    type: ADD,
    data: { book }
  };
}

export function removeBook(bookUUID) {
  return {
    type: REMOVE,
    data: { bookUUID }
  };
}
