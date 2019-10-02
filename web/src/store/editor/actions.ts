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

// import { receiveNote } from '../note';
// import { refreshNoteInList } from '../notes';

import {
  FLUSH_CONTENT,
  UPDATE_BOOK,
  RESET,
  CREATE_SESSION,
  MARK_PERSISTED,
  MarkPersistedAction,
  CreateSessionAction,
  FlushContentAction,
  UpdateBookAction,
  ResetAction
} from './type';

export function createSession({
  noteUUID,
  bookUUID,
  bookLabel,
  content
}): CreateSessionAction {
  return {
    type: CREATE_SESSION,
    data: { noteUUID, bookUUID, bookLabel, content }
  };
}

export function flushContent(
  sessionKey: string,
  content: string
): FlushContentAction {
  return {
    type: FLUSH_CONTENT,
    data: { sessionKey, content }
  };
}

export interface UpdateBookActionParam {
  sessionKey: string;
  uuid: string;
  label: string;
}

export function updateBook({
  sessionKey,
  uuid,
  label
}: UpdateBookActionParam): UpdateBookAction {
  return {
    type: UPDATE_BOOK,
    data: {
      sessionKey,
      uuid,
      label
    }
  };
}

export function resetEditor(sessionKey: string): ResetAction {
  return {
    type: RESET,
    data: {
      sessionKey
    }
  };
}

export function markPersisted(): MarkPersistedAction {
  return {
    type: MARK_PERSISTED
  };
}
