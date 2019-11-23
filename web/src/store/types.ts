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

import { Action, Store } from 'redux';
import { ThunkDispatch, ThunkAction } from 'redux-thunk';

import { AuthState } from './auth/type';
import { FormState } from './form/type';
import { BooksState } from './books/type';
import { EditorState } from './editor/type';
import { NoteState } from './note/type';
import { NotesState } from './notes/type';
import { UIState } from './ui/type';
import { RouteState } from './route/type';
import { FiltersState } from './filters/type';
import { RepetitionRulesState } from './repetitionRules/type';

// RemoteData represents a data in Redux store that is fetched from a remote source.
// It contains the state related to the fetching of the data as well as the data itself.
export interface RemoteData<T> {
  isFetching: boolean;
  isFetched: boolean;
  data: T;
  errorMessage: string;
}

// AppState represents the application state
export interface AppState {
  auth: AuthState;
  form: FormState;
  books: BooksState;
  editor: EditorState;
  note: NoteState;
  notes: NotesState;
  ui: UIState;
  route: RouteState;
  filters: FiltersState;
  repetitionRules: RepetitionRulesState;
}

// ThunkAction is a thunk action type
export type ThunkAction<T = void> = ThunkAction<
  Promise<T>,
  AppState,
  void,
  Action
>;

// AppStore represents the store for the app
export type AppStore = Store<AppState>;

export type ReduxDispatch = ThunkDispatch<AppState, any, Action>;
