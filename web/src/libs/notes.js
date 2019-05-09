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

import { nanosecToMillisec } from '../helpers/time';

// parsePrevDate returns an object containing year and month for previous note group
export function parsePrevDate(ts) {
  const ms = nanosecToMillisec(ts);
  const date = new Date(ms);

  const year = date.getUTCFullYear();
  const month = date.getUTCMonth() + 1;

  return { year, month };
}
