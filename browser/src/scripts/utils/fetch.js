/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import qs from 'qs';

function checkStatus(response) {
  if (response.status >= 200 && response.status < 300) {
    return response;
  }
  return response.text().then(body => {
    const error = new Error(body);
    error.response = response;

    throw error;
  });
}

function parseJSON(response) {
  if (response.headers.get('Content-Type') === 'application/json') {
    return response.json();
  }

  return Promise.resolve();
}

function request(url, options) {
  return fetch(url, options).then(checkStatus).then(parseJSON);
}

export function post(url, data, options = {}) {
  return request(url, {
    method: 'POST',
    body: JSON.stringify(data),
    ...options
  });
}

export function get(url, options = {}) {
  let endpoint = url;

  if (options.params) {
    endpoint = `${endpoint}?${qs.stringify(options.params)}`;
  }

  return request(endpoint, {
    method: 'GET',
    ...options
  });
}
