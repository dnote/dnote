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

export interface CreateDeleteNoteReviewPayload {
  digestUUID: string;
  noteUUID: string;
}

export default function init(config: HttpClientConfig) {
  const client = getHttpClient(config);

  return {
    create: ({
      digestUUID,
      noteUUID
    }: CreateDeleteNoteReviewPayload): Promise<void> => {
      const endpoint = '/note_review';
      const payload = {
        digest_uuid: digestUUID,
        note_uuid: noteUUID
      };

      return client.post(endpoint, payload);
    },

    remove: ({
      digestUUID,
      noteUUID
    }: CreateDeleteNoteReviewPayload): Promise<void> => {
      const endpoint = '/note_review';
      const payload = {
        digest_uuid: digestUUID,
        note_uuid: noteUUID
      };

      return client.del(endpoint, payload);
    }
  };
}
