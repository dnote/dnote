import { RemoteData } from '../types';
import { BookData } from 'jslib/operations/books';

export type BooksState = RemoteData<BookData[]>;

export const RECEIVE = 'books/RECEIVE';
export const ADD = 'books/ADD';
export const REMOVE = 'books/REMOVE';
export const START_FETCHING = 'books/START_FETCHING';
export const FINISH_FETCHING = 'books/FINISH_FETCHING';

export interface ReceiveBooks {
  type: typeof RECEIVE;
  data: {
    books: BookData[];
  };
}

export interface StartFetchingBooks {
  type: typeof START_FETCHING;
}

export interface FinishFetchingBooks {
  type: typeof FINISH_FETCHING;
}

export interface AddBook {
  type: typeof ADD;
  data: {
    book: BookData;
  };
}

export interface RemoveBook {
  type: typeof REMOVE;
  data: {
    bookUUID: string;
  };
}

export type BooksActionType =
  | ReceiveBooks
  | StartFetchingBooks
  | FinishFetchingBooks
  | AddBook
  | RemoveBook;
