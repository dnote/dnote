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

import { LOGIN, LOGOUT, LogoutAction } from './types';
import { ThunkAction } from '../types';
import initServices from '../../utils/services';

export function login({ email, password }): ThunkAction<void> {
  return (dispatch, getState) => {
    const { settings } = getState();
    const { apiUrl } = settings;

    return initServices(apiUrl)
      .users.signin({ email, password })
      .then(resp => {
        dispatch({
          type: LOGIN,
          data: {
            sessionKey: resp.key,
            sessionKeyExpiry: resp.expiresAt
          }
        });
      });
  };
}

export function logout(): LogoutAction {
  return {
    type: LOGOUT
  };
}
