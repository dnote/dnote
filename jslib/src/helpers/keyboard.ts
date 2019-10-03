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

export const KEYCODE_DOWN = 40;
export const KEYCODE_UP = 38;
export const KEYCODE_ENTER = 13;
export const KEYCODE_ESC = 27;
export const KEYCODE_TAB = 9;
export const KEYCODE_BACKSPACE = 8;

// alphabet
export const KEYCODE_LOWERCASE_B = 66;

// isPrintableKey returns if the key represented in the given event is printable.
export function isPrintableKey(e: KeyboardEvent) {
  return e.key.length === 1;
}
