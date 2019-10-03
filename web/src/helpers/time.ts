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

const fullMonthNames = [
  'January',
  'February',
  'March',
  'April',
  'May',
  'June',
  'July',
  'August',
  'September',
  'October',
  'November',
  'Dececember'
];

/******* durations in milliseconds */
export const DAY = 86400000;

// nanosecToSec converts a given nanoseconds to seconds by dropping surplus digits
export function nanosecToSec(t: number): number {
  const truncated = String(t).slice(0, -9);

  return parseInt(truncated, 10);
}

// nanosecToMillisec converts a given nanoseconds to milliseconds by dropping surplus digits
export function nanosecToMillisec(t: number): number {
  const truncated = String(t).slice(0, -6);

  return parseInt(truncated, 10);
}

// getShortMonthName returns the shortened month name of the given date
export function getShortMonthName(date: Date) {
  const month = date.getMonth();

  return shortMonthNames[month];
}

// monthNumToFullName returns a full month name based on the number denoting the month,
// ranging from 1 to 12 corresponding to each month of a year.
export function monthNumToFullName(num: number): string {
  if (num > 12 || num < 1) {
    throw new Error(`invalid month number ${num}`);
  }

  return fullMonthNames[num - 1];
}

// presentNoteTS presents a note's added_on timestamp which is in unix nano
export function presentNoteTS(t: number): string {
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

// getUTCOffset returns the UTC offset string for the client. The returned
// value is in the format of '+08:00'
export function getUTCOffset(): string {
  function pad(value: number): string {
    return value < 10 ? `0${value}` : `${value}`;
  }

  const date = new Date();

  let sign;
  if (date.getTimezoneOffset() > 0) {
    sign = '+';
  } else {
    sign = '-';
  }

  const offset = Math.abs(date.getTimezoneOffset());
  const hours = pad(Math.floor(offset / 60));
  const minutes = pad(offset % 60);

  return sign + hours + ':' + minutes;
}

// daysToSec translates the given number of days to seconds
export function daysToSec(numDays: number) {
  const dayInSeconds = DAY / 1000;

  return dayInSeconds * numDays;
}
