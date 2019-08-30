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
import { matchPath } from 'react-router-dom';
import { Location } from 'history';

// path definitions
export const homePathDef = '/';
export const notePathDef = '/notes/:noteUUID';
export const noteEditPathDef = '/notes/:noteUUID/edit';
export const noteNewPathDef = '/new';
export const booksPathDef = '/books';
export const loginPathDef = '/login';
export const joinPathDef = '/join';

// filterSearchObj filters the given search object and returns a new object
function filterSearchObj(obj) {
  const ret: any = {};

  const keys = Object.keys(obj);
  for (let i = 0; i < keys.length; ++i) {
    const key = keys[i];
    const val = obj[key];

    // reject empty string
    if (val !== '') {
      ret[key] = val;
    }
  }

  // page is implicitly 1
  if (ret.page === 1) {
    delete ret.page;
  }

  return ret;
}

interface GetLocationParams {
  pathname: string;
  searchObj?: any;
  state?: any;
}

function getLocation({
  pathname,
  searchObj,
  state
}: GetLocationParams): Location<any> {
  const ret: Location<any> = { pathname, search: '', state, hash: '' };

  if (searchObj) {
    const o = filterSearchObj(searchObj);

    ret.search = qs.stringify(o, { arrayFormat: 'repeat' });
  }
  if (state) {
    ret.state = state;
  }

  return ret;
}

export function getNewPath(searchObj = {}): Location {
  return getLocation({ pathname: '/new', searchObj });
}

export function getRandomPath(searchObj = {}): Location {
  return getLocation({ pathname: '/random', searchObj });
}

export function getHomePath(searchObj = {}): Location {
  return getLocation({ pathname: homePathDef, searchObj });
}

export function getBooksPath(searchObj = {}): Location {
  const basePath = '/books';

  return getLocation({ pathname: basePath, searchObj });
}

export function getNotePath(noteUUID: string, searchObj = {}): Location {
  const path = `/notes/${noteUUID}`;

  return getLocation({
    pathname: path,
    searchObj
  });
}

export function getNoteEditPath(noteUUID: string): Location {
  const path = `/notes/${noteUUID}/edit`;

  return getLocation({
    pathname: path
  });
}

export function getJoinPath(searchObj): Location {
  return getLocation({ pathname: '/join', searchObj });
}

export function getLoginPath(searchObj): Location {
  return getLocation({ pathname: '/login', searchObj });
}

export function getSubscriptionPath(searchObj = {}): Location {
  return getLocation({ pathname: '/subscriptions', searchObj });
}

export function getSubscriptionCheckoutPath(searchObj = {}): Location {
  return getLocation({ pathname: '/subscriptions/checkout', searchObj });
}

export enum SettingSections {
  account = 'account',
  notification = 'notification',
  billing = 'billing'
}

export function getSettingsPath(section: SettingSections) {
  return `/settings/${section}`;
}

// checkCurrentPath checks if the current path is the given path
export function checkCurrentPath(location: Location, path: string): boolean {
  const match = matchPath(location.pathname, {
    path,
    exact: true
  });

  return Boolean(match);
}
