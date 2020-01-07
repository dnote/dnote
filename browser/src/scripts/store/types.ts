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

import { Action } from 'redux';
import { ThunkAction as ReduxThunkAction } from 'redux-thunk';

import { AuthState } from './auth/types';
import { ComposerState } from './composer/types';
import { LocationState } from './location/types';
import { SettingsState } from './settings/types';
import { BooksState } from './books/types';

// AppState represents the application state
export interface AppState {
  auth: AuthState;
  composer: ComposerState;
  location: LocationState;
  settings: SettingsState;
  books: BooksState;
}

// ThunkAction is a thunk action type
export type ThunkAction<T = void> = ReduxThunkAction<
  Promise<T>,
  AppState,
  void,
  Action
>;
