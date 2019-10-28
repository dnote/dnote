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

import { RemoteData } from '../types';
import { BookData } from 'jslib/operations/types';

export type BooksState = RemoteData<BookData[]>;

export const RECEIVE = 'books/RECEIVE';
export const ADD = 'books/ADD';
export const REMOVE = 'books/REMOVE';
export const START_FETCHING = 'books/START_FETCHING';
export const FINISH_FETCHING = 'books/FINISH_FETCHING';

export interface ReceiveBooks {
  type: typeof RECEIVE;
  data: {
    books: BookData[];
  };
}

export interface StartFetchingBooks {
  type: typeof START_FETCHING;
}

export interface FinishFetchingBooks {
  type: typeof FINISH_FETCHING;
}

export interface AddBook {
  type: typeof ADD;
  data: {
    book: BookData;
  };
}

export interface RemoveBook {
  type: typeof REMOVE;
  data: {
    bookUUID: string;
  };
}

export type BooksActionType =
  | ReceiveBooks
  | StartFetchingBooks
  | FinishFetchingBooks
  | AddBook
  | RemoveBook;
