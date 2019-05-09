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

import * as notesService from '../services/notes';
import { decryptNote } from '../crypto/notes';
import { parsePrevDate } from '../libs/notes';

export const ADD = 'notes/ADD';
export const REFRESH = 'notes/REFRESH';
export const RECEIVE = 'notes/RECEIVE';
export const RECEIVE_MORE = 'notes/RECEIVE_MORE';
export const START_FETCHING = 'notes/START_FETCHING';
export const START_FETCHING_MORE = 'notes/START_FETCHING_MORE';
export const RECEIVE_ERROR = 'notes/RECEIVE_ERROR';
export const RESET = 'notes/RESET';
export const REMOVE = 'notes/REMOVE';

export function addNote(note, year, month) {
  return {
    type: ADD,
    data: {
      note: {
        data: note,
        errorMessage: null
      },
      year,
      month
    }
  };
}

export function removeNote({ year, month, noteUUID }) {
  return {
    type: REMOVE,
    data: {
      year,
      month,
      noteUUID
    }
  };
}

export function refreshNoteInList({ noteUUID, refreshedNote }) {
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

function receiveNotes(notes, total, year, month, prevDate) {
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

function receiveMoreNotes(notes, total, year, month, prevDate) {
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

function startFetchingNotes(year, month) {
  return {
    type: START_FETCHING,
    data: {
      year,
      month
    }
  };
}

function startFetchingMoreNotes(year, month) {
  return {
    type: START_FETCHING_MORE,
    data: {
      year,
      month
    }
  };
}

function receiveError(year, month, error) {
  return {
    type: RECEIVE_ERROR,
    data: {
      year,
      month,
      error
    }
  };
}

export function resetNotes() {
  return {
    type: RESET
  };
}

function decryptNotes(notes, cipherKeyBuf) {
  return notes.map(note => {
    return decryptNote(note, cipherKeyBuf)
      .then(noteDec => {
        return {
          data: noteDec,
          errorMessage: null
        };
      })
      .catch(err => {
        return {
          data: note,
          errorMessage: err.message
        };
      });
  });
}

export function getNotes(cipherKeyBuf, year, month, queryObj, demo) {
  return async dispatch => {
    dispatch(startFetchingNotes(year, month));

    const { q, book } = queryObj;

    const query = {
      year,
      month
    };
    if (q) {
      query.q = q;
    }
    if (book) {
      query.book = book;
    }

    return notesService
      .fetch(query, demo)
      .then(res => {
        const { notes, total, prev_date: prevDate } = res;

        const p = decryptNotes(notes, cipherKeyBuf);

        return Promise.all(p).then(notesDec => {
          dispatch(receiveNotes(notesDec, total, year, month, prevDate));
        });
      })
      .catch(err => {
        console.log('Error fetching notes', err.stack);
        dispatch(receiveError(year, month, err.message));
      });
  };
}

export function getMoreNotes(
  cipherKeyBuf,
  year,
  month,
  page,
  queryObj = {},
  demo = false
) {
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

    return notesService
      .fetch(q, demo)
      .then(res => {
        const { notes, total, prev_date: prevDate } = res;

        const p = decryptNotes(notes, cipherKeyBuf);

        return Promise.all(p).then(notesDec => {
          dispatch(receiveMoreNotes(notesDec, total, year, month, prevDate));
        });
      })
      .catch(err => {
        console.log('err', err);
        // dispatch(receiveError(year, month, err));
      });
  };
}

export function getInitialNotes({ facets, year, month, cipherKeyBuf, demo }) {
  return async (dispatch, getState) => {
    const { notes } = getState();
    const hasError = notes.groups.some(group => {
      return Boolean(group.error);
    });

    if (!notes.prevDate || hasError) {
      return;
    }

    await dispatch(getNotes(cipherKeyBuf, year, month, facets, demo));

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
          cipherKeyBuf,
          demo,
          year: prevYear,
          month: prevMonth
        })
      );
    }
  };
}
