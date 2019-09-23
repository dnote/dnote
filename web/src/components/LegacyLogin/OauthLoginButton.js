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

import React from 'react';

import { handleLogin } from 'jslib/helpers/auth';
import google from '../../img/google.png';
import github from '../../img/github.png';

export default class OauthLoginButton extends React.Component {
  getLogo = () => {
    const { provider } = this.props;

    switch (provider) {
      case 'github':
        return github;
      case 'gplus':
        return google;
      default:
        return null;
    }
  };

  render() {
    const { referrer, provider, text } = this.props;

    return (
      <button
        type="button"
        className="button oauth-button"
        onClick={() => {
          handleLogin({ provider, referrer });
        }}
      >
        <span className="oauth-button-content">
          <img src={this.getLogo()} alt={provider} className="provider-logo" />
          <span className="oauth-text">{text}</span>
        </span>
      </button>
    );
  }
}
