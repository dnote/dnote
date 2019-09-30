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

export interface MessageData {
  content: string;
  kind: string;
}

export interface MessageState {
  [path: string]: MessageData;
}

export interface UIState {
  message: MessageState;
}

export const SET_MESSAGE = 'ui/SET_MESSAGE';
export const UNSET_MESSAGE = 'ui/UNSET_MESSAGE';

export interface SetMessageAction {
  type: typeof SET_MESSAGE;
  data: {
    message: string;
    kind: string;
    path: string;
  };
}

export interface UnsetMessageAction {
  type: typeof UNSET_MESSAGE;
  data: {
    path: string;
  };
}

export type UIActionType = SetMessageAction | UnsetMessageAction;
