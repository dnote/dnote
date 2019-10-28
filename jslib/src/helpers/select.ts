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

// Option represents an option in a selection list
export interface Option {
  label: string;
  value: string;
}

// optionValueCreate is the value of the option for creating a new option
export const optionValueCreate = 'create-new-option';

// filterOptions returns a new array of options based on the given filter criteria
export function filterOptions(
  options: Option[],
  term: string,
  creatable: boolean
): Option[] {
  if (!term) {
    return options;
  }

  const ret = [];
  const searchReg = new RegExp(`${term}`, 'i');
  let hit = null;

  for (let i = 0; i < options.length; i++) {
    const option = options[i];

    if (option.label === term) {
      hit = option;
    } else if (searchReg.test(option.label) && option.value !== '') {
      ret.push(option);
    }
  }

  // if there is an exact match, display the option at the top
  // otherwise, display a creatable option at the bottom
  if (hit) {
    ret.unshift(hit);
  } else if (creatable) {
    // creatable option has a value of an empty string
    ret.push({ label: term, value: '' });
  }

  return ret;
}

// booksToOptions returns an array of options for select ui, given an array of books
export function booksToOptions(books: BookData[]): Option[] {
  const ret = [];

  for (let i = 0; i < books.length; ++i) {
    const book = books[i];

    ret.push({
      label: book.label,
      value: book.uuid
    });
  }

  return ret;
}
