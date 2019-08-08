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

import {
  ADD,
  REFRESH,
  RECEIVE,
  RECEIVE_MORE,
  START_FETCHING,
  START_FETCHING_MORE,
  RECEIVE_ERROR,
  RESET,
  REMOVE,
  AddAction,
  RemoveAction,
  RefreshAction,
  ReceiveAction,
  ReceiveMoreAction,
  StartFetchingAction,
  StartFetchingMoreAction,
  ReceiveErrorAction,
  ResetAction
} from './type';
import { ThunkAction } from '../types';
import * as notesOperation from '../../operations/notes';
import { parsePrevDate } from '../../libs/notes';
import { Facets } from '../../libs/facets';
// import { NoteData } from '../../operations/types';

export function addNote(note, year, month): AddAction {
  return {
    type: ADD,
    data: {
      note,
      year,
      month
    }
  };
}

export function removeNote({ year, month, noteUUID }): RemoveAction {
  return {
    type: REMOVE,
    data: {
      year,
      month,
      noteUUID
    }
  };
}

export function refreshNoteInList({ noteUUID, refreshedNote }): RefreshAction {
  const date = new Date();
  const year = date.getUTCFullYear();
  const month = date.getUTCMonth() + 1;

  return {
    type: REFRESH,
    data: {
      year,
      month,
      noteUUID,
      book: refreshedNote.book,
      content: refreshedNote.content,
      isPublic: refreshedNote.public
    }
  };
}

function receiveNotes(notes, total, year, month, prevDate): ReceiveAction {
  return {
    type: RECEIVE,
    data: {
      notes,
      total,
      year,
      month,
      prevDate
    }
  };
}

function receiveMoreNotes(
  notes,
  total,
  year,
  month,
  prevDate
): ReceiveMoreAction {
  return {
    type: RECEIVE_MORE,
    data: {
      notes,
      year,
      month,
      prevDate
    }
  };
}

function startFetchingNotes(year, month): StartFetchingAction {
  return {
    type: START_FETCHING,
    data: {
      year,
      month
    }
  };
}

function startFetchingMoreNotes(year, month): StartFetchingMoreAction {
  return {
    type: START_FETCHING_MORE,
    data: {
      year,
      month
    }
  };
}

function receiveError(year, month, error): ReceiveErrorAction {
  return {
    type: RECEIVE_ERROR,
    data: {
      year,
      month,
      error
    }
  };
}

export function resetNotes(): ResetAction {
  return {
    type: RESET
  };
}

export function getNotes(year, month, facets: Facets): ThunkAction<void> {
  return async dispatch => {
    dispatch(startFetchingNotes(year, month));

    const query = {
      year,
      month,
      ...facets
    };

    return notesOperation
      .fetch(query)
      .then(res => {
        const { notes, total, prev_date: prevDate } = res;
        dispatch(receiveNotes(notes, total, year, month, prevDate));
      })
      .catch(err => {
        console.log('Error fetching notes', err.stack);
        dispatch(receiveError(year, month, err.message));
      });
  };
}

export function getMoreNotes(year, month, page, queryObj = {}) {
  return (dispatch, getState) => {
    const state = getState();
    if (state.notes.isFetchingMore || !state.notes.prevDate) {
      return null;
    }

    dispatch(startFetchingMoreNotes(year, month));

    const q = {
      ...queryObj,
      year,
      month,
      page
    };

    return notesOperation
      .fetch(q)
      .then(res => {
        const { notes, total, prev_date: prevDate } = res;
        dispatch(receiveMoreNotes(notes, total, year, month, prevDate));
      })
      .catch(err => {
        console.log('err', err);
        // dispatch(receiveError(year, month, err));
      });
  };
}

interface GetInitialNotesParams {
  facets: Facets;
  year: number;
  month: number;
}

export function getInitialNotes({
  facets,
  year,
  month
}: GetInitialNotesParams) {
  return async (dispatch, getState) => {
    const { notes } = getState();
    const hasError = notes.groups.some(group => {
      return Boolean(group.error);
    });

    if (!notes.prevDate || hasError) {
      return;
    }

    await dispatch(getNotes(year, month, facets));

    const { notes: updatedNotes } = getState();

    const notesTotal = updatedNotes.groups.reduce((acc, group) => {
      return acc + group.total;
    }, 0);

    if (notesTotal < 8) {
      const { year: prevYear, month: prevMonth } = parsePrevDate(
        updatedNotes.prevDate
      );

      await dispatch(
        getInitialNotes({
          facets,
          year: prevYear,
          month: prevMonth
        })
      );
    }
  };
}
