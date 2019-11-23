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

import { removeKey } from 'jslib/helpers/obj';
import { getEditorSessionkey } from 'web/libs/editor';
import {
  EditorState,
  EditorSession,
  EditorActionType,
  FLUSH_CONTENT,
  UPDATE_BOOK,
  RESET,
  CREATE_SESSION,
  MARK_PERSISTED
} from './type';

function makeSession(key: string): EditorSession {
  return {
    sessionKey: key,
    noteUUID: null,
    bookUUID: null,
    bookLabel: null,
    content: ''
  };
}

const initialState: EditorState = {
  persisted: false,
  sessions: {}
};

export default function(
  state = initialState,
  action: EditorActionType
): EditorState {
  switch (action.type) {
    case CREATE_SESSION: {
      const { data } = action;

      const sessionKey = getEditorSessionkey(data.noteUUID);

      return {
        ...state,
        persisted: false,
        sessions: {
          ...state.sessions,
          [sessionKey]: {
            sessionKey,
            noteUUID: data.noteUUID,
            bookUUID: data.bookUUID,
            bookLabel: data.bookLabel,
            content: data.content
          }
        }
      };
    }
    case FLUSH_CONTENT: {
      const { data } = action;

      return {
        ...state,
        persisted: false,
        sessions: {
          ...state.sessions,
          [data.sessionKey]: {
            ...state.sessions[data.sessionKey],
            content: data.content
          }
        }
      };
    }
    case UPDATE_BOOK: {
      const { data } = action;

      return {
        ...state,
        persisted: false,
        sessions: {
          ...state.sessions,
          [data.sessionKey]: {
            ...state.sessions[data.sessionKey],
            bookUUID: action.data.uuid,
            bookLabel: action.data.label
          }
        }
      };
    }
    case MARK_PERSISTED: {
      return {
        ...state,
        persisted: true
      };
    }
    case RESET: {
      const { data } = action;

      return {
        ...state,
        persisted: false,
        sessions: removeKey(state.sessions, data.sessionKey)
      };
    }
    default:
      return state;
  }
}
