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

// module https.ts provides an interface to make HTTP requests and receive responses

class ResponseError extends Error {
  response: Response;
}

function checkStatus(response: Response): Response | Promise<Response> {
  if (response.status >= 200 && response.status < 300) {
    return response;
  }

  return response.text().then(body => {
    const error = new ResponseError(body);
    error.response = response;

    throw error;
  });
}

function parseJSON<T>(response: Response): Promise<T> {
  if (response.headers.get('Content-Type') === 'application/json') {
    return response.json() as Promise<T>;
  }

  return Promise.resolve(null);
}

function request<T>(path: string, options: RequestInit) {
  return fetch(path, {
    ...options
  })
    .then(checkStatus)
    .then(res => {
      return parseJSON<T>(res);
    });
}

function get<T>(path: string, options = {}): Promise<T> {
  return request<T>(path, {
    method: 'GET',
    ...options
  });
}

function post<T>(path: string, data: any, options = {}): Promise<T> {
  return request<T>(path, {
    method: 'POST',
    body: JSON.stringify(data),
    ...options
  });
}

function patch<T>(path: string, data: any, options = {}) {
  return request<T>(path, {
    method: 'PATCH',
    body: JSON.stringify(data),
    ...options
  });
}

function put(path: string, data: any, options = {}) {
  return request(path, {
    method: 'PUT',
    body: JSON.stringify(data),
    ...options
  });
}

function del(path: string, options = {}) {
  return request(path, {
    method: 'DELETE',
    ...options
  });
}

export interface HttpClientConfig {
  pathPrefix: string;
  baseUrl: string;
}

// getHttpClient returns an http client
export function getHttpClient(c: HttpClientConfig) {
  function transformPath(path: string): string {
    let ret = path;

    if (c.pathPrefix !== '') {
      ret = `${c.pathPrefix}${ret}`;
    }
    if (c.baseUrl !== '') {
      ret = `${c.baseUrl}${ret}`;
    }

    return ret;
  }

  return {
    get: <T = any>(path: string, options = {}) => {
      return get<T>(transformPath(path), options);
    },
    post: <T>(path: string, data = {}, options = {}) => {
      return post<T>(transformPath(path), data, options);
    },
    patch: <T = any>(path: string, data, options = {}) => {
      return patch<T>(transformPath(path), data, options);
    },
    put: (path: string, data, options = {}) => {
      return put(transformPath(path), data, options);
    },
    del: (path: string, options = {}) => {
      return del(transformPath(path), options);
    }
  };
}
