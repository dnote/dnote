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
import { Redirect } from 'react-router-dom';

import { ClassicMigrationSteps, getClassicMigrationPath } from 'web/libs/paths';
import services from 'web/libs/services';
import { bufToB64, b64ToBuf } from 'web/libs/encoding';
import { useSelector, useDispatch } from '../../store';
import { loginHelper, aes256GcmDecrypt } from '../../crypto';
import { getCurrentUser } from '../../store/auth';
import authStyles from '../Common/Auth.scss';
import Logo from '../Icons/Logo';
import Flash from '../Common/Flash';
import LoginForm from '../Login/LoginForm';

interface Props {}

const ClassicLogin: React.SFC<Props> = () => {
  const [errMsg, setErrMsg] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [emailVal, setEmailVal] = useState('');
  const dispatch = useDispatch();

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
      const { iteration } = await services.users.classicPresignin({ email });

      if (iteration === 0) {
        throw new Error('Please login from /login');
      }

      const { masterKey, authKey } = await loginHelper({
        email,
        password,
        iteration
      });
      const signinResp = await services.users.classicSignin({ email, authKey });

      const cipherKey = await aes256GcmDecrypt(
        b64ToBuf(masterKey),
        b64ToBuf(signinResp.cipherKeyEnc)
      );
      localStorage.setItem('cipherKey', bufToB64(cipherKey));

      await dispatch(getCurrentUser());
    } catch (err) {
      console.log(err);
      setErrMsg(err.message);
      setSubmitting(false);
    }
  }

  const { user } = useSelector(state => {
    return {
      user: state.auth.user
    };
  });

  const userData = user.data;
  const loggedIn = userData.uuid !== '';

  if (loggedIn && userData.classic) {
    return (
      <Redirect
        to={getClassicMigrationPath(ClassicMigrationSteps.setPassword)}
      />
    );
  }

  return (
    <div className={authStyles.page}>
      <Helmet>
        <title>Sign In (Classic)</title>
      </Helmet>
      <div className="container">
        <a href="/">
          <Logo fill="#252833" width={60} height={60} />
        </a>
        <h1 className={authStyles.heading}>Sign in to Dnote Classic</h1>
        <div className={authStyles.body}>
          <div className={authStyles.panel}>
            {errMsg && (
              <Flash
                id="T-login-error"
                kind="danger"
                wrapperClassName={authStyles['error-flash']}
              >
                {errMsg}
              </Flash>
            )}

            <LoginForm
              onLogin={handleLogin}
              submitting={submitting}
              email={emailVal}
              onUpdateEmail={val => setEmailVal(val)}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default ClassicLogin;
