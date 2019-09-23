import { Store, Action } from 'redux';
import {
  useDispatch as useReduxDispatch,
  useStore as useReduxStore,
  useSelector as useReduxSelector
} from 'react-redux';
import { ThunkDispatch } from 'redux-thunk';

import { ComposerState } from './composer/types';
import { LocationState } from './location/types';
import { SettingsState } from './settings/types';
import { BooksState } from './books/types';

// AppState represents the application state
interface AppState {
  composer: ComposerState;
  location: LocationState;
  settings: SettingsState;
  books: BooksState;
}

type ReduxDispatch = ThunkDispatch<AppState, any, Action>;

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
