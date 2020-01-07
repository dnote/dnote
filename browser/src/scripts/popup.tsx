/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';

import { debounce } from 'jslib/helpers/perf';
import configureStore from './store';
import { loadState, saveState } from './utils/storage';
import App from './components/App';
import ext from './utils/ext';

const appContainer = document.getElementById('app');

loadState(items => {
  if (ext.runtime.lastError) {
    appContainer.innerText = `Failed to retrieve previous app state ${ext.runtime.lastError.message}`;
    return;
  }

  let initialState;
  const prevState = items.state;
  if (prevState) {
    // rehydrate
    initialState = prevState;
  }

  const store = configureStore(initialState);

  store.subscribe(
    debounce(() => {
      const state = store.getState();

      saveState(state);
    }, 100)
  );

  ReactDOM.render(
    <Provider store={store}>
      <App />
    </Provider>,
    appContainer,
    () => {
      // On Chrome, popup window size is kept at minimum if app render is delayed
      // Therefore add minimum dimension to body until app is rendered
      document.getElementsByTagName('body')[0].className = '';
    }
  );
});
