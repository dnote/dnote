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

import { fetchCalendar } from '../services/users';

export const RECEIVE_CALENDAR = 'calendar/RECEIVE_CALENDAR';
export const RECEIVE_CALENDAR_ERROR = 'calendar/RECEIVE_CALENDAR_ERROR';
export const START_FETCHING = 'calendar/START_FETCHING';

function receiveCalendar(items) {
  return {
    type: RECEIVE_CALENDAR,
    data: { items }
  };
}

function startFetchingCalendar() {
  return {
    type: START_FETCHING
  };
}

function receiveCalendarError(error) {
  return {
    type: RECEIVE_CALENDAR_ERROR,
    data: {
      error: {
        status: error.status,
        message: error.message
      }
    }
  };
}

export function getLessonCalendar(options = { demo: false }) {
  return dispatch => {
    dispatch(startFetchingCalendar());

    const { demo } = options;

    return fetchCalendar({ demo })
      .then(res => {
        dispatch(receiveCalendar(res));
      })
      .catch(err => {
        dispatch(receiveCalendarError(err));
      });
  };
}

// helpers

export function isCalendarLoaded(state) {
  return state.calendar.isFetched && !state.calendar.error;
}
