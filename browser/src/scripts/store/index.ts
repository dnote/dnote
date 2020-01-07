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

import { combineReducers, createStore, applyMiddleware, compose } from 'redux';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';

import location from './location/reducers';
import settings from './settings/reducers';
import books from './books/reducers';
import composer from './composer/reducers';
import auth from './auth/reducers';
import { AppState } from './types';
import config from '../utils/config';

const rootReducer = combineReducers({
  auth,
  location,
  settings,
  books,
  composer
});

// initState returns a new state with any missing values populated
// if a state is given.
function initState(s: AppState | undefined): AppState {
  if (s === undefined) {
    return undefined;
  }

  const { settings: settingsState } = s;

  return {
    ...s,
    settings: {
      ...settingsState,
      apiUrl: settingsState.apiUrl || config.defaultApiEndpoint,
      webUrl: settingsState.webUrl || config.defaultWebUrl
    }
  };
}

// configureStore returns a new store that contains the appliation state
export default function configureStore(state: AppState | undefined) {
  const typedWindow = window as any;

  const composeEnhancers =
    typedWindow.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

  return createStore(
    rootReducer,
    initState(state),
    composeEnhancers(applyMiddleware(createLogger, thunkMiddleware))
  );
}
