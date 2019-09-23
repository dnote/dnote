export interface ComposerState {
  content: string;
  bookUUID: string;
  bookLabel: string;
}

export const UPDATE_CONTENT = 'composer/UPDATE_CONTENT';
export const UPDATE_BOOK = 'composer/UPDATE_BOOK';
export const RESET = 'composer/RESET';
export const RESET_BOOK = 'composer/RESET_BOOK';

export interface UpdateContentAction {
  type: typeof UPDATE_CONTENT;
  data: {
    content: string;
  };
}

export interface UpdateBookAction {
  type: typeof UPDATE_BOOK;
  data: {
    uuid: string;
    label: string;
  };
}

export interface ResetAction {
  type: typeof RESET;
}

export interface ResetBookAction {
  type: typeof RESET_BOOK;
}

export type ComposerActionType =
  | UpdateContentAction
  | UpdateBookAction
  | ResetAction
  | ResetBookAction;
