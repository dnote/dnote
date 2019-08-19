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

import { parseSearchString } from './url';
import { Queries } from './queries';

export interface Filters {
  queries: Queries;
  page: number;
}

function compareBookArr(a1: string[], a2: string[]) {
  if (a1.length !== a2.length) {
    return false;
  }

  const a1Sorted = a1.sort();
  const a2Sorted = a2.sort();

  for (let i = 0; i < a1Sorted.length; ++i) {
    if (a1Sorted[i] !== a2Sorted[i]) {
      return false;
    }
  }

  return true;
}

// getFiltersFromSearchStr unmarshals the given search string from the URL
// into an object
export function getFiltersFromSearchStr(search: string): Filters {
  const searchObj = parseSearchString(search);

  let bookVal;
  if (typeof searchObj.book === 'string') {
    bookVal = [searchObj.book];
  } else if (searchObj.book === undefined) {
    bookVal = [];
  } else {
    bookVal = searchObj.book;
  }

  const ret: Filters = {
    queries: {
      q: searchObj.q || '',
      book: bookVal
    },
    page: parseInt(searchObj.page, 10) || 1
  };

  return ret;
}

// checkFilterEqual checks that the two given filters are equal
export function checkFilterEqual(a: Filters, b: Filters): boolean {
  return (
    a.page === b.page &&
    a.queries.q === b.queries.q &&
    compareBookArr(a.queries.book, b.queries.book)
  );
}

// toSearchObj transforms the filters into a search obj to be marshaled to a URL search string
export function toSearchObj(filters: Filters): any {
  const ret: any = {};

  const { queries } = filters;

  if (filters.page) {
    ret.page = filters.page;
  }
  if (queries.q !== '') {
    ret.q = queries.q;
  }
  if (queries.book.length > 0) {
    ret.book = queries.book;
  }

  return ret;
}
