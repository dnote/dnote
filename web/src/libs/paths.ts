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

// filterSearchObj filters the given search object and returns a new object
function filterSearchObj(obj) {
  const ret = {};

  const keys = Object.keys(obj);
  for (let i = 0; i < keys.length; ++i) {
    const key = keys[i];
    const val = obj[key];

    // reject empty string
    if (val !== '') {
      ret[key] = val;
    }
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

    ret.search = qs.stringify(o);
  }
  if (state) {
    ret.state = state;
  }

  return ret;
}

export function getNewPath() {
  return getLocation({ pathname: '/new' });
}

export function getRandomPath() {
  return getLocation({ pathname: '/random' });
}

export function getHomePath(searchObj = {}, options = { demo: false }) {
  const { demo } = options;

  let basePath;
  if (demo) {
    basePath = '/demo';
  } else {
    basePath = '/';
  }

  return getLocation({ pathname: basePath, searchObj });
}

export function getBooksPath(options = { demo: false }) {
  const { demo } = options;

  let basePath;
  if (demo) {
    basePath = '/demo/books';
  } else {
    basePath = '/books';
  }

  return getLocation({ pathname: basePath });
}

export function getDigestsPath(options = { demo: false }) {
  const { demo } = options;

  let basePath;
  if (demo) {
    basePath = '/demo/digests';
  } else {
    basePath = '/digests';
  }

  return getLocation({ pathname: basePath });
}

export function getDigestPath(digestUUID, options = { demo: false }) {
  const { demo } = options;

  let basePath;
  if (demo) {
    basePath = '/demo/digests';
  } else {
    basePath = '/digests';
  }

  const path = `${basePath}/${digestUUID}`;
  return getLocation({ pathname: path });
}

export function getNotePath(noteUUID: string) {
  const path = `/notes/${noteUUID}`;

  return getLocation({
    pathname: path
  });
}

export function getNoteEditPath(noteUUID: string) {
  const path = `/notes/${noteUUID}/edit`;

  return getLocation({
    pathname: path
  });
}

export function getJoinPath(searchObj) {
  return getLocation({ pathname: '/join', searchObj });
}

export function getLoginPath(searchObj) {
  return getLocation({ pathname: '/login', searchObj });
}

export function getSubscriptionPath(searchObj) {
  return getLocation({ pathname: '/subscriptions', searchObj });
}

export function getSubscriptionCheckoutPath(searchObj) {
  return getLocation({ pathname: '/subscriptions/checkout', searchObj });
}

export function getSettingsPath(section) {
  return `/settings/${section}`;
}

// mainSidebarPaths are paths that have the main sidebar and the main footer
export const mainSidebarPaths = [
  '/',
  '/books',
  '/digests',
  '/demo',
  '/demo/books',
  '/demo/digests'
];

// noteSidebarPaths are paths that have the note sidebar
export const noteSidebarPaths = ['/', '/demo'];

// footerPaths are paths that have footers
export const footerPaths = [
  ...mainSidebarPaths,
  '/settings/:section',
  '/notes/:noteUUID',
  '/demo/notes/:noteUUID'
];

function isEmailPreferencePath(pathname) {
  const match = matchPath(pathname, {
    path: '/email-preference',
    exact: true
  });

  return Boolean(match);
}

// isNotePath checks if the given pathname is for the note path
export function isNotePath(pathname) {
  const match = matchPath(pathname, {
    path: ['/notes/:noteUUID', '/demo/notes/:noteUUID'],
    exact: true
  });

  return Boolean(match);
}

// isSubscriptionsPath checks if the given pathname is for the subscriptions path
export function isSubscriptionsPath(pathname) {
  const match = matchPath(pathname, {
    path: '/subscriptions',
    exact: true
  });

  return Boolean(match);
}

// isSubscriptionsCheckoutPath checks if the given pathname is for the subscriptions path
export function isSubscriptionsCheckoutPath(pathname) {
  const match = matchPath(pathname, {
    path: '/subscriptions/checkout',
    exact: true
  });

  return Boolean(match);
}

// isDigestPath checks if the given pathname is for the digest path
export function isDigestPath(pathname) {
  const match = matchPath(pathname, {
    path: ['/digests/:digestUUID', '/demo/digests/:digestUUID'],
    exact: true
  });

  return Boolean(match);
}

// isLegacyPath checks if the given pathname is for the legacy path
export function isLegacyPath(pathname) {
  const match = matchPath(pathname, {
    path: '/legacy'
  });

  return Boolean(match);
}

// isHomePath checks if the given pathname is for the home path
export function isHomePath(pathname, demo = false) {
  let path;
  if (demo) {
    path = '/demo';
  } else {
    path = '/';
  }

  const match = matchPath(pathname, {
    path,
    exact: true
  });

  return Boolean(match);
}

const demoPaths = [
  '/demo',
  '/demo/books',
  '/demo/digests',
  '/demo/digests/:digestUUID',
  '/demo/notes/:noteUUID'
];

// isDemoPath checks if the given pathname is for the demo path
export function isDemoPath(pathname) {
  for (let i = 0; i < demoPaths.length; ++i) {
    const p = demoPaths[i];

    const match = matchPath(pathname, {
      path: p,
      exact: true
    });

    if (match) {
      return true;
    }
  }

  return false;
}

// checkBoxedLayout determines if the layout for the given location is boxed
export function checkBoxedLayout(location, isEditor) {
  const { pathname } = location;

  if (isNotePath(pathname)) {
    return isEditor;
  }
  if (isLegacyPath(pathname)) {
    return false;
  }
  if (isEmailPreferencePath(pathname)) {
    return false;
  }
  if (isDigestPath(pathname)) {
    return false;
  }
  if (isSubscriptionsCheckoutPath(pathname)) {
    return false;
  }

  return !isSubscriptionsPath(pathname);
}

// path definitions
export const notePath = '/notes/:noteUUID';
