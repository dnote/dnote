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

import initBooksService from '../services/books';
import { HttpClientConfig } from '../helpers/http';
import { BookData } from './types';

export interface CreateParams {
  name: string;
}

export default function init(c: HttpClientConfig) {
  const booksService = initBooksService(c);

  return {
    get: (bookUUID: string) => {
      return booksService.get(bookUUID);
    },

    // create creates an encrypted book. It returns a promise that resolves with
    // a decrypted book.
    create: (payload: CreateParams): Promise<BookData> => {
      return booksService.create(payload).then(res => {
        return res.book;
      });
    },

    fetch: (params = {}) => {
      return booksService.fetch(params);
    },

    // remove deletes the book with the given uuid
    remove: (bookUUID: string) => {
      return booksService.remove(bookUUID);
    }
  };
}
