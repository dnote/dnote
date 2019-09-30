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

import { parseSearchString } from 'jslib/helpers/url';
import { removeKey } from 'jslib/helpers/obj';
import { Queries, keywordBook } from 'jslib/helpers/queries';
import { getHomePath } from './paths';

export function getSearchDest(location: Location, queries: Queries) {
  let searchObj: any = parseSearchString(location.search);

  if (queries.q !== '') {
    searchObj.q = queries.q;
  } else {
    searchObj = removeKey(searchObj, 'q');
  }

  if (queries.book.length > 0) {
    searchObj.book = queries.book;
  } else {
    searchObj = removeKey(searchObj, keywordBook);
  }

  searchObj = removeKey(searchObj, 'page');

  return getHomePath(searchObj);
}
