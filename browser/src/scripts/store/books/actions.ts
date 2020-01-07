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

import initServices from '../../utils/services';

import {
  START_FETCHING,
  RECEIVE,
  RECEIVE_ERROR,
  StartFetchingAction,
  ReceiveAction,
  ReceiveErrorAction
} from './types';

function startFetchingBooks(): StartFetchingAction {
  return {
    type: START_FETCHING
  };
}

function receiveBooks(books): ReceiveAction {
  return {
    type: RECEIVE,
    data: {
      books
    }
  };
}

function receiveBooksError(error: string): ReceiveErrorAction {
  return {
    type: RECEIVE_ERROR,
    data: {
      error
    }
  };
}

export function fetchBooks() {
  return (dispatch, getState) => {
    dispatch(startFetchingBooks());

    const { settings, auth } = getState();
    const services = initServices(settings.apiUrl);

    services.books
      .fetch(
        {},
        {
          headers: {
            Authorization: `Bearer ${auth.sessionKey}`
          }
        }
      )
      .then(books => {
        dispatch(receiveBooks(books));
      })
      .catch(err => {
        console.log('error fetching books', err);
        dispatch(receiveBooksError(err));
      });
  };
}
