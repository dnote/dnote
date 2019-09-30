import { NoteData } from 'jslib/operations/types';
import { BookData } from 'jslib/operations/books';
import { RemoteData } from '../types';

export interface NotesState extends RemoteData<NoteData[]> {
  total: number;
}

export const ADD = 'notes/ADD';
export const REFRESH = 'notes/REFRESH';
export const RECEIVE = 'notes/RECEIVE';
export const START_FETCHING = 'notes/START_FETCHING';
export const RECEIVE_ERROR = 'notes/RECEIVE_ERROR';
export const RESET = 'notes/RESET';
export const REMOVE = 'notes/REMOVE';

export interface AddAction {
  type: typeof ADD;
  data: {
    note: NoteData;
  };
}

export interface RefreshAction {
  type: typeof REFRESH;
  data: {
    noteUUID: string;
    book: BookData;
    content: string;
    isPublic: boolean;
  };
}

export interface ReceiveAction {
  type: typeof RECEIVE;
  data: {
    notes: NoteData[];
    total: number;
  };
}

export interface StartFetchingAction {
  type: typeof START_FETCHING;
}

export interface ReceiveErrorAction {
  type: typeof RECEIVE_ERROR;
  data: {
    error: string;
  };
}

export interface ResetAction {
  type: typeof RESET;
}

export interface RemoveAction {
  type: typeof REMOVE;
  data: {
    noteUUID: string;
  };
}

export type NotesActionType =
  | AddAction
  | RefreshAction
  | ReceiveAction
  | StartFetchingAction
  | ReceiveErrorAction
  | ResetAction
  | RemoveAction;
