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

export function getBrowserType() {
  if (
    (navigator.userAgent.indexOf('Opera') ||
      navigator.userAgent.indexOf('OPR')) !== -1
  ) {
    return 'Opera';
  }
  if (navigator.userAgent.indexOf('Chrome') !== -1) {
    return 'Chrome';
  }
  if (navigator.userAgent.indexOf('Safari') !== -1) {
    return 'Safari';
  }
  if (navigator.userAgent.indexOf('Firefox') !== -1) {
    return 'Firefox';
  }
  if (
    navigator.userAgent.indexOf('MSIE') !== -1 ||
    !!document.documentMode === true
  ) {
    // IF IE > 10
    return 'IE';
  }

  return 'unknown';
}
