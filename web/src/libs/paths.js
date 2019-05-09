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

function getPathObj({ pathname, searchObj, state }) {
  const ret = { pathname };

  if (searchObj) {
    const o = filterSearchObj(searchObj);

    ret.search = qs.stringify(o);
  }
  if (state) {
    ret.state = state;
  }

  return ret;
}

export function homePath(searchObj = {}, options = { demo: false }) {
  const { demo } = options;

  let basePath;
  if (demo) {
    basePath = '/demo';
  } else {
    basePath = '/';
  }

  return getPathObj({ pathname: basePath, searchObj });
}

export function booksPath(options = { demo: false }) {
  const { demo } = options;

  let basePath;
  if (demo) {
    basePath = '/demo/books';
  } else {
    basePath = '/books';
  }

  return getPathObj({ pathname: basePath });
}

export function dashboardPath(options = { demo: false }) {
  const { demo } = options;

  let basePath;
  if (demo) {
    basePath = '/demo/dashboard';
  } else {
    basePath = '/dashboard';
  }

  return getPathObj({ pathname: basePath });
}

export function notePath(noteUUID, searchObj, { demo, isEditor }) {
  const basePath = `/notes/${noteUUID}`;

  let path;
  if (demo) {
    path = `/demo${basePath}`;
  } else {
    path = basePath;
  }

  return getPathObj({
    pathname: path,
    searchObj,
    state: { editor: isEditor }
  });
}

export function clientsPath(client) {
  if (client) {
    return `/apps/${client}`;
  }

  return '/apps';
}

export function joinPath(searchObj) {
  return getPathObj({ pathname: '/join', searchObj });
}

export function loginPath(searchObj) {
  return getPathObj({ pathname: '/login', searchObj });
}

export function subscriptionsPath(searchObj) {
  return getPathObj({ pathname: '/subscriptions', searchObj });
}

export function settingsPath(section) {
  return `/settings/${section}`;
}

// mainSidebarPaths are paths that have the main sidebar and the main footer
export const mainSidebarPaths = [
  '/',
  '/books',
  '/dashboard',
  '/demo',
  '/demo/books',
  '/demo/dashboard'
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
  '/demo/dashboard',
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

  return !isSubscriptionsPath(pathname);
}
