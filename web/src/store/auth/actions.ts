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

import services from 'web/libs/services';
import { UserData, EmailPrefData, SourceData, SubscriptionData } from './type';
import { ThunkAction } from '../types';

import {
  RECEIVE_EMAIL_PREFERENCE,
  RECEIVE_EMAIL_PREFERENCE_ERROR,
  START_FETCHING_USER,
  RECEIVE_USER,
  RECEIVE_USER_ERROR,
  RECEIVE_SUBSCRIPTION,
  RECEIVE_SUBSCRIPTION_ERROR,
  START_FETCHING_SUBSCRIPTION,
  CLEAR_SUBSCRIPTION,
  START_FETCHING_SOURCE,
  RECEIVE_SOURCE,
  CLEAR_SOURCE,
  RECEIVE_SOURCE_ERROR,
  StartFetchingUserAction,
  ReceiveUserAction,
  ReceiveUserErrorAction,
  ReceiveEmailPreferenceAction,
  ReceiveEmailPreferenceErrorAction,
  StartFetchingSubscriptionAction,
  ReceiveSubscriptionAction,
  ClearSubscriptionAction,
  ReceiveSubscriptionErrorAction,
  StartFetchingSourceAction,
  ReceiveSourceAction,
  ClearSourceAction,
  ReceiveSourceErrorAction
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
              pro: false,
              classic: false
            })
          );
          return;
        }

        dispatch(receiveUserError(err.message));
      });
  };
}

export function startFetchingSubscription(): StartFetchingSubscriptionAction {
  return {
    type: START_FETCHING_SUBSCRIPTION
  };
}

export function receiveSubscription(subscription): ReceiveSubscriptionAction {
  return {
    type: RECEIVE_SUBSCRIPTION,
    data: { subscription }
  };
}

export function clearSubscription(): ClearSubscriptionAction {
  return {
    type: CLEAR_SUBSCRIPTION
  };
}

export function receiveSubscriptionError(
  errorMessage
): ReceiveSubscriptionErrorAction {
  return {
    type: RECEIVE_SUBSCRIPTION_ERROR,
    data: { errorMessage }
  };
}

export function getSubscription(): ThunkAction<SubscriptionData> {
  return dispatch => {
    dispatch(startFetchingSubscription());

    return services.payment
      .getSubscription()
      .then(subscription => {
        dispatch(receiveSubscription(subscription));
      })
      .catch(err => {
        console.log('error fetching subscription', err.message);
        dispatch(receiveSubscriptionError(err.message));
      });
  };
}

export function startFetchingSource(): StartFetchingSourceAction {
  return {
    type: START_FETCHING_SOURCE
  };
}

export function receiveSource(source): ReceiveSourceAction {
  return {
    type: RECEIVE_SOURCE,
    data: { source }
  };
}

export function clearSource(): ClearSourceAction {
  return {
    type: CLEAR_SOURCE
  };
}

export function receiveSourceError(
  errorMessage: string
): ReceiveSourceErrorAction {
  return {
    type: RECEIVE_SOURCE_ERROR,
    data: { errorMessage }
  };
}

export function getSource(): ThunkAction<SourceData> {
  return dispatch => {
    dispatch(startFetchingSource());

    return services.payment
      .getSource()
      .then(source => {
        console.log('source', source);
        dispatch(receiveSource(source));
      })
      .catch(err => {
        console.log('error fetching source', err.message);
        dispatch(receiveSourceError(err.message));
      });
  };
}
