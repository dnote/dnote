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

export interface EditorSession {
  sessionKey: string;
  noteUUID: string | null;
  bookUUID: string | null;
  bookLabel: string | null;
  content: string;
}

export interface EditorState {
  persisted: boolean;
  sessions: {
    [key: string]: EditorSession;
  };
}

export const MARK_PERSISTED = 'editor/MARK_PERSISTED';
export const CREATE_SESSION = 'editor/CREATE_SESSION';
export const FLUSH_CONTENT = 'editor/FLUSH_CONTENT';
export const UPDATE_BOOK = 'editor/UPDATE_BOOK';
export const RESET = 'editor/RESET';

export interface MarkPersistedAction {
  type: typeof MARK_PERSISTED;
}

export interface CreateSessionAction {
  type: typeof CREATE_SESSION;
  data: {
    noteUUID: string;
    bookUUID: string;
    bookLabel: string;
    content: string;
  };
}

export interface FlushContentAction {
  type: typeof FLUSH_CONTENT;
  data: {
    sessionKey: string;
    content: string;
  };
}

export interface UpdateBookAction {
  type: typeof UPDATE_BOOK;
  data: {
    sessionKey: string;
    uuid: string;
    label: string;
  };
}

export interface ResetAction {
  type: typeof RESET;
  data: {
    sessionKey: string;
  };
}

export type EditorActionType =
  | MarkPersistedAction
  | CreateSessionAction
  | FlushContentAction
  | UpdateBookAction
  | ResetAction;
