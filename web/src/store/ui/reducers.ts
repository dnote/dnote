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
import { SET_MESSAGE, UNSET_MESSAGE, UIState, UIActionType } from './type';

export const initialState: UIState = {
  message: {}
};

export default function(state = initialState, action: UIActionType): UIState {
  switch (action.type) {
    case SET_MESSAGE: {
      return {
        ...state,
        message: {
          ...state.message,
          [action.data.path]: {
            content: action.data.message,
            kind: action.data.kind
          }
        }
      };
    }
    case UNSET_MESSAGE: {
      return {
        ...state,
        message: removeKey(state.message, action.data.path)
      };
    }
    default:
      return state;
  }
}
