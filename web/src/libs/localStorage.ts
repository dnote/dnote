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

const stateKey = 'state';

// loadState parses the serialized state tree stored in the localStorage
// and returns it
export function loadState(): JSON {
  try {
    const serialized = localStorage.getItem(stateKey);

    if (serialized === null) {
      return undefined;
    }

    return JSON.parse(serialized);
  } catch (e) {
    console.log('Unable load state from the localStorage', e.message);
    return undefined;
  }
}

// saveState writes the given state to localStorage
export function saveState(state: object) {
  try {
    const serialized = JSON.stringify(state);

    localStorage.setItem(stateKey, serialized);
  } catch (e) {
    console.log('Unable save to the localStorage', e.message);
  }
}
