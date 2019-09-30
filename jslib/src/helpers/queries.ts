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

import { Location } from 'history';

import { parseSearchString } from './url';
import { removeKey } from './obj';
import * as searchLib from './search';

export interface Queries {
  q: string;
  book: string[];
}

function encodeQuery(keyword: string, value: string): string {
  return `${keyword}:${value}`;
}

export const keywordBook = 'book';

export const keywords = [keywordBook];

// parse unmarshals the given string represesntation of the queries into an object
export function parse(s: string): Queries {
  const result = searchLib.parse(s, keywords);

  let bookValue: string[];
  const { book } = result.filters;
  if (!book) {
    bookValue = [];
  } else if (typeof book === 'string') {
    bookValue = [book];
  } else {
    bookValue = book;
  }

  let qValue: string;
  if (result.text) {
    qValue = result.text;
  } else {
    qValue = '';
  }

  return {
    q: qValue,
    book: bookValue
  };
}

// stringify marshals the givne queries into a string format
export function stringify(queries: Queries): string {
  let ret = '';

  if (queries.book.length > 0) {
    for (let i = 0; i < queries.book.length; i++) {
      const book = queries.book[i];

      ret += encodeQuery(keywordBook, book);
      ret += ' ';
    }
  }

  if (queries.q !== '') {
    const result = searchLib.parse(queries.q, keywords);
    ret += result.text;
  }

  return ret;
}
