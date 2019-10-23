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

import { pluralize } from '../../libs/string';

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

const shortDayNames = ['Sun', 'Mon', 'Tue', 'Wed', 'Thur', 'Fri', 'Sat'];

const fullDayNames = [
  'Sunday',
  'Monday',
  'Tuesday',
  'Wednesday',
  'Thursday',
  'Friday',
  'Saturday'
];

/******* durations in milliseconds */
export const SECOND = 1000;
export const MINUTE = 60 * SECOND;
export const HOUR = 60 * MINUTE;
export const DAY = 24 * HOUR;
export const WEEK = 7 * DAY;

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

// getDayName returns the shortened month name of the given date
export function getDayName(date: Date, short: boolean = false) {
  const day = date.getDay();

  if (short) {
    return shortDayNames[day];
  }

  return fullDayNames[day];
}

// getMonthName returns the shortened month name of the given date
export function getMonthName(date: Date, short: boolean = false) {
  const month = date.getMonth();

  if (short) {
    return shortMonthNames[month];
  }

  return monthNumToFullName(month + 1);
}

// monthNumToFullName returns a full month name based on the number denoting the month,
// ranging from 1 to 12 corresponding to each month of a year.
export function monthNumToFullName(num: number): string {
  if (num > 12 || num < 1) {
    throw new Error(`invalid month number ${num}`);
  }

  return fullMonthNames[num - 1];
}

export function pad(value: number): string {
  return value < 10 ? `0${value}` : `${value}`;
}

// getUTCOffset returns the UTC offset string for the client. The returned
// value is in the format of '+08:00'
export function getUTCOffset(): string {
  const date = new Date();

  let sign;
  if (date.getTimezoneOffset() > 0) {
    sign = '-';
  } else {
    sign = '+';
  }

  const offset = Math.abs(date.getTimezoneOffset());
  const hours = Math.floor(offset / 60);

  const rawMinutes = offset % 60;
  if (rawMinutes === 0) {
    return `${sign}${hours}`;
  }

  const minutes = pad(rawMinutes);
  return `${sign}${hours}:${minutes}`;
}

// daysToMs translates the given number of days to seconds
export function daysToMs(numDays: number) {
  const dayInSeconds = DAY;

  return dayInSeconds * numDays;
}

function parseMs(ms: number) {
  const weeks = Math.floor(ms / WEEK);
  const days = Math.floor((ms % WEEK) / DAY);
  const hours = Math.floor(((ms % WEEK) % DAY) / HOUR);
  const minutes = Math.floor((((ms % WEEK) % DAY) % HOUR) / MINUTE);
  const seconds = (((ms % WEEK) % DAY) % HOUR) % MINUTE;

  return {
    weeks,
    days,
    hours,
    minutes,
    seconds
  };
}

// msToHTMLTimeDuration converts the given number of seconds into a valid
// time duration string as defined by the W3C HTML5 recommendation
export function msToHTMLTimeDuration(ms: number): string {
  const { weeks, days, hours, minutes, seconds } = parseMs(ms);

  let ret = 'P';

  const numDays = weeks * 7 + days;
  if (numDays > 0) {
    ret += `${numDays}D`;
  }

  if (hours > 0) {
    ret += `${hours}H`;
  }
  if (minutes > 0) {
    ret += `${minutes}M`;
  }
  if (seconds > 0) {
    ret += `${seconds}S`;
  }

  return ret;
}

// msToDuration translates the given time in seconds into a human-readable duration
export function msToDuration(ms: number): string {
  const { weeks, days, hours, minutes, seconds } = parseMs(ms);

  let ret = '';

  if (weeks > 0) {
    ret += `${weeks} ${pluralize('week', weeks)} `;
  }
  if (days > 0) {
    ret += `${days} ${pluralize('day', days)} `;
  }
  if (hours > 0) {
    ret += `${hours} ${pluralize('hour', hours)} `;
  }
  if (minutes > 0) {
    ret += `${minutes} ${pluralize('minute', minutes)} `;
  }

  return ret.trim();
}

export function timeAgo(ms: number, simple: boolean = false): string {
  const shortNounMap = {
    year: 'y',
    week: 'w',
    month: 'm',
    day: 'd',
    hour: 'h',
    minute: 'min'
  };

  function getStr(interval: number, noun: string): string {
    // if (simple) {
    //   const n = shortNounMap[noun];
    //   return `${interval}${n}`;
    // }

    return `${interval} ${pluralize(noun, interval)} ago`;
  }

  const ts = Math.floor(new Date().getTime() - ms);

  let interval = Math.floor(ts / (52 * WEEK));
  if (interval > 1) {
    return getStr(interval, 'year');
  }

  interval = Math.floor(ts / (4 * WEEK));
  if (interval >= 1) {
    return getStr(interval, 'month');
  }

  interval = Math.floor(ts / WEEK);
  if (interval >= 1) {
    return getStr(interval, 'week');
  }

  interval = Math.floor(ts / DAY);
  if (interval >= 1) {
    return getStr(interval, 'day');
  }

  interval = Math.floor(ts / HOUR);
  if (interval >= 1) {
    return getStr(interval, 'hour');
  }

  interval = Math.floor(ts / MINUTE);
  if (interval >= 1) {
    return getStr(interval, 'minute');
  }

  return 'Just now';
}
