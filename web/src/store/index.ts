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

import { combineReducers, createStore, applyMiddleware, compose } from 'redux';
import thunkMiddleware from 'redux-thunk';

import auth from './auth/reducers';
import form from './form/reducers';
import books from './books/reducers';
import editor from './editor/reducers';
import note from './note/reducers';
import ui from './ui/reducers';
import route from './route/reducers';
import notes from './notes/reducers';

const rootReducer = combineReducers({
  auth,
  form,
  books,
  editor,
  notes,
  note,
  ui,
  route
});

// configuruStore returns a new store that contains the appliation state
export default function configureStore(initialState) {
  const typedWindow = window as any;

  return createStore(
    rootReducer,
    initialState,
    compose(
      applyMiddleware(thunkMiddleware),
      typedWindow.__REDUX_DEVTOOLS_EXTENSION__ &&
        typedWindow.__REDUX_DEVTOOLS_EXTENSION__()
    )
  );
}

export * from './types';
