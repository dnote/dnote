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

import { BookDomain } from 'jslib/operations/types';
import {
  RECEIVE,
  START_FETCHING,
  ERROR,
  RESET,
  SET_NOTE_REVIEWED,
  DigestState,
  DigestActionType
} from './type';

const initialState: DigestState = {
  data: {
    uuid: '',
    createdAt: '',
    updatedAt: '',
    version: 0,
    notes: [],
    isRead: false,
    repetitionRule: {
      uuid: '',
      title: '',
      enabled: false,
      hour: 0,
      minute: 0,
      bookDomain: BookDomain.All,
      frequency: 0,
      books: [],
      lastActive: 0,
      nextActive: 0,
      noteCount: 0,
      createdAt: '',
      updatedAt: ''
    }
  },
  isFetching: false,
  isFetched: false,
  errorMessage: null
};

export default function(
  state = initialState,
  action: DigestActionType
): DigestState {
  switch (action.type) {
    case START_FETCHING: {
      return {
        ...state,
        isFetching: true,
        isFetched: false
      };
    }
    case ERROR: {
      return {
        ...state,
        isFetching: false,
        errorMessage: action.data.errorMessage
      };
    }
    case RECEIVE: {
      return {
        ...state,
        data: action.data.digest,
        isFetching: false,
        isFetched: true,
        errorMessage: null
      };
    }
    case SET_NOTE_REVIEWED: {
      return {
        ...state,
        data: {
          ...state.data,
          notes: state.data.notes.map(note => {
            if (action.data.noteUUID === note.uuid) {
              const isReviewed = action.data.isReviewed;

              return {
                ...note,
                isReviewed
              };
            }

            return note;
          })
        }
      };
    }
    case RESET: {
      return initialState;
    }
    default:
      return state;
  }
}
