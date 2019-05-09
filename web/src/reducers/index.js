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

import { combineReducers } from 'redux';

import auth from './auth';
import notes from './notes';
import note from './note';
import books from './books';
import calendar from './calendar';
import ui from './ui';
import form from './form';
import editor from './editor';
import digest from './digest';

const rootReducer = combineReducers({
  auth,
  notes,
  note,
  books,
  calendar,
  ui,
  form,
  editor,
  digest
});

export default rootReducer;
