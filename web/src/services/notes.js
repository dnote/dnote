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

import { getPath } from '../libs/url';
import { apiClient } from '../libs/http';

export function create({ bookUUID, content }) {
  const payload = {
    book_uuid: bookUUID,
    content
  };

  return apiClient.post('/v2/notes', payload);
}

export function update(noteUUID, { bookUUID, content, isPublic }) {
  const endpoint = `/v1/notes/${noteUUID}`;

  const payload = {};

  if (bookUUID) {
    payload.book_uuid = bookUUID;
  }
  if (content) {
    payload.content = content;
  }
  if (isPublic !== undefined) {
    payload.public = isPublic;
  }

  return apiClient.patch(endpoint, payload);
}

export function remove(noteUUID) {
  const endpoint = `/v1/notes/${noteUUID}`;

  return apiClient.del(endpoint);
}

export function fetch(queryObj, demo) {
  const { year, month, q, book, page } = queryObj;

  let endpoint = getPath('/notes', {
    year,
    month,
    q,
    book,
    page
  });

  if (demo) {
    endpoint = `/demo${endpoint}`;
  }

  return apiClient.get(endpoint);
}

export function fetchOne(noteUUID, options = {}) {
  const { demo } = options;

  let endpoint = `/notes/${noteUUID}`;

  if (demo) {
    endpoint = `/demo${endpoint}`;
  }

  return apiClient.get(endpoint);
}

export function legacyFetchNotes() {
  const endpoint = '/legacy/notes';

  return apiClient.get(endpoint, { credentials: 'include' });
}
