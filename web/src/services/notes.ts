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
import { NoteData } from '../operations/types';
import { Filters } from '../libs/filters';

interface CreateParams {
  book_uuid: string;
  content: string;
}

interface CreateResponse {
  result: NoteData;
}

export function create(params: CreateParams): Promise<CreateResponse> {
  return apiClient.post<CreateResponse>('/v2/notes', params);
}

interface UpdateParams {
  book_uuid?: string;
  content?: string;
  public?: boolean;
}

interface UpdateNoteResp {
  status: number;
  result: NoteData;
}

export function update(noteUUID: string, params: UpdateParams) {
  const endpoint = `/v1/notes/${noteUUID}`;

  return apiClient.patch<UpdateNoteResp>(endpoint, params);
}

export function remove(noteUUID: string) {
  const endpoint = `/v1/notes/${noteUUID}`;

  return apiClient.del(endpoint, {});
}

interface FetchResponse {
  notes: NoteData[];
  total: number;
}

export function fetch(filters: Filters) {
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

  return apiClient.get<FetchResponse>(endpoint, {});
}

type FetchOneResponse = NoteData;

interface FetchOneQuery {
  q?: string;
}

export function fetchOne(
  noteUUID: string,
  params: FetchOneQuery
): Promise<FetchOneResponse> {
  const endpoint = getPath(`/notes/${noteUUID}`, params);

  return apiClient.get<FetchOneResponse>(endpoint, {});
}

export function legacyFetchNotes() {
  const endpoint = '/legacy/notes';

  return apiClient.get(endpoint, { credentials: 'include' });
}
