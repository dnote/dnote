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

import {
  UPDATE_CONTENT,
  UPDATE_BOOK,
  RESET,
  RESET_BOOK,
  UpdateContentAction,
  UpdateBookAction,
  ResetBookAction,
  ResetAction
} from './types';

export function updateContent(content: string): UpdateContentAction {
  return {
    type: UPDATE_CONTENT,
    data: { content }
  };
}

export interface UpdateBookActionParam {
  uuid: string;
  label: string;
}

export function updateBook({
  uuid,
  label
}: UpdateBookActionParam): UpdateBookAction {
  return {
    type: UPDATE_BOOK,
    data: {
      uuid,
      label
    }
  };
}

export function resetBook(): ResetBookAction {
  return {
    type: RESET_BOOK
  };
}

export function resetComposer(): ResetAction {
  return {
    type: RESET
  };
}
