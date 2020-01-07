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

export type BookData = any;

export interface BooksState {
  items: BookData[];
  isFetching: boolean;
  isFetched: boolean;
  error: string | null;
}

export const START_FETCHING = 'books/START_FETCHING';
export const RECEIVE = 'books/RECEIVE';
export const RECEIVE_ERROR = 'books/RECEIVE_ERROR';

export interface StartFetchingAction {
  type: typeof START_FETCHING;
}

export interface ReceiveAction {
  type: typeof RECEIVE;
  data: {
    books: BookData[];
  };
}

export interface ReceiveErrorAction {
  type: typeof RECEIVE_ERROR;
  data: {
    error: string;
  };
}

export type BooksActionType =
  | StartFetchingAction
  | ReceiveAction
  | ReceiveErrorAction;
