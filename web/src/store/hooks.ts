/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

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
