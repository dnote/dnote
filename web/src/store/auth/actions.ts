/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

import services from 'web/libs/services';
import { EmailPrefData } from 'jslib/operations/types';
import { UserData } from './type';
import { ThunkAction } from '../types';

import {
  RECEIVE_EMAIL_PREFERENCE,
  RECEIVE_EMAIL_PREFERENCE_ERROR,
  START_FETCHING_USER,
  RECEIVE_USER,
  RECEIVE_USER_ERROR,
  StartFetchingUserAction,
  ReceiveUserAction,
  ReceiveUserErrorAction,
  ReceiveEmailPreferenceAction,
  ReceiveEmailPreferenceErrorAction
} from './type';

function startFetchingUser(): StartFetchingUserAction {
  return {
    type: START_FETCHING_USER
  };
}

export function receiveUser(user: UserData): ReceiveUserAction {
  return {
    type: RECEIVE_USER,
    data: { user }
  };
}

function receiveUserError(errorMessage): ReceiveUserErrorAction {
  return {
    type: RECEIVE_USER_ERROR,
    data: { errorMessage }
  };
}

export function receiveEmailPreference(
  emailPreference: EmailPrefData
): ReceiveEmailPreferenceAction {
  return {
    type: RECEIVE_EMAIL_PREFERENCE,
    data: { emailPreference }
  };
}

export function getEmailPreferenceError(
  errorMessage: string
): ReceiveEmailPreferenceErrorAction {
  return {
    type: RECEIVE_EMAIL_PREFERENCE_ERROR,
    data: { errorMessage }
  };
}

export function getEmailPreference(token?: string) {
  return dispatch => {
    return services.users
      .getEmailPreference({ token })
      .then(emailPreference => {
        dispatch(receiveEmailPreference(emailPreference));
      })
      .catch(err => {
        console.log('error fetching email preference', err.message);
        dispatch(getEmailPreferenceError(err.message));
      });
  };
}

interface GetCurrentUserOptions {
  refresh?: boolean;
}

export function getCurrentUser(
  options: GetCurrentUserOptions = {}
): ThunkAction<UserData | void> {
  return dispatch => {
    if (!options.refresh) {
      dispatch(startFetchingUser());
    }

    return services.users
      .getMe()
      .then(user => {
        dispatch(receiveUser(user));

        return user;
      })
      .catch(err => {
        // 401 if not logged in
        if (err.response.status === 401) {
          dispatch(
            receiveUser({
              uuid: '',
              email: '',
              emailVerified: false,
              pro: false
            })
          );
          return;
        }

        dispatch(receiveUserError(err.message));
      });
  };
}
