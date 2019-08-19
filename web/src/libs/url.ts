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

import qs from 'qs';
import isArray from 'lodash/isArray';
import omitBy from 'lodash/omitBy';
import { Location } from 'history';

// getPath returns a path optionally suffixed by query string
export function getPath(path, queryObj): string {
  const queryStr = qs.stringify(queryObj, { arrayFormat: 'repeat' });

  if (!queryStr) {
    return path;
  }

  return `${path}?${queryStr}`;
}

// getPathFromLocation returns a full path based on the location object used by
// React Router
export function getPathFromLocation(location): string {
  const { pathname, search } = location;

  return `${pathname}${search}`;
}

// parseSearchString parses the 'search' string in `location` object provided
// by React Router.
export function parseSearchString(search: string): any {
  if (!search || search === '') {
    return {};
  }

  // drop the leading '?'
  const queryStr = search.substring(1);
  return qs.parse(queryStr);
}

// addQueryToLocation returns a new location object for react-router given the
// new `queryKey` and `val` to be set in loation.query.
// If there exists the given key in the query object, addQueryToLocation sets its
// value to be an array containing the old value and the new value.
// Otherwise the value for the key is set to the `val`.
export function addQueryToLocation(
  location: Location,
  queryKey: string,
  val: string,
  override = true
): Location {
  const queryObj = parseSearchString(location.search);
  const existingParam = queryObj[queryKey];

  let updatedQueryVal;
  if (existingParam && !override) {
    if (isArray(existingParam)) {
      updatedQueryVal = [...existingParam, val];
    } else {
      updatedQueryVal = [existingParam, val];
    }
  } else {
    updatedQueryVal = val;
  }

  const newQueryObj = {
    ...queryObj,
    [queryKey]: updatedQueryVal
  };

  return {
    ...location,
    search: `?${qs.stringify(newQueryObj)}`
  };
}

// removeQueryFromLocation returns a new location object without the queryKey
// and val
export function removeQueryFromLocation(
  location: Location,
  queryKey: string,
  val?: string
): Location {
  const queryObj = parseSearchString(location.search);
  const existingParam = queryObj[queryKey];
  if (!existingParam) {
    return location;
  }

  let newQueryObj;
  if (val === undefined) {
    newQueryObj = omitBy(queryObj, (v, k) => k === queryKey);
  } else {
    const queryVal = val.toString(); // stringify because query params only store string

    if (isArray(existingParam)) {
      const updatedQueryVal = existingParam.filter(elm => elm !== queryVal);
      newQueryObj = {
        ...queryObj,
        [queryKey]: updatedQueryVal
      };
    } else {
      newQueryObj = omitBy(
        queryObj,
        (v, k) => k === queryKey && queryVal === v
      );
    }
  }

  return {
    ...location,
    search: `?${qs.stringify(newQueryObj)}`
  };
}

export function getReferrer(location: Location): string {
  const queryObj = parseSearchString(location.search);
  const { referrer } = queryObj;

  if (referrer) {
    return decodeURIComponent(referrer);
  }

  if (location.state && location.state.referrer) {
    return location.state.referrer;
  }

  return '';
}
