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

import React, { useState } from 'react';

import Flash from './Flash';
import { updateSettings, resetSettings } from '../store/settings/actions';
import { useDispatch, useSelector, useStore } from '../store/hooks';

interface Props {}

// isValidURL checks if the given string is a valid URL
function isValidURL(url: string): boolean {
  const a = document.createElement('a');
  a.href = url;
  return a.host && a.host !== window.location.host;
}

// validateFormState validates the given form state. If any input is
// invalid, it throws an error.
function validateFormState({ apiUrl, webUrl }) {
  if (!isValidURL(apiUrl)) {
    throw new Error('Invalid URL for the API URL');
  }

  if (!isValidURL(webUrl)) {
    throw new Error('Invalid URL for the web URL');
  }
}

const Settings: React.FunctionComponent<Props> = () => {
  const { settings } = useSelector(state => ({
    settings: state.settings
  }));
  const store = useStore();

  const [apiUrl, setAPIUrl] = useState(settings.apiUrl);
  const [webUrl, setWebUrl] = useState(settings.webUrl);
  const [errMsg, setErrMsg] = useState('');
  const [successMsg, setSuccessMsg] = useState('');
  const dispatch = useDispatch();

  function handleRestore() {
    dispatch(resetSettings());
    setSuccessMsg('Restored the default settings');

    const { settings: settingsState } = store.getState();

    setAPIUrl(settingsState.apiUrl);
    setWebUrl(settingsState.webUrl);
  }

  function handleSubmit(e) {
    e.preventDefault();
    setSuccessMsg('');
    setErrMsg('');

    try {
      validateFormState({ apiUrl, webUrl });
    } catch (err) {
      setErrMsg(err.message);
      return;
    }

    dispatch(
      updateSettings({
        apiUrl,
        webUrl
      })
    );
    setSuccessMsg('Succesfully updated the settings.');
  }

  return (
    <div>
      <Flash kind="error" when={errMsg !== ''} message={errMsg} />
      <Flash kind="info" when={successMsg !== ''} message={successMsg} />

      <div className="settings page">
        <h1 className="heading">Settings</h1>

        <p className="lead">Customize your Dnote extension</p>

        <form id="settings-form" onSubmit={handleSubmit}>
          <div className="input-row">
            <label htmlFor="api-url-input" className="label">
              API URL
            </label>

            <input
              type="api-url"
              placeholder="https://api.getdnote.com"
              className="input"
              id="api-url-input"
              value={apiUrl}
              onChange={e => {
                setAPIUrl(e.target.value);
              }}
            />
          </div>

          <div className="input-row">
            <label htmlFor="web-url-input" className="label">
              Web URL
            </label>

            <input
              type="web-url"
              placeholder="https://app.getdnote.com"
              className="input"
              id="web-url-input"
              value={webUrl}
              onChange={e => {
                setWebUrl(e.target.value);
              }}
            />
          </div>

          <div className="actions">
            <button
              type="submit"
              className="button button-first button-small button-stretch"
            >
              Save
            </button>

            <button
              type="button"
              onClick={handleRestore}
              className="restore button-no-ui"
            >
              Restore default
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Settings;
