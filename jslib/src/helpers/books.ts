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

import { BookData } from '../operations/types';

// errBookNameNumeric is an error for book names that only contain numbers
export const errBookNameNumeric = new Error(
  'The book name cannot contain only numbers'
);

// errBookNameHasSpace is an error for book names that have any space
export const errBookNameHasSpace = new Error(
  'The book name cannot contain spaces'
);

// errBookNameHasComma is an error for book names that have any comma
export const errBookNameHasComma = new Error('The book name has comma');

// errBookNameReserved is an error incidating that the specified book name is reserved
export const errBookNameReserved = new Error('The book name is reserved');

const numberRegex = /^\d+$/;

const reservedBookNames = ['trash', 'conflicts'];

// validateBookName validates the given book name and throws error if not valid
export function validateBookName(bookName: string) {
  if (reservedBookNames.indexOf(bookName) > -1) {
    throw errBookNameReserved;
  }

  if (numberRegex.test(bookName)) {
    throw errBookNameNumeric;
  }

  if (bookName.indexOf(' ') > -1) {
    throw errBookNameHasSpace;
  }

  if (bookName.indexOf(',') > -1) {
    throw errBookNameHasComma;
  }
}

// checkDuplicate checks if the given book name has a duplicate in the given array
// of books
export function checkDuplicate(books: BookData[], bookName: string): boolean {
  for (let i = 0; i < books.length; i++) {
    const book = books[i];

    if (book.label === bookName) {
      return true;
    }
  }

  return false;
}
