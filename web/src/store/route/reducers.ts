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

import { SET_PREV_LOCATION, RouteActionType, RouteState } from './type';

export const initialState: RouteState = {
  prevLocation: {
    pathname: '',
    hash: '',
    search: '',
    state: {}
  }
};

export default function(
  state = initialState,
  action: RouteActionType
): RouteState {
  switch (action.type) {
    case SET_PREV_LOCATION: {
      return {
        ...state,
        prevLocation: {
          ...state.prevLocation,
          pathname: action.data.pathname,
          search: action.data.search,
          hash: action.data.hash,
          state: action.data.state
        }
      };
    }
    default:
      return state;
  }
}
