/* Copyright (C) 2019, 2020, 2021 Monomax Software Pty Ltd
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

export interface PresentedNote {
  uuid: string;
  content: string;
  updated_at: string;
  created_at: string;
  user: {
    name: '';
    uuid: '';
  };
  public: boolean;
  book: {
    label: '';
    uuid: '';
  };
  usn: number;
  added_on: number;
}

export function mapNote(item: PresentedNote): NoteData {
  return {
    uuid: item.uuid,
    content: item.content,
    createdAt: item.created_at,
    updatedAt: item.updated_at,
    public: item.public,
    user: {
      name: item.user.name,
      uuid: item.user.uuid
    },
    book: {
      label: item.book.label,
      uuid: item.book.uuid
    },
    usn: item.usn,
    addedOn: item.added_on
  };
}

export interface CreateParams {
  book_uuid: string;
  content: string;
}

export interface CreateResponse {
  result: PresentedNote;
}

export interface CreateResult {
  result: NoteData;
}

export interface UpdateParams {
  book_uuid?: string;
  content?: string;
  public?: boolean;
}

export interface UpdateResponse {
  status: number;
  result: PresentedNote;
}

export interface UpdateResult {
  status: number;
  result: NoteData;
}

export interface FetchResponse {
  notes: PresentedNote[];
  total: number;
}

export interface FetchResult {
  notes: NoteData[];
  total: number;
}

export interface FetchOneQuery {
  q?: string;
}

type FetchOneResponse = PresentedNote;
type FetchOneResult = NoteData;

export default function init(config: HttpClientConfig) {
  const client = getHttpClient(config);

  return {
    create: (params: CreateParams, opts = {}): Promise<CreateResult> => {
      return client
        .post<CreateResponse>('/v3/notes', params, opts)
        .then(res => {
          return {
            result: mapNote(res.result)
          };
        });
    },

    update: (noteUUID: string, params: UpdateParams): Promise<UpdateResult> => {
      const endpoint = `/v3/notes/${noteUUID}`;

      return client.patch<UpdateResponse>(endpoint, params).then(res => {
        return {
          status: res.status,
          result: mapNote(res.result)
        };
      });
    },

    remove: (noteUUID: string) => {
      const endpoint = `/v3/notes/${noteUUID}`;

      return client.del(endpoint, {});
    },

    fetch: (filters: Filters): Promise<FetchResult> => {
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

      return client.get<FetchResponse>(endpoint, {}).then(res => {
        return {
          total: res.total,
          notes: res.notes.map(mapNote)
        };
      });
    },

    fetchOne: (
      noteUUID: string,
      params: FetchOneQuery
    ): Promise<FetchOneResult> => {
      const endpoint = getPath(`/notes/${noteUUID}`, params);

      return client.get<FetchOneResponse>(endpoint, {}).then(mapNote);
    },

    classicFetch: () => {
      const endpoint = '/classic/notes';

      return client.get(endpoint, { credentials: 'include' });
    }
  };
}
