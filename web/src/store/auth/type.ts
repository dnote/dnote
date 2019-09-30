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

import { RemoteData } from '../types';

export interface UserData {
  uuid: string;
  email: string;
  emailVerified: boolean;
  pro: boolean;
  // TODO: remove once all classic users have been migrated
  classic: boolean;
}

export interface EmailPrefData {
  digestWeekly: boolean;
}

// TODO: type
export type SubscriptionData = any;
export type SourceData = any;

export type UserState = RemoteData<UserData>;
export type EmailPrefState = RemoteData<EmailPrefData>;
export type SubscriptionState = RemoteData<SubscriptionData>;
export type SourceState = RemoteData<SourceData>;

export interface AuthState {
  user: UserState;
  emailPreference: EmailPrefState;
  subscription: SubscriptionState;
  source: SourceState;
}

export const START_FETCHING_USER = 'auth/START_FETCHING_USER';
export const RECEIVE_USER = 'auth/RECEIVE_USER';
export const RECEIVE_USER_ERROR = 'auth/RECEIVE_USER_ERROR';
export const RECEIVE_EMAIL_PREFERENCE = 'auth/RECEIVE_EMAIL_PREFERENCE';
export const RECEIVE_EMAIL_PREFERENCE_ERROR =
  'auth/RECEIVE_EMAIL_PREFERENCE_ERROR';
export const START_FETCHING_SUBSCRIPTION = 'auth/START_FETCHING_SUBSCRIPTION';
export const RECEIVE_SUBSCRIPTION = 'auth/RECEIVE_SUBSCRIPTION';
export const CLEAR_SUBSCRIPTION = 'auth/CLEAR_SUBSCRIPTION';
export const RECEIVE_SUBSCRIPTION_ERROR = 'auth/RECEIVE_SUBSCRIPTION_ERROR';
export const START_FETCHING_SOURCE = 'auth/START_FETCHING_SOURCE';
export const RECEIVE_SOURCE = 'auth/RECEIVE_SOURCE';
export const CLEAR_SOURCE = 'auth/CLEAR_SOURCE';
export const RECEIVE_SOURCE_ERROR = 'auth/RECEIVE_SOURCE_ERROR';

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
export interface StartFetchingSubscriptionAction {
  type: typeof START_FETCHING_SUBSCRIPTION;
}
export interface ReceiveSubscriptionAction {
  type: typeof RECEIVE_SUBSCRIPTION;
  data: {
    subscription: SubscriptionData;
  };
}
export interface ClearSubscriptionAction {
  type: typeof CLEAR_SUBSCRIPTION;
}
export interface ReceiveSubscriptionErrorAction {
  type: typeof RECEIVE_SUBSCRIPTION_ERROR;
  data: {
    errorMessage: string;
  };
}
export interface StartFetchingSourceAction {
  type: typeof START_FETCHING_SOURCE;
}
export interface ReceiveSourceAction {
  type: typeof RECEIVE_SOURCE;
  data: {
    source: SourceData;
  };
}
export interface ClearSourceAction {
  type: typeof CLEAR_SOURCE;
}
export interface ReceiveSourceErrorAction {
  type: typeof RECEIVE_SOURCE_ERROR;
  data: {
    errorMessage: string;
  };
}

export type AuthActionType =
  | StartFetchingUserAction
  | ReceiveUserAction
  | ReceiveUserErrorAction
  | ReceiveEmailPreferenceAction
  | ReceiveEmailPreferenceErrorAction
  | StartFetchingSubscriptionAction
  | ReceiveSubscriptionAction
  | ClearSubscriptionAction
  | ReceiveSubscriptionErrorAction
  | StartFetchingSourceAction
  | ReceiveSourceAction
  | ClearSourceAction
  | ReceiveSourceErrorAction;
