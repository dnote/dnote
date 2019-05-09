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

import { getMe } from '../services/users';
import { apiClient } from '../libs/http';
import * as paymentService from '../services/payment';
import * as usersService from '../services/users';

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

function startFetchingUser() {
  return {
    type: START_FETCHING_USER
  };
}

export function receiveUser(user) {
  return {
    type: RECEIVE_USER,
    data: { user }
  };
}

function receiveUserError(errorMessage) {
  return {
    type: RECEIVE_USER_ERROR,
    data: { errorMessage }
  };
}

export function receiveEmailPreference(emailPreference) {
  return {
    type: RECEIVE_EMAIL_PREFERENCE,
    data: { emailPreference }
  };
}

export function getEmailPreferenceError(errorMessage) {
  return {
    type: RECEIVE_EMAIL_PREFERENCE_ERROR,
    data: { errorMessage }
  };
}

export function getEmailPreference(token) {
  return dispatch => {
    return usersService
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

export function getCurrentUser(options = {}) {
  return dispatch => {
    if (!options.refresh) {
      dispatch(startFetchingUser());
    }

    return getMe()
      .then(user => {
        dispatch(receiveUser(user));

        return user;
      })
      .catch(err => {
        // 401 if not logged in
        if (err.response.status === 401) {
          dispatch(receiveUser({}));
          return;
        }

        dispatch(receiveUserError(err.message));
      });
  };
}

export function startFetchingSubscription() {
  return {
    type: START_FETCHING_SUBSCRIPTION
  };
}

export function receiveSubscription(subscription) {
  return {
    type: RECEIVE_SUBSCRIPTION,
    data: { subscription }
  };
}

export function clearSubscription() {
  return {
    type: CLEAR_SUBSCRIPTION
  };
}

export function receiveSubscriptionError(errorMessage) {
  return {
    type: RECEIVE_SUBSCRIPTION_ERROR,
    data: { errorMessage }
  };
}

export function getSubscription() {
  return dispatch => {
    dispatch(startFetchingSubscription());

    return paymentService
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

export function startFetchingSource() {
  return {
    type: START_FETCHING_SOURCE
  };
}

export function receiveSource(source) {
  return {
    type: RECEIVE_SOURCE,
    data: { source }
  };
}

export function clearSource() {
  return {
    type: CLEAR_SOURCE
  };
}

export function receiveSourceError(errorMessage) {
  return {
    type: RECEIVE_SOURCE_ERROR,
    data: { errorMessage }
  };
}

export function getSource() {
  return dispatch => {
    dispatch(startFetchingSource());

    return paymentService
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

export function legacyGetCurrentUser() {
  return dispatch => {
    return apiClient
      .get('/legacy/me')
      .then(res => {
        const { user } = res;

        dispatch(receiveUser(user));
      })
      .catch(err => {
        // 401 if not logged in
        if (err.status === 401) {
          return;
        }

        console.log('getUser error', err);
      });
  };
}
