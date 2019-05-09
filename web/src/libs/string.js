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

// excerpt trims the given string up to the last word that makes the string
// exceed the maxLength, and attaches ellipses at the end. If the string is
// shorter than the given maxLength, it returns the original string.
export function excerpt(str, maxLength) {
  if (str.length < maxLength) {
    return str;
  }

  let ret;

  ret = str.substr(0, maxLength);
  ret = ret.substr(0, Math.min(ret.length, ret.lastIndexOf(' ')));
  ret += '...';

  return ret;
}

// escapesRegExp escapes the regular expression special characters.
export function escapesRegExp(s) {
  return s.replace(/[-/\\^$*+?.()|[\]{}]/g, '\\$&');
}
