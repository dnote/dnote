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
import Helmet from 'react-helmet';
import { Link } from 'react-router-dom';

import { getLoginPath } from 'web/libs/paths';
import services from 'web/libs/services';
import Flash from '../../Common/Flash';
import Form from './Form';
import Logo from '../../Icons/Logo';
import authStyles from '../../Common/Auth.scss';
import styles from './Request.scss';

interface Props {}

const PasswordResetRequest: React.SFC<Props> = () => {
  const [errorMsg, setErrorMsg] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [processed, setProcessed] = useState(false);

  function onSubmit(email) {
    if (!email) {
      setErrorMsg('Please enter email');
      return;
    }

    setErrorMsg('');
    setSubmitting(true);

    services.users
      .sendResetPasswordEmail({ email })
      .then(() => {
        setSubmitting(false);
        setProcessed(true);
      })
      .catch(err => {
        setSubmitting(false);
        setErrorMsg(err.message);
      });
  }

  return (
    <div className={authStyles.page}>
      <Helmet>
        <title>Reset Password</title>
      </Helmet>

      <div className="container">
        <Link to="/">
          <Logo fill="#252833" width={60} height={60} />
        </Link>
        <h1 className={authStyles.heading}>Reset Password</h1>

        <div className={authStyles.body}>
          <div className={authStyles.panel}>
            <Flash kind="danger" when={errorMsg !== ''}>
              {errorMsg}
            </Flash>

            {processed ? (
              <div>
                <div className={styles['success-msg']}>
                  Check your email for a link to reset your password.
                </div>
                <Link
                  to={getLoginPath()}
                  className="button button-first button-normal button-stretch"
                >
                  Back to login
                </Link>
              </div>
            ) : (
              <Form onSubmit={onSubmit} submitting={submitting} />
            )}
          </div>

          <div className={authStyles.footer}>
            <div className={authStyles.callout}>Remember your password?</div>
            <Link to={getLoginPath()} className="auth-cta">
              Login
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PasswordResetRequest;
