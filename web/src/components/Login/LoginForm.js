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

import React, { useState } from 'react';

import authStyles from '../Common/Auth.module.scss';
import Button from '../Common/Button';

function LoginForm({ onLogin, submitting, email, onUpdateEmail }) {
  const [password, setPassword] = useState('');

  return (
    <form
      onSubmit={e => {
        e.preventDefault();

        onLogin(email, password);
      }}
      id="T-login-form"
      className={authStyles.form}
    >
      <div className={authStyles['input-row']}>
        <label htmlFor="email-input" className={authStyles.label}>
          Email
          <input
            id="email-input"
            type="email"
            placeholder="you@example.com"
            className="form-control"
            value={email}
            onChange={e => {
              const val = e.target.value;

              onUpdateEmail(val);
            }}
            autoComplete="on"
            autoFocus
          />
        </label>
      </div>

      <div className={authStyles['input-row']}>
        <label htmlFor="password-input" className={authStyles.label}>
          Password
          <input
            id="password-input"
            type="password"
            placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;"
            className="form-control"
            value={password}
            onChange={e => {
              const val = e.target.value;

              setPassword(val);
            }}
          />
        </label>
      </div>

      <Button
        type="submit"
        kind="first"
        stretch
        className={authStyles['auth-button']}
        isBusy={submitting}
      >
        Sign in
      </Button>
    </form>
  );
}

export default LoginForm;
