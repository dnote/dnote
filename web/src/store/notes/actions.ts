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

import operations from 'web/libs/operations';
import { Filters } from 'jslib/helpers/filters';
import {
  ADD,
  REFRESH,
  RECEIVE,
  START_FETCHING,
  RECEIVE_ERROR,
  RESET,
  REMOVE,
  AddAction,
  RemoveAction,
  RefreshAction,
  ReceiveAction,
  StartFetchingAction,
  ReceiveErrorAction,
  ResetAction
} from './type';
import { ThunkAction } from '../types';

export function addNote(note): AddAction {
  return {
    type: ADD,
    data: {
      note
    }
  };
}

export function removeNote({ noteUUID }): RemoveAction {
  return {
    type: REMOVE,
    data: {
      noteUUID
    }
  };
}

export function refreshNoteInList({ noteUUID, refreshedNote }): RefreshAction {
  return {
    type: REFRESH,
    data: {
      noteUUID,
      book: refreshedNote.book,
      content: refreshedNote.content,
      isPublic: refreshedNote.public
    }
  };
}

function receiveNotes(notes, total): ReceiveAction {
  return {
    type: RECEIVE,
    data: {
      notes,
      total
    }
  };
}

function startFetchingNotes(): StartFetchingAction {
  return {
    type: START_FETCHING
  };
}

function receiveError(error: string): ReceiveErrorAction {
  return {
    type: RECEIVE_ERROR,
    data: {
      error
    }
  };
}

export function resetNotes(): ResetAction {
  return {
    type: RESET
  };
}

export function getNotes(filters: Filters): ThunkAction<void> {
  return async dispatch => {
    dispatch(startFetchingNotes());

    return operations.notes
      .fetch(filters)
      .then(res => {
        const { notes, total } = res;
        dispatch(receiveNotes(notes, total));
      })
      .catch(err => {
        console.log('Error fetching notes', err.stack);
        dispatch(receiveError(err.message));
      });
  };
}
