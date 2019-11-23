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
import { NoteData } from 'jslib/operations/types';

export type NoteState = RemoteData<NoteData>;

export const RECEIVE = 'note/RECEIVE';
export const START_FETCHING = 'note/START_FETCHING';
export const ERROR = 'note/ERROR';
export const RESET = 'note/RESET';

export interface ReceiveNote {
  type: typeof RECEIVE;
  data: {
    note: NoteData;
  };
}

export interface StartFetchingNote {
  type: typeof START_FETCHING;
}

export interface ResetNote {
  type: typeof RESET;
}

export interface ReceiveNoteError {
  type: typeof ERROR;
  data: {
    errorMessage: string;
  };
}

export type NoteActionType =
  | ReceiveNote
  | StartFetchingNote
  | ReceiveNoteError
  | ResetNote;
