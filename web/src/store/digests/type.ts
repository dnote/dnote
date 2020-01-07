/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

import { DigestData } from 'jslib/operations/types';

import { RemoteData } from '../types';

export interface DigestsState extends RemoteData<DigestData[]> {
  total: number;
  page: number;
}

export const START_FETCHING = 'digests/START_FETCHING';
export const RECEIVE = 'digests/RECEIVE';
export const RECEIVE_ERROR = 'digests/RECEIVE_ERROR';
export const RESET = 'digests/RESET';

export interface StartFetchingAction {
  type: typeof START_FETCHING;
}

export interface ResetAction {
  type: typeof RESET;
}

export interface ReceiveErrorAction {
  type: typeof RECEIVE_ERROR;
  data: {
    error: string;
  };
}

export interface ReceiveAction {
  type: typeof RECEIVE;
  data: {
    items: DigestData[];
    total: number;
    page: number;
  };
}

export type DigestsActionType =
  | StartFetchingAction
  | ReceiveAction
  | ReceiveErrorAction
  | ResetAction;
