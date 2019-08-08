import { RemoteData } from '../types';
import { NoteData } from '../../operations/types';

export type NoteState = RemoteData<NoteData>;

export const RECEIVE = 'note/RECEIVE';
export const START_FETCHING = 'note/START_FETCHING';
export const ERROR = 'note/ERROR';
export const RESET = 'note/RESET';

export interface ReceiveNote {
  type: typeof RECEIVE;
  data: {
    note: NoteData;
  };
}

export interface StartFetchingNote {
  type: typeof START_FETCHING;
}

export interface ResetNote {
  type: typeof RESET;
}

export interface ReceiveNoteError {
  type: typeof ERROR;
  data: {
    errorMessage: string;
  };
}

export type NoteActionType =
  | ReceiveNote
  | StartFetchingNote
  | ReceiveNoteError
  | ResetNote;
