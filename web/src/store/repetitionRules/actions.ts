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

import services from 'web/libs/services';
import { RepetitionRuleData } from 'jslib/operations/types';
import { CreateParams } from 'jslib/services/repetitionRules';
import {
  RECEIVE,
  ADD,
  REMOVE,
  START_FETCHING,
  FINISH_FETCHING,
  RECEIVE_ERROR,
  ReceiveRepetitionRulesAction,
  ReceiveRepetitionRulesErrorAction,
  StartFetchingRepetitionRulesAction,
  FinishFetchingRepetitionRulesAction,
  AddRepetitionRuleAction,
  RemoveRepetitionRuleAction
} from './type';
import { ThunkAction } from '../types';

function receiveRepetitionRules(
  repetitionRules: RepetitionRuleData[]
): ReceiveRepetitionRulesAction {
  return {
    type: RECEIVE,
    data: { repetitionRules }
  };
}

function receiveRepetitionRulesError(err: string): ReceiveRepetitionRulesErrorAction {
  return {
    type: RECEIVE_ERROR,
    data: { err }
  };
}

function startFetchingRepetitionRules(): StartFetchingRepetitionRulesAction {
  return {
    type: START_FETCHING
  };
}

function finishFetchingRepetitionRules(): FinishFetchingRepetitionRulesAction {
  return {
    type: FINISH_FETCHING
  };
}

export const getRepetitionRules = (): ThunkAction<void> => {
  return dispatch => {
    dispatch(startFetchingRepetitionRules());

    return services.repetitionRules
      .fetchAll()
      .then(data => {
        dispatch(receiveRepetitionRules(data));
        dispatch(finishFetchingRepetitionRules());
      })
      .catch(err => {
        console.log('getRepetitionRules error', err);
        dispatch(receiveRepetitionRulesError(err));
      });
  };
};

export function addRepetitionRule(repetitionRule: RepetitionRuleData): AddRepetitionRuleAction {
  return {
    type: ADD,
    data: { repetitionRule }
  };
}

export const createRepetitionRule = (
  p: CreateParams
): ThunkAction<RepetitionRuleData> => {
  return dispatch => {
    return services.repetitionRules.create(p).then(data => {
      dispatch(addRepetitionRule(data));

      return data;
    });
  };
};

export function removeRepetitionRule(uuid: string): RemoveRepetitionRuleAction {
  return {
    type: REMOVE,
    data: { uuid }
  };
}
