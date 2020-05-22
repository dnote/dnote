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

import { EmailPrefData } from 'jslib/operations/types';
import { RemoteData } from '../types';

export interface UserData {
  uuid: string;
  email: string;
  emailVerified: boolean;
  pro: boolean;
}

export type UserState = RemoteData<UserData>;
export type EmailPrefState = RemoteData<EmailPrefData>;

export interface AuthState {
  user: UserState;
  emailPreference: EmailPrefState;
}

export const START_FETCHING_USER = 'auth/START_FETCHING_USER';
export const RECEIVE_USER = 'auth/RECEIVE_USER';
export const RECEIVE_USER_ERROR = 'auth/RECEIVE_USER_ERROR';
export const RECEIVE_EMAIL_PREFERENCE = 'auth/RECEIVE_EMAIL_PREFERENCE';
export const RECEIVE_EMAIL_PREFERENCE_ERROR =
  'auth/RECEIVE_EMAIL_PREFERENCE_ERROR';

export interface StartFetchingUserAction {
  type: typeof START_FETCHING_USER;
}
export interface ReceiveUserAction {
  type: typeof RECEIVE_USER;
  data: { user: UserData };
}
export interface ReceiveUserErrorAction {
  type: typeof RECEIVE_USER_ERROR;
  data: {
    errorMessage: string;
  };
}
export interface ReceiveEmailPreferenceAction {
  type: typeof RECEIVE_EMAIL_PREFERENCE;
  data: {
    emailPreference: EmailPrefData;
  };
}
export interface ReceiveEmailPreferenceErrorAction {
  type: typeof RECEIVE_EMAIL_PREFERENCE_ERROR;
  data: {
    errorMessage: string;
  };
}

export type AuthActionType =
  | StartFetchingUserAction
  | ReceiveUserAction
  | ReceiveUserErrorAction
  | ReceiveEmailPreferenceAction
  | ReceiveEmailPreferenceErrorAction;
