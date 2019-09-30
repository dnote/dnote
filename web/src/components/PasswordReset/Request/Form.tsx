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
import classnames from 'classnames';

import Button from '../../Common/Button';
import authStyles from '../../Common/Auth.scss';
import styles from './Form.scss';

interface Props {
  onSubmit: (email: string) => void;
  submitting: boolean;
}

const PasswordResetRequestForm: React.SFC<Props> = ({
  onSubmit,
  submitting
}) => {
  const [email, setEmail] = useState('');

  return (
    <form
      onSubmit={e => {
        e.preventDefault();

        onSubmit(email);
      }}
      className={authStyles.form}
    >
      <div className={authStyles['input-row']}>
        <label htmlFor="email-input" className={styles.label}>
          Enter your email and we will send you a link to reset your password
          <input
            id="email-input"
            type="email"
            placeholder="you@example.com"
            className={classnames('form-control', styles['email-input'])}
            value={email}
            onChange={e => {
              const val = e.target.value;

              setEmail(val);
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
        Send password reset email
      </Button>
    </form>
  );
};

export default PasswordResetRequestForm;
