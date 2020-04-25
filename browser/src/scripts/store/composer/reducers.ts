/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import {
  UPDATE_CONTENT,
  UPDATE_BOOK,
  RESET,
  RESET_BOOK,
  ComposerActionType,
  ComposerState
} from './types';

const initialState: ComposerState = {
  content: '',
  bookUUID: '',
  bookLabel: ''
};

export default function (
  state = initialState,
  action: ComposerActionType
): ComposerState {
  switch (action.type) {
    case UPDATE_CONTENT: {
      return {
        ...state,
        content: action.data.content
      };
    }
    case UPDATE_BOOK: {
      return {
        ...state,
        bookUUID: action.data.uuid,
        bookLabel: action.data.label
      };
    }
    case RESET_BOOK: {
      return {
        ...state,
        bookUUID: '',
        bookLabel: ''
      };
    }
    case RESET: {
      return initialState;
    }
    default:
      return state;
  }
}
