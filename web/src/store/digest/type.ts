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

export type DigestState = RemoteData<DigestData>;

export const RECEIVE = 'digest/RECEIVE';
export const START_FETCHING = 'digest/START_FETCHING';
export const ERROR = 'digest/ERROR';
export const RESET = 'digest/RESET';
export const SET_NOTE_REVIEWED = 'digest/SET_NOTE_REVIEWED';

export interface ReceiveDigest {
  type: typeof RECEIVE;
  data: {
    digest: DigestData;
  };
}

export interface StartFetchingDigest {
  type: typeof START_FETCHING;
}

export interface ResetDigest {
  type: typeof RESET;
}

export interface ReceiveDigestError {
  type: typeof ERROR;
  data: {
    errorMessage: string;
  };
}

export interface SetNoteReviewed {
  type: typeof SET_NOTE_REVIEWED;
  data: {
    noteUUID: string;
    isReviewed: boolean;
  };
}

export type DigestActionType =
  | ReceiveDigest
  | StartFetchingDigest
  | ReceiveDigestError
  | ResetDigest
  | SetNoteReviewed;
