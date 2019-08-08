import { NoteData } from '../../operations/types';
import { BookData } from '../../operations/books';

export interface NotesGroupItems {
  [uuid: string]: NoteData;
}

export interface NotesGroup {
  year: number;
  month: number;
  uuids: string[];
  items: NotesGroupItems;
  total: number;
  isFetching: boolean;
  isFetched: boolean;
  isFetchingMore: boolean;
  hasFetchedMore: boolean;
  page: number;
  error: string;
}

export interface NotesState {
  groups: NotesGroup[];
  initialized: boolean;
  prevDate: number;
}

//= RemoteData<NoteData>;

export const ADD = 'notes/ADD';
export const REFRESH = 'notes/REFRESH';
export const RECEIVE = 'notes/RECEIVE';
export const RECEIVE_MORE = 'notes/RECEIVE_MORE';
export const START_FETCHING = 'notes/START_FETCHING';
export const START_FETCHING_MORE = 'notes/START_FETCHING_MORE';
export const RECEIVE_ERROR = 'notes/RECEIVE_ERROR';
export const RESET = 'notes/RESET';
export const REMOVE = 'notes/REMOVE';

export interface AddAction {
  type: typeof ADD;
  data: {
    note: NoteData;
    year: number;
    month: number;
  };
}

export interface RefreshAction {
  type: typeof REFRESH;
  data: {
    year: number;
    month: number;
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
    year: number;
    month: number;
    prevDate: number;
  };
}

export interface ReceiveMoreAction {
  type: typeof RECEIVE_MORE;
  data: {
    notes: NoteData[];
    year: number;
    month: number;
    prevDate: number;
  };
}

export interface StartFetchingAction {
  type: typeof START_FETCHING;
  data: {
    year: number;
    month: number;
  };
}

export interface StartFetchingMoreAction {
  type: typeof START_FETCHING_MORE;
  data: {
    year: number;
    month: number;
  };
}

export interface ReceiveErrorAction {
  type: typeof RECEIVE_ERROR;
  data: {
    year: number;
    month: number;
    error: string;
  };
}

export interface ResetAction {
  type: typeof RESET;
}

export interface RemoveAction {
  type: typeof REMOVE;
  data: {
    year: number;
    month: number;
    noteUUID: string;
  };
}

export type NotesActionType =
  | AddAction
  | RefreshAction
  | ReceiveAction
  | ReceiveMoreAction
  | StartFetchingAction
  | StartFetchingMoreAction
  | ReceiveErrorAction
  | ResetAction
  | RemoveAction;
