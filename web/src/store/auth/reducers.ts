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
import {
  AuthState,
  AuthActionType,
  UserData,
  RECEIVE_EMAIL_PREFERENCE,
  RECEIVE_EMAIL_PREFERENCE_ERROR,
  START_FETCHING_USER,
  RECEIVE_USER,
  RECEIVE_USER_ERROR
} from './type';

export const initialState: AuthState = {
  user: {
    isFetching: false,
    isFetched: false,
    data: {
      uuid: '',
      email: '',
      emailVerified: false,
      pro: false
    },
    errorMessage: ''
  },
  emailPreference: {
    isFetching: false,
    isFetched: false,
    data: {
      inactiveReminder: false,
      productUpdate: false
    },
    errorMessage: ''
  }
};

function reduceUsers(
  state = initialState.user,
  action: AuthActionType
): RemoteData<UserData> {
  switch (action.type) {
    case START_FETCHING_USER: {
      return {
        ...state,
        isFetching: true,
        isFetched: false
      };
    }
    case RECEIVE_USER: {
      const { user } = action.data;
      return {
        ...state,
        data: {
          uuid: user.uuid,
          email: user.email,
          emailVerified: user.emailVerified,
          pro: user.pro
        },
        errorMessage: '',
        isFetching: false,
        isFetched: true
      };
    }
    case RECEIVE_USER_ERROR: {
      return {
        ...state,
        isFetching: false,
        isFetched: false,
        errorMessage: action.data.errorMessage
      };
    }
    default:
      return state;
  }
}

function reducerEmailPreference(
  state = initialState.emailPreference,
  action: AuthActionType
): RemoteData<EmailPrefData> {
  switch (action.type) {
    case RECEIVE_EMAIL_PREFERENCE:
      return {
        ...state,
        errorMessage: '',
        isFetching: false,
        isFetched: true,
        data: action.data.emailPreference
      };
    case RECEIVE_EMAIL_PREFERENCE_ERROR: {
      return {
        ...state,
        isFetching: false,
        isFetched: false,
        errorMessage: action.data.errorMessage
      };
    }
    default:
      return state;
  }
}

export default function (
  state = initialState,
  action: AuthActionType
): AuthState {
  switch (action.type) {
    case START_FETCHING_USER:
    case RECEIVE_USER_ERROR:
    case RECEIVE_USER: {
      return {
        ...state,
        user: reduceUsers(state.user, action)
      };
    }
    case RECEIVE_EMAIL_PREFERENCE:
    case RECEIVE_EMAIL_PREFERENCE_ERROR: {
      return {
        ...state,
        emailPreference: reducerEmailPreference(state.emailPreference, action)
      };
    }
    default:
      return state;
  }
}
