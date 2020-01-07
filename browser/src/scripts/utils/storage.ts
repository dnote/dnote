/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import ext from './ext';

const stateKey = 'state';

// filterState filters the given state to be suitable for reuse upon next app
// load
function filterState(state) {
  return {
    ...state,
    location: {
      ...state.location,
      path: '/'
    }
  };
}

function parseStorageItem(item) {
  if (!item) {
    return null;
  }

  return JSON.parse(item);
}

// saveState writes the given state to storage
export function saveState(state) {
  const filtered = filterState(state);
  const serialized = JSON.stringify(filtered);

  ext.storage.local.set({ [stateKey]: serialized }, () => {
    console.log('synced state');
  });
}

// loadState loads and parses serialized state stored in ext.storage
export function loadState(done) {
  ext.storage.local.get('state', items => {
    const parsed = {
      ...items,
      state: parseStorageItem(items.state)
    };

    return done(parsed);
  });
}
