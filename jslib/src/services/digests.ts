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

import { getHttpClient, HttpClientConfig } from '../helpers/http';
import { getPath } from '../helpers/url';
import { DigestData } from '../operations/types';
import { RepetitionRuleRespData } from './repetitionRules';

export interface FetchAllResponse {
  total: number;
  items: {
    uuid: string;
    version: number;
    repetition_rule: RepetitionRuleRespData;
    created_at: string;
    updated_at: string;
    receipts: {
      created_at: string;
      updated_at: string;
    }[];
  }[];
}

export interface FetchAllResult {
  total: number;
  items: DigestData[];
}

export default function init(config: HttpClientConfig) {
  const client = getHttpClient(config);

  return {
    fetch: (digestUUID: string) => {
      const endpoint = `/digests/${digestUUID}`;

      return client.get(endpoint);
    },

    fetchAll: ({ page }): Promise<FetchAllResult> => {
      const path = '/digests';

      const endpoint = getPath(path, { page });

      return client.get<FetchAllResponse>(endpoint).then(res => {
        return {
          total: res.total,
          items: res.items.map(item => {
            return {
              uuid: item.uuid,
              createdAt: item.created_at,
              updatedAt: item.updated_at,
              version: item.version,
              notes: [],
              repetitionRule: {
                uuid: item.repetition_rule.uuid,
                title: item.repetition_rule.title,
                enabled: item.repetition_rule.enabled,
                hour: item.repetition_rule.hour,
                minute: item.repetition_rule.minute,
                bookDomain: item.repetition_rule.book_domain,
                frequency: item.repetition_rule.frequency,
                books: item.repetition_rule.books,
                lastActive: item.repetition_rule.last_active,
                nextActive: item.repetition_rule.next_active,
                noteCount: item.repetition_rule.note_count,
                createdAt: item.repetition_rule.created_at,
                updatedAt: item.repetition_rule.updated_at
              },
              receipts: item.receipts.map(receipt => {
                return {
                  createdAt: receipt.created_at,
                  updatedAt: receipt.updated_at
                };
              })
            };
          })
        };
      });
    }
  };
}
