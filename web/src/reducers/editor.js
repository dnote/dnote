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

import {
  UPDATE_CONTENT,
  UPDATE_BOOK_UUID,
  RESET,
  STAGE_NOTE,
  COMMIT_NOTE,
  MARK_DIRTY
} from '../actions/editor';

const initialState = {
  noteUUID: null,
  bookUUID: null,
  content: '',
  dirty: false
};

export default function(state = initialState, action) {
  switch (action.type) {
    case STAGE_NOTE: {
      return {
        ...state,
        noteUUID: action.data.noteUUID,
        bookUUID: action.data.bookUUID,
        content: action.data.content,
        dirty: false
      };
    }
    case COMMIT_NOTE: {
      return {
        ...state,
        dirty: false
      };
    }
    case UPDATE_CONTENT: {
      return {
        ...state,
        content: action.data.content
      };
    }
    case UPDATE_BOOK_UUID: {
      return {
        ...state,
        bookUUID: action.data.uuid
      };
    }
    case MARK_DIRTY: {
      return {
        ...state,
        dirty: true
      };
    }
    case RESET: {
      return initialState;
    }
    default:
      return state;
  }
}
