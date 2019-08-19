import { Action, Store } from 'redux';
import { ThunkDispatch, ThunkAction } from 'redux-thunk';

import { AuthState } from './auth/type';
import { FormState } from './form/type';
import { BooksState } from './books/type';
import { EditorState } from './editor/type';
import { NoteState } from './note/type';
import { NotesState } from './notes/type';
import { UIState } from './ui/type';
import { RouteState } from './route/type';
import { FiltersState } from './filters/type';

// RemoteData represents a data in Redux store that is fetched from a remote source.
// It contains the state related to the fetching of the data as well as the data itself.
export interface RemoteData<T> {
  isFetching: boolean;
  isFetched: boolean;
  data: T;
  errorMessage: null | string;
}

// AppState represents the application state
export interface AppState {
  auth: AuthState;
  form: FormState;
  books: BooksState;
  editor: EditorState;
  note: NoteState;
  notes: NotesState;
  ui: UIState;
  route: RouteState;
  filters: FiltersState;
}

// ThunkAction is a thunk action type
export type ThunkAction<T = void> = ThunkAction<
  Promise<T>,
  AppState,
  void,
  Action
>;

// AppStore represents the store for the app
export type AppStore = Store<AppState>;

export type ReduxDispatch = ThunkDispatch<AppState, any, Action>;
