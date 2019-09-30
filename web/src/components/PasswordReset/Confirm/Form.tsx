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

import authStyles from '../../Common/Auth.scss';
import Button from '../../Common/Button';

interface Props {
  onResetPassword: (password: string, passwordConfirmation: string) => void;
  submitting: boolean;
}

const PasswordResetConfirmForm: React.SFC<Props> = ({
  onResetPassword,
  submitting
}) => {
  const [password, setPassword] = useState('');
  const [passwordConfirmation, setPasswordConfirmation] = useState('');

  return (
    <form
      onSubmit={e => {
        e.preventDefault();

        onResetPassword(password, passwordConfirmation);
      }}
      className={authStyles.form}
    >
      <div className={authStyles['input-row']}>
        <label htmlFor="password-input" className={authStyles.label}>
          Password
          <input
            id="password-input"
            type="password"
            placeholder="********"
            className="form-control"
            value={password}
            onChange={e => {
              const val = e.target.value;

              setPassword(val);
            }}
          />
        </label>
      </div>

      <div className={authStyles['input-row']}>
        <label
          htmlFor="password-confirmation-input"
          className={authStyles.label}
        >
          Password Confirmation
          <input
            id="password-confirmation-input"
            type="password"
            placeholder="********"
            className="form-control"
            value={passwordConfirmation}
            onChange={e => {
              const val = e.target.value;

              setPasswordConfirmation(val);
            }}
          />
        </label>
      </div>

      <Button
        type="submit"
        size="normal"
        kind="first"
        stretch
        className={authStyles['auth-button']}
        isBusy={submitting}
      >
        Reset password
      </Button>
    </form>
  );
};

export default PasswordResetConfirmForm;
