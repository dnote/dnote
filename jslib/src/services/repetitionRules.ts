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

export interface UpdateParams {
  title?: string;
  hour?: number;
  minute?: number;
  frequency?: number;
  book_domain?: BookDomain;
  note_count?: number;
  book_uuids?: string[];
  enabled?: boolean;
}

export default function init(config: HttpClientConfig) {
  const client = getHttpClient(config);

  return {
    fetchAll: () => {
      const endpoint = '/repetition_rules';

      return client.get<RepetitionRuleData[]>(endpoint);
    },
    create: (params: CreateParams) => {
      const endpoint = '/repetition_rules';

      console.log(params);

      return client.post<RepetitionRuleData>(endpoint, params);
    },
    update: (uuid: string, params: UpdateParams) => {
      const endpoint = `/repetition_rules/${uuid}`;

      return client.patch<RepetitionRuleData>(endpoint, params);
    }
  };
}
