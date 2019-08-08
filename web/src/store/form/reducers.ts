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

import { FormState, FormActionType, UPDATE_AUTH_EMAIL } from './type';

export const initialState: FormState = {
  auth: {
    email: ''
  }
};

export default function(
  state = initialState,
  action: FormActionType
): FormState {
  switch (action.type) {
    case UPDATE_AUTH_EMAIL: {
      const { data } = action;

      return {
        ...state,
        auth: {
          ...state.auth,
          email: data.email
        }
      };
    }
    default:
      return state;
  }
}
