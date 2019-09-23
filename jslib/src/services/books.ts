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

import qs from 'qs';
import { getHttpClient, HttpClientConfig } from '../helpers/http';

export interface BookFetchParams {
  name?: string;
  encrypted?: boolean;
}

export interface CreateParams {
  name: string;
}

export interface CreatePayload {
  book: {
    uuid: string;
    usn: number;
    created_at: string;
    updated_at: string;
    label: string;
  };
}

// TODO: type
type updateParams = any;

export default function init(config: HttpClientConfig) {
  const client = getHttpClient(config);

  return {
    fetch: (queryObj: BookFetchParams = {}, opts = {}) => {
      const baseURL = '/v3/books';

      const queryStr = qs.stringify(queryObj);

      let endpoint;
      if (queryStr) {
        endpoint = `${baseURL}?${queryStr}`;
      } else {
        endpoint = baseURL;
      }

      return client.get(endpoint, opts);
    },

    create: (payload: CreateParams, opts = {}) => {
      return client.post<CreatePayload>('/v3/books', payload, opts);
    },

    remove: (uuid: string) => {
      return client.del(`/v3/books/${uuid}`);
    },

    update: (uuid: string, payload: updateParams) => {
      return client.patch(`/v3/books/${uuid}`, payload);
    },

    get: (bookUUID: string) => {
      const endpoint = `/v3/books/${bookUUID}`;

      return client.get(endpoint);
    }
  };
}
