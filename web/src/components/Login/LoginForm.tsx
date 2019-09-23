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
import { Link } from 'react-router-dom';

import styles from '../Common/Auth.scss';
import Button from '../Common/Button';
import { getPasswordResetRequestPath } from 'web/libs/paths';

interface Props {
  email: string;
  submitting: boolean;
  onLogin: (email: string, password: string) => void;
  onUpdateEmail: (string) => void;
}

const LoginForm: React.SFC<Props> = ({
  onLogin,
  submitting,
  email,
  onUpdateEmail
}) => {
  const [password, setPassword] = useState('');

  return (
    <form
      onSubmit={e => {
        e.preventDefault();

        onLogin(email, password);
      }}
      id="T-login-form"
      className={styles.form}
    >
      <div className={styles['input-row']}>
        <label htmlFor="email-input" className={styles.label}>
          Email
          <input
            tabIndex={1}
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

      <div className={styles['input-row']}>
        <label htmlFor="password-input" className={styles.label}>
          Password
          <Link to={getPasswordResetRequestPath()} className={styles.forgot}>
            Forgot?
          </Link>
          <input
            tabIndex={2}
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
        tabIndex={3}
        type="submit"
        size="normal"
        kind="first"
        stretch
        className={styles['auth-button']}
        isBusy={submitting}
      >
        Sign in
      </Button>
    </form>
  );
};

export default LoginForm;
