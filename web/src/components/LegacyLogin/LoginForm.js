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

/* eslint-disable jsx-a11y/label-has-associated-control */

import React from 'react';

export default class LoginForm extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      password: ''
    };
  }

  render() {
    const { email, onLogin, onEmailChange, submitting } = this.props;
    const { password } = this.state;

    return (
      <form
        onSubmit={e => {
          e.preventDefault();

          onLogin(email, password);
        }}
        className="auth-form"
      >
        <div className="input-row">
          <label htmlFor="email-input" className="label">
            Email
          </label>
          <input
            id="email-input"
            type="email"
            placeholder="you@example.com"
            className="form-control"
            value={email}
            onChange={e => {
              const val = e.target.value;

              onEmailChange(val);
            }}
            autoComplete="on"
          />
        </div>

        <div className="input-row">
          <div className="label-row">
            <label htmlFor="password-input" className="label">
              Password
            </label>
          </div>
          <input
            id="password-input"
            type="password"
            placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;"
            className="form-control"
            value={password}
            onChange={e => {
              const val = e.target.value;

              this.setState({ password: val });
            }}
          />
        </div>

        <button
          type="submit"
          className="button button-first button-stretch auth-button"
          disabled={submitting}
        >
          {submitting ? <i className="fa fa-spinner fa-spin" /> : 'Sign in'}
        </button>
      </form>
    );
  }
}
