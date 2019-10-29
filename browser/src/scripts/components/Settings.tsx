/* Copyright (C) 2019 Monomax Software Pty Ltd
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
import { connect } from 'react-redux';
import { findDOMNode } from 'react-dom';

import Link from './Link';
import Flash from './Flash';
import config from '../utils/config';
import { updateSettings, login } from '../store/settings/actions';
import { useDispatch, useSelector } from '../store/hooks';
import services from '../utils/services';

interface Props {}

function isValidURL(url: string): boolean {
  var a = document.createElement('a');
  a.href = url;
  return a.host && a.host != window.location.host;
}

function validateFormState({ apiUrl, webUrl }) {
  if (!isValidURL(apiUrl)) {
    throw new Error('Invalid URL for the API URL');
  }

  if (!isValidURL(webUrl)) {
    throw new Error('Invalid URL for the web URL');
  }
}

const Settings: React.FunctionComponent<Props> = () => {
  const { settings } = useSelector(state => {
    return {
      settings: state.settings
    };
  });

  const [apiUrl, setAPIUrl] = useState(settings.apiUrl);
  const [webUrl, setWebUrl] = useState(settings.webUrl);
  const [errMsg, setErrMsg] = useState('');
  const [successMsg, setSuccessMsg] = useState('');
  const dispatch = useDispatch();

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

        <p className="lead">Customize Dnote browser extension</p>

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
          </div>
        </form>
      </div>
    </div>
  );
};

export default Settings;
