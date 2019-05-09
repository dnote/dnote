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
import { connect } from 'react-redux';
import Helmet from 'react-helmet';
import { Link } from 'react-router-dom';

import LoginForm from './LoginForm';

import { getReferrer } from '../../libs/url';
import Logo from '../Icons/Logo';
import Flash from '../Common/Flash';
import { presignin, signin } from '../../services/users';
import { loginHelper, aes256GcmDecrypt } from '../../crypto';
import { bufToB64, b64ToBuf } from '../../libs/encoding';
import { getCurrentUser } from '../../actions/auth';
import { updateAuthEmail } from '../../actions/form';

import authStyles from '../Common/Auth.module.scss';

function Login({ doGetCurrentUser, doUpdateAuthEmail, formData, location }) {
  const [errMsg, setErrMsg] = useState('');
  const [submitting, setSubmitting] = useState(false);

  async function handleLogin(email, password) {
    if (!email) {
      setErrMsg('Please enter email');
      return;
    }
    if (!password) {
      setErrMsg('Please enter password');
      return;
    }

    setErrMsg('');
    setSubmitting(true);

    try {
      const { iteration } = await presignin({ email, password });

      if (iteration === 0) {
        throw new Error('Please login from /app/legacy/login');
      }

      const { masterKey, authKey } = await loginHelper({
        email,
        password,
        iteration
      });
      const signinResp = await signin({ email, authKey });

      const cipherKey = await aes256GcmDecrypt(
        b64ToBuf(masterKey),
        b64ToBuf(signinResp.cipher_key_enc)
      );
      localStorage.setItem('cipherKey', bufToB64(cipherKey));

      // guestOnly HOC will redirect the user accordingly after the current user is fetched
      await doGetCurrentUser();
      doUpdateAuthEmail('');
    } catch (err) {
      console.log(err);
      setErrMsg(err.message);
      setSubmitting(false);
    }
  }
  const referrer = getReferrer(location);

  return (
    <div className={authStyles.page}>
      <Helmet>
        <title>Sign In</title>
      </Helmet>
      <div className="container">
        <a href="/">
          <Logo fill="#252833" width="60" height="60" />
        </a>
        <h1 className={authStyles.heading}>Sign in to Dnote</h1>

        <div className={authStyles.body}>
          {referrer && (
            <Flash type="info" wrapperClassName={authStyles['referrer-flash']}>
              Please sign in to continue to that page.
            </Flash>
          )}

          <div className={authStyles.panel}>
            {errMsg && (
              <Flash
                id="T-login-error"
                type="danger"
                wrapperClassName={authStyles['error-flash']}
              >
                {errMsg}
              </Flash>
            )}

            <LoginForm
              onLogin={handleLogin}
              submitting={submitting}
              onUpdateEmail={doUpdateAuthEmail}
              email={formData.auth.email}
            />
          </div>

          <div className={authStyles.footer}>
            <div className={authStyles.callout}>Don&#39;t have an account?</div>
            <Link to="/join" className={authStyles.cta}>
              Create account
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
function mapStateToProps(state) {
  return {
    formData: state.form
  };
}

const mapDispatchToProps = {
  doGetCurrentUser: getCurrentUser,
  doUpdateAuthEmail: updateAuthEmail
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Login);
