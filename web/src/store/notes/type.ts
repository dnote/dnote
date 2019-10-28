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

import { NoteData, BookData } from 'jslib/operations/types';
import { RemoteData } from '../types';

export interface NotesState extends RemoteData<NoteData[]> {
  total: number;
}

export const ADD = 'notes/ADD';
export const REFRESH = 'notes/REFRESH';
export const RECEIVE = 'notes/RECEIVE';
export const START_FETCHING = 'notes/START_FETCHING';
export const RECEIVE_ERROR = 'notes/RECEIVE_ERROR';
export const RESET = 'notes/RESET';
export const REMOVE = 'notes/REMOVE';

export interface AddAction {
  type: typeof ADD;
  data: {
    note: NoteData;
  };
}

export interface RefreshAction {
  type: typeof REFRESH;
  data: {
    noteUUID: string;
    book: BookData;
    content: string;
    isPublic: boolean;
  };
}

export interface ReceiveAction {
  type: typeof RECEIVE;
  data: {
    notes: NoteData[];
    total: number;
  };
}

export interface StartFetchingAction {
  type: typeof START_FETCHING;
}

export interface ReceiveErrorAction {
  type: typeof RECEIVE_ERROR;
  data: {
    error: string;
  };
}

export interface ResetAction {
  type: typeof RESET;
}

export interface RemoveAction {
  type: typeof REMOVE;
  data: {
    noteUUID: string;
  };
}

export type NotesActionType =
  | AddAction
  | RefreshAction
  | ReceiveAction
  | StartFetchingAction
  | ReceiveErrorAction
  | ResetAction
  | RemoveAction;
