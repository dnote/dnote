import { Store } from 'redux';
import {
  useDispatch as useReduxDispatch,
  useStore as useReduxStore,
  useSelector as useReduxSelector
} from 'react-redux';

import { ReduxDispatch, AppState } from './types';
import { FiltersState } from './filters/type';

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

// custom hooks
export function useFilters(): FiltersState {
  const { filters } = useSelector(state => {
    return {
      filters: state.filters
    };
  });

  return filters;
}
