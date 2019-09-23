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

import { getPath } from '../helpers/url';
import { getHttpClient, HttpClientConfig } from '../helpers/http';
import { NoteData } from '../operations/types';
import { Filters } from '../helpers/filters';

export interface CreateParams {
  book_uuid: string;
  content: string;
}

export interface CreateResponse {
  result: NoteData;
}

export interface UpdateParams {
  book_uuid?: string;
  content?: string;
  public?: boolean;
}

export interface UpdateNoteResp {
  status: number;
  result: NoteData;
}

export interface FetchResponse {
  notes: NoteData[];
  total: number;
}

export interface FetchOneQuery {
  q?: string;
}

type FetchOneResponse = NoteData;

export default function init(config: HttpClientConfig) {
  const client = getHttpClient(config);

  return {
    create: (params: CreateParams, opts = {}): Promise<CreateResponse> => {
      return client.post<CreateResponse>('/v3/notes', params, opts);
    },

    update: (noteUUID: string, params: UpdateParams) => {
      const endpoint = `/v3/notes/${noteUUID}`;

      return client.patch<UpdateNoteResp>(endpoint, params);
    },

    remove: (noteUUID: string) => {
      const endpoint = `/v3/notes/${noteUUID}`;

      return client.del(endpoint, {});
    },

    fetch: (filters: Filters) => {
      const params: any = {
        page: filters.page
      };

      const { queries } = filters;
      if (queries.q) {
        params.q = queries.q;
      }
      if (queries.book) {
        params.book = queries.book;
      }

      const endpoint = getPath('/notes', params);

      return client.get<FetchResponse>(endpoint, {});
    },

    fetchOne: (
      noteUUID: string,
      params: FetchOneQuery
    ): Promise<FetchOneResponse> => {
      const endpoint = getPath(`/notes/${noteUUID}`, params);

      return client.get<FetchOneResponse>(endpoint, {});
    },

    classicFetch: () => {
      const endpoint = '/classic/notes';

      return client.get(endpoint, { credentials: 'include' });
    }
  };
}
