export type BookData = any;

export interface BooksState {
  items: BookData[];
  isFetching: boolean;
  isFetched: boolean;
  error: string | null;
}

export const START_FETCHING = 'books/START_FETCHING';
export const RECEIVE = 'books/RECEIVE';
export const RECEIVE_ERROR = 'books/RECEIVE_ERROR';

export interface StartFetchingAction {
  type: typeof START_FETCHING;
}

export interface ReceiveAction {
  type: typeof RECEIVE;
  data: {
    books: BookData[];
  };
}

export interface ReceiveErrorAction {
  type: typeof RECEIVE_ERROR;
  data: {
    error: string;
  };
}

export type BooksActionType =
  | StartFetchingAction
  | ReceiveAction
  | ReceiveErrorAction;
