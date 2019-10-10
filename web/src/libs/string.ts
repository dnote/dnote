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
export function excerpt(s: string, maxLength: number) {
  if (s.length < maxLength) {
    return s;
  }

  let ret;

  ret = s.substr(0, maxLength);
  ret = ret.substr(0, Math.min(ret.length, ret.lastIndexOf(' ')));
  ret += '...';

  return ret;
}

// escapesRegExp escapes the regular expression special characters.
export function escapesRegExp(s: string) {
  return s.replace(/[-/\\^$*+?.()|[\]{}]/g, '\\$&');
}

// pluralize pluralizes the given singular noun depending on the given count
// by naively concatenating a suffix 's'.
export function pluralize(singular: string, count: number): string {
  if (count === 1) {
    return singular;
  }

  return `${singular}s`;
}

// addOrdinalSuffix append appropriate suffix to the given number to
// represent it as an ordinal number
export function addOrdinalSuffix(i: number): string {
  const j = i % 10;
  const k = i % 100;

  if (j == 1 && k != 11) {
    return i + 'st';
  }
  if (j == 2 && k != 12) {
    return i + 'nd';
  }
  if (j == 3 && k != 13) {
    return i + 'rd';
  }
  return i + 'th';
}
