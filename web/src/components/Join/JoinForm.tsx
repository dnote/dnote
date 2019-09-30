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

import Button from '../Common/Button';
import { useDispatch } from '../../store';
import { updateAuthEmail } from '../../store/form';
import authStyles from '../Common/Auth.scss';

interface Props {
  onJoin: (
    email: string,
    password: string,
    passwordConfirmation: string
  ) => void;
  submitting: boolean;
  email?: string;
  cta?: string;
}

const JoinForm: React.SFC<Props> = ({
  onJoin,
  submitting,
  email,
  cta = 'Join'
}) => {
  const [password, setPassword] = useState('');
  const [passwordConfirmation, setPasswordConfirmation] = useState('');
  const dispatch = useDispatch();

  return (
    <form
      onSubmit={e => {
        e.preventDefault();

        onJoin(email, password, passwordConfirmation);
      }}
      id="T-join-form"
      className={authStyles.form}
    >
      {email !== undefined && (
        <div className={authStyles['input-row']}>
          <label htmlFor="email-input" className={authStyles.label}>
            Your email
            <input
              autoFocus
              id="email-input"
              type="email"
              placeholder="you@example.com"
              className="form-control"
              value={email}
              onChange={e => {
                const val = e.target.value;

                dispatch(updateAuthEmail(val));
              }}
            />
          </label>
        </div>
      )}

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

      <div className={authStyles['input-row']}>
        <label
          htmlFor="password-confirmation-input"
          className={authStyles.label}
        >
          Password confirmation
          <input
            id="password-confirmation-input"
            type="password"
            placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;"
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
        kind="third"
        size="normal"
        stretch
        className={authStyles['auth-button']}
        isBusy={submitting}
      >
        {cta}
      </Button>
    </form>
  );
};

export default JoinForm;
