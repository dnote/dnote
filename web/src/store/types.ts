import { Action, Store } from 'redux';
import { ThunkDispatch, ThunkAction } from 'redux-thunk';
import {
  useDispatch as useReduxDispatch,
  useStore as useReduxStore,
  useSelector as useReduxSelector
} from 'react-redux';

import { AuthState } from './auth/type';
import { FormState } from './form/type';
import { BooksState } from './books/type';
import { EditorState } from './editor/type';
import { NoteState } from './note/type';
import { NotesState } from './notes/type';
import { UIState } from './ui/type';
import { RouteState } from './route/type';

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

export function useDispatch(): ReduxDispatch {
  return useReduxDispatch<ReduxDispatch>();
}

export function useStore(): Store<AppState> {
  return useReduxStore<AppState>();
}

export function useSelector<TSelected>(
  selector: (state: AppState) => TSelected
) {
  return useReduxSelector<AppState, TSelected>(selector);
}
