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

export const validFacets = ['book', 'q'];

// TODO: replace it with getFacetsFromSearchStr
export function getFacetsFromQueryObj(queryObj) {
  const ret = {};

  validFacets.forEach(key => {
    ret[key] = queryObj[key];
  });

  return ret;
}

export function isFacetActive(searchStr) {
  const query = parseSearchString(searchStr);

  for (let i = 0; i < validFacets.length; i++) {
    const key = validFacets[i];

    if (query[key] !== undefined) {
      return true;
    }
  }

  return false;
}

export function getFacetsFromSearchStr(searchStr) {
  const searchObj = parseSearchString(searchStr);

  return getFacetsFromQueryObj(searchObj);
}
