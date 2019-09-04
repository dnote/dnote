import React, { useState } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router-dom';

import Flash from '../../Common/Flash';
import Form from './Form';
import Logo from '../../Icons/Logo';
import * as usersService from '../../../services/users';
import authStyles from '../../Common/Auth.scss';
import { getLoginPath } from '../../../libs/paths';
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

    usersService
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
