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

export default class JoinForm extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      password: '',
      passwordConfirmation: ''
    };
  }

  render() {
    const { onJoin, submitting, onEmailChange, email } = this.props;
    const { password, passwordConfirmation } = this.state;

    return (
      <form
        onSubmit={e => {
          e.preventDefault();

          onJoin(email, password, passwordConfirmation);
        }}
        className="auth-form"
      >
        <div className="input-row">
          <label htmlFor="email-input" className="label">
            Your Email
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
          />
        </div>

        <div className="input-row">
          <label htmlFor="password-input" className="label">
            Create a Password
          </label>
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

        <div className="input-row">
          <label htmlFor="password-input" className="label">
            Confirm your Password
          </label>
          <input
            id="password-input"
            type="password"
            placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;"
            className="form-control"
            value={passwordConfirmation}
            onChange={e => {
              const val = e.target.value;

              this.setState({ passwordConfirmation: val });
            }}
          />
        </div>

        <button
          type="submit"
          className="button button-third button-stretch auth-button"
          disabled={submitting}
        >
          {submitting ? <i className="fa fa-spinner fa-spin" /> : 'Confirm'}
        </button>
      </form>
    );
  }
}
