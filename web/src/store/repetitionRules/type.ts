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

import { RepetitionRuleData } from 'jslib/operations/types';
import { RemoteData } from '../types';

export type RepetitionRulesState = RemoteData<RepetitionRuleData[]>;

export const RECEIVE = 'repetitionRules/RECEIVE';
export const RECEIVE_ERROR = 'repetitionRules/RECEIVE_ERROR';
export const ADD = 'repetitionRules/ADD';
export const REMOVE = 'repetitionRules/REMOVE';
export const START_FETCHING = 'repetitionRules/START_FETCHING';
export const FINISH_FETCHING = 'repetitionRules/FINISH_FETCHING';

export interface ReceiveRepetitionRulesAction {
  type: typeof RECEIVE;
  data: {
    repetitionRules: RepetitionRuleData[];
  };
}

export interface ReceiveRepetitionRulesErrorAction {
  type: typeof RECEIVE_ERROR;
  data: {
    err: string;
  };
}

export interface StartFetchingRepetitionRulesAction {
  type: typeof START_FETCHING;
}

export interface FinishFetchingRepetitionRulesAction {
  type: typeof FINISH_FETCHING;
}

export interface AddRepetitionRuleAction {
  type: typeof ADD;
  data: {
    repetitionRule: RepetitionRuleData;
  };
}

export interface RemoveRepetitionRuleAction {
  type: typeof REMOVE;
  data: {
    uuid: string;
  };
}

export type RepetitionRulesActionType =
  | ReceiveRepetitionRulesAction
  | ReceiveRepetitionRulesErrorAction
  | StartFetchingRepetitionRulesAction
  | FinishFetchingRepetitionRulesAction
  | AddRepetitionRuleAction
  | RemoveRepetitionRuleAction;
