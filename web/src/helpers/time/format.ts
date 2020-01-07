/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

import { getMonthName, getUTCOffset, pad, getDayName } from './index';
import { addOrdinalSuffix } from '../../libs/string';

// format verbs
const YYYY = '%YYYY';
const YY = '%YY';
const MMMM = '%MMMM';
const MMM = '%MMM';
const MM = '%MM';
const M = '%M';
const DD = '%DD';
const D = '%D';
const Do = '%Do';
const hh = '%hh';
const h = '%h';
const mm = '%mm';
const m = '%m';
const A = '%A';
const a = '%a';
const Z = '%Z';
const dddd = '%dddd';

// getPeriod returns the period for the time for the given date
function getPeriod(date: Date) {
  const hour = date.getHours();

  let ret;
  if (hour > 12) {
    ret = 'PM';
  } else {
    ret = 'AM';
  }

  return ret;
}

// formatTime formats time to a human readable string based on the given format string
export default function formatTime(date: Date, format: string): string {
  let ret = format;

  if (ret.indexOf(YYYY) > -1) {
    ret = ret.replace(new RegExp(YYYY, 'g'), date.getFullYear().toString());
  }
  if (ret.indexOf(YY) > -1) {
    const year = date.getFullYear().toString();
    const newSubstr = year.substring(2, 4);

    ret = ret.replace(new RegExp(YY, 'g'), newSubstr);
  }

  if (ret.indexOf(MMMM) > -1) {
    ret = ret.replace(new RegExp(MMMM, 'g'), getMonthName(date));
  }
  if (ret.indexOf(MMM) > -1) {
    ret = ret.replace(new RegExp(MMM, 'g'), getMonthName(date, true));
  }
  if (ret.indexOf(MM) > -1) {
    const monthIdx = date.getMonth() + 1;
    const newSubstr = pad(monthIdx);

    ret = ret.replace(new RegExp(MM, 'g'), newSubstr);
  }
  if (ret.indexOf(M) > -1) {
    const monthIdx = `${date.getMonth() + 1}`;

    ret = ret.replace(new RegExp(M, 'g'), monthIdx);
  }

  if (ret.indexOf(DD) > -1) {
    const day = date.getDate();
    const newSubstr = pad(day);

    ret = ret.replace(new RegExp(DD, 'g'), newSubstr);
  }
  if (ret.indexOf(Do) > -1) {
    const day = date.getDate();
    const newSubstr = addOrdinalSuffix(day);

    ret = ret.replace(new RegExp(Do, 'g'), newSubstr);
  }
  if (ret.indexOf(D) > -1) {
    ret = ret.replace(new RegExp(D, 'g'), date.getDate().toString());
  }

  if (ret.indexOf(hh) > -1) {
    const hour = date.getHours();

    ret = ret.replace(new RegExp(hh, 'g'), pad(hour));
  }
  if (ret.indexOf(h) > -1) {
    let hour = date.getHours();
    if (hour > 12) {
      hour -= 12;
    }

    ret = ret.replace(new RegExp(h, 'g'), hour.toString());
  }

  if (ret.indexOf(mm) > -1) {
    const minute = date.getMinutes();

    ret = ret.replace(new RegExp(mm, 'g'), pad(minute));
  }
  if (ret.indexOf(m) > -1) {
    ret = ret.replace(/m/g, date.getMinutes().toString());
  }

  if (ret.indexOf(A) > -1) {
    const period = getPeriod(date);

    ret = ret.replace(new RegExp(A, 'g'), period);
  }
  if (ret.indexOf(a) > -1) {
    const period = getPeriod(date).toLowerCase();

    ret = ret.replace(new RegExp(a, 'g'), period);
  }

  if (ret.indexOf(dddd) > -1) {
    ret = ret.replace(new RegExp(dddd, 'g'), getDayName(date, false));
  }

  if (ret.indexOf(Z) > -1) {
    const offset = getUTCOffset();

    ret = ret.replace(new RegExp(Z, 'g'), offset);
  }

  return ret;
}
