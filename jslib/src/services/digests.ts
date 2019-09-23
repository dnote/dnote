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

export default function init(config: HttpClientConfig) {
  const client = getHttpClient(config);

  return {
    fetch: (digestUUID, { demo }) => {
      let endpoint;
      if (demo) {
        endpoint = `/demo/digests/${digestUUID}`;
      } else {
        endpoint = `/digests/${digestUUID}`;
      }

      return client.get(endpoint);
    },

    fetchAll: ({ page, demo }) => {
      let path;
      if (demo) {
        path = `/demo/digests`;
      } else {
        path = '/digests';
      }

      const endpoint = getPath(path, { page });

      return client.get(endpoint);
    }
  };
}
