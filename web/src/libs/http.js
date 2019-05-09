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

// module https.js provides an interface to make HTTP requests and receive responses

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

function request(path, options) {
  return fetch(path, {
    ...options
  })
    .then(checkStatus)
    .then(parseJSON);
}

function get(path, options = {}) {
  return request(path, {
    method: 'GET',
    ...options
  });
}

function post(path, data, options = {}) {
  return request(path, {
    method: 'POST',
    body: JSON.stringify(data),
    ...options
  });
}

function patch(path, data, options = {}) {
  return request(path, {
    method: 'PATCH',
    body: JSON.stringify(data),
    ...options
  });
}

function put(path, data, options = {}) {
  return request(path, {
    method: 'PUT',
    body: JSON.stringify(data),
    ...options
  });
}

function del(path, options = {}) {
  return request(path, {
    method: 'DELETE',
    ...options
  });
}

// httpClient implements basic http verbs
export const httpClient = {
  get,
  post,
  patch,
  put,
  del
};

function prependApi(path) {
  return `/api${path}`;
}

// apiClient is a special case of http client that prepends '/api' to the path
export const apiClient = {
  get: (path, options) => {
    return get(prependApi(path), options);
  },
  post: (path, options) => {
    return post(prependApi(path), options);
  },
  patch: (path, data, options) => {
    return patch(prependApi(path), data, options);
  },
  put: (path, data, options) => {
    return put(prependApi(path), data, options);
  },
  del: (path, options) => {
    return del(prependApi(path), options);
  }
};
