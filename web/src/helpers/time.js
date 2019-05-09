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

import moment from 'moment';

const shortMonthNames = [
  'Jan',
  'Feb',
  'Mar',
  'Apr',
  'May',
  'Jun',
  'Jul',
  'Aug',
  'Sep',
  'Oct',
  'Nov',
  'Dec'
];

/** ***** durations in milliseconds */
const DAY = 86400000;

// nanosecToSec converts a given nanoseconds to seconds by dropping surplus digits
export function nanosecToSec(t) {
  const truncated = String(t).slice(0, -9);

  return parseInt(truncated, 10);
}

// nanosecToMillisec converts a given nanoseconds to milliseconds by dropping surplus digits
export function nanosecToMillisec(t) {
  const truncated = String(t).slice(0, -6);

  return parseInt(truncated, 10);
}

// getShortMonthName returns the shortened month name of the given date
export function getShortMonthName(date) {
  const month = date.getMonth();

  return shortMonthNames[month];
}

// presentNoteTS presents a note's added_on timestamp which is in unix nano
export function presentNoteTS(t) {
  const time = nanosecToSec(t);
  const past = moment.unix(time);

  const now = new Date();
  const diff = -past.diff(now);

  if (diff < DAY) {
    return `today ${past.format('h:mm a')}`;
  }

  if (diff < 2 * DAY) {
    return `yesterday ${past.format('h:mm a')}`;
  }

  if (diff < 7 * DAY) {
    return past.format('dddd h:mm a');
  }

  return `${past.format('MMM D')} (${past.fromNow()})`;
}
