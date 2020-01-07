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

export interface ComposerState {
  content: string;
  bookUUID: string;
  bookLabel: string;
}

export const UPDATE_CONTENT = 'composer/UPDATE_CONTENT';
export const UPDATE_BOOK = 'composer/UPDATE_BOOK';
export const RESET = 'composer/RESET';
export const RESET_BOOK = 'composer/RESET_BOOK';

export interface UpdateContentAction {
  type: typeof UPDATE_CONTENT;
  data: {
    content: string;
  };
}

export interface UpdateBookAction {
  type: typeof UPDATE_BOOK;
  data: {
    uuid: string;
    label: string;
  };
}

export interface ResetAction {
  type: typeof RESET;
}

export interface ResetBookAction {
  type: typeof RESET_BOOK;
}

export type ComposerActionType =
  | UpdateContentAction
  | UpdateBookAction
  | ResetAction
  | ResetBookAction;
