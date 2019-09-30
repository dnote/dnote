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
import { Link, RouteComponentProps } from 'react-router-dom';

import { getReferrer } from 'jslib/helpers/url';
import { getRootUrl } from 'web/libs/paths';
import services from 'web/libs/services';
import LoginForm from './LoginForm';
import Logo from '../Icons/Logo';
import Flash from '../Common/Flash';
import { getCurrentUser } from '../../store/auth';
import { updateAuthEmail } from '../../store/form';
import authStyles from '../Common/Auth.scss';
import { useSelector, useDispatch } from '../../store';

interface Props extends RouteComponentProps<any> {}

const Login: React.SFC<Props> = ({ location }) => {
  const [errMsg, setErrMsg] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const dispatch = useDispatch();
  const { formData } = useSelector(state => {
    return {
      formData: state.form
    };
  });

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
      await services.users.signin({ email, password });

      // guestOnly HOC will redirect the user accordingly after the current user is fetched
      await dispatch(getCurrentUser());
      dispatch(updateAuthEmail(''));
    } catch (err) {
      console.log(err);
      setErrMsg(err.message);
      setSubmitting(false);
    }
  }
  const referrer = getReferrer(location);

  return (
    <div id="T-login-page" className={authStyles.page}>
      <Helmet>
        <title>Sign In</title>
      </Helmet>
      <div className="container">
        <a href={getRootUrl()}>
          <Logo fill="#252833" width={60} height={60} />
        </a>
        <h1 className={authStyles.heading}>Sign in to Dnote</h1>

        <div className={authStyles.body}>
          {referrer && (
            <Flash kind="info" wrapperClassName={authStyles['referrer-flash']}>
              Please sign in to continue to that page.
            </Flash>
          )}

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
              email={formData.auth.email}
              onUpdateEmail={val => {
                dispatch(updateAuthEmail(val));
              }}
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
};

export default Login;
