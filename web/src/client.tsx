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

import 'core-js/stable';
import 'regenerator-runtime/runtime';

import React from 'react';
import { render } from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { Provider } from 'react-redux';

import { debounce } from 'jslib/helpers/perf';
import App from './components/App';
import configureStore from './store';
import { markPersisted } from './store/editor';
import { loadState, saveState } from './libs/localStorage';
import './libs/restoreScroll';

const persistedState = loadState();
const store = configureStore(persistedState);

store.subscribe(
  debounce(() => {
    const state = store.getState();

    saveState({
      editor: state.editor
    });

    if (!state.editor.persisted) {
      store.dispatch(markPersisted());
    }
  }, 50)
);

function renderApp() {
  render(
    <Provider store={store}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </Provider>,
    document.getElementById('app')
  );
}

const splashEl = document.getElementById('splash');
if (splashEl && splashEl.parentNode) {
  splashEl.parentNode.removeChild(splashEl);
}
renderApp();
