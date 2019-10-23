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

import { BookData, RepetitionRuleData, BookDomain } from '../operations/types';
import { getHttpClient, HttpClientConfig } from '../helpers/http';
import { getPath } from '../helpers/url';

export type FetchResponse = RepetitionRuleData[];

export interface CreateParams {
  title: string;
  hour: number;
  minute: number;
  book_domain: BookDomain;
  frequency: number;
  note_count: number;
  book_uuids: string[];
  enabled: boolean;
}

export type UpdateParams = Partial<CreateParams>;

export interface RepetitionRuleRespData {
  uuid: string;
  title: string;
  enabled: boolean;
  hour: number;
  minute: number;
  book_domain: BookDomain;
  frequency: number;
  books: BookData[];
  note_count: number;
  last_active: number;
  created_at: string;
  updated_at: string;
}

function mapData(d: RepetitionRuleRespData): RepetitionRuleData {
  return {
    uuid: d.uuid,
    title: d.title,
    enabled: d.enabled,
    hour: d.hour,
    minute: d.minute,
    bookDomain: d.book_domain,
    frequency: d.frequency,
    books: d.books,
    noteCount: d.note_count,
    lastActive: d.last_active,
    createdAt: d.created_at,
    updatedAt: d.updated_at
  };
}

export default function init(config: HttpClientConfig) {
  const client = getHttpClient(config);

  return {
    fetch: (uuid: string): Promise<RepetitionRuleData> => {
      const endpoint = `/repetition_rules/${uuid}`;

      return client.get<RepetitionRuleRespData>(endpoint).then(resp => {
        return mapData(resp);
      });
    },
    fetchAll: (): Promise<RepetitionRuleData[]> => {
      const endpoint = '/repetition_rules';

      return client.get<RepetitionRuleRespData[]>(endpoint).then(resp => {
        return resp.map(mapData);
      });
    },
    create: (params: CreateParams) => {
      const endpoint = '/repetition_rules';

      return client
        .post<RepetitionRuleRespData>(endpoint, params)
        .then(resp => {
          return mapData(resp);
        });
    },
    update: (uuid: string, params: UpdateParams) => {
      const endpoint = `/repetition_rules/${uuid}`;

      return client
        .patch<RepetitionRuleRespData>(endpoint, params)
        .then(resp => {
          return mapData(resp);
        });
    },
    remove: (uuid: string) => {
      const endpoint = `/repetition_rules/${uuid}`;

      return client.del(endpoint);
    }
  };
}
