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
import { connect } from 'react-redux';

import JoinForm from './JoinForm';
import Logo from '../Icons/Logo';
import Flash from '../Common/Flash';
import { getReferrer } from '../../libs/url';
import { FormState, updateAuthEmail } from '../../store/form';
import { AppState } from '../../store';
import * as usersService from '../../services/users';
import { getCurrentUser } from '../../store/auth';
import authStyles from '../Common/Auth.scss';

interface Props extends RouteComponentProps<any> {
  doGetCurrentUser: () => Promise<any>;
  doUpdateAuthEmail: (string) => null;
  formData: FormState;
}

const Join: React.SFC<Props> = ({
  doGetCurrentUser,
  formData,
  doUpdateAuthEmail,
  location
}) => {
  const [errMsg, setErrMsg] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const referrer = getReferrer(location);

  async function handleJoin(email, password, passwordConfirmation) {
    if (!email) {
      setErrMsg('Please enter an email.');
      return;
    }
    if (!password) {
      setErrMsg('Please enter a pasasword.');
      return;
    }
    if (!passwordConfirmation) {
      setErrMsg('The passwords do not match.');
      return;
    }

    setErrMsg('');
    setSubmitting(true);

    try {
      await usersService.register({ email, password });

      // guestOnly HOC will redirect the user accordingly after the current user is fetched
      await doGetCurrentUser();
      doUpdateAuthEmail('');
    } catch (err) {
      console.log(err);
      setErrMsg(err.message);
      setSubmitting(false);
    }
  }

  return (
    <div className={authStyles.page}>
      <Helmet>
        <title>Join</title>
      </Helmet>
      <div className="container">
        <a href="/">
          <Logo fill="#252833" width={60} height={60} />
        </a>
        <h1 className={authStyles.heading}>Join Dnote</h1>

        <div className={authStyles.body}>
          {referrer && (
            <Flash kind="info" wrapperClassName={authStyles['referrer-flash']}>
              Please join to continue.
            </Flash>
          )}

          <div className={authStyles.panel}>
            {errMsg && (
              <Flash kind="danger" wrapperClassName={authStyles['error-flash']}>
                {errMsg}
              </Flash>
            )}

            <JoinForm
              onJoin={handleJoin}
              submitting={submitting}
              email={formData.auth.email}
              onUpdateEmail={doUpdateAuthEmail}
            />
          </div>

          <div className={authStyles.footer}>
            <div className={authStyles.callout}>Already have an account?</div>
            <Link to="/login" className={authStyles.cta}>
              Sign in
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

function mapStateToProps(state: AppState) {
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
)(Join);
