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
