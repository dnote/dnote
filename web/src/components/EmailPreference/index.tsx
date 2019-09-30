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

import React, { useState, useEffect } from 'react';
import classnames from 'classnames';
import { Link, withRouter, RouteComponentProps } from 'react-router-dom';
import Helmet from 'react-helmet';

import EmailPreferenceForm from '../Common/EmailPreferenceForm';
import Logo from '../Icons/Logo';
import Flash from '../Common/Flash';
import { parseSearchString } from 'jslib/helpers/url';
import { getEmailPreference } from '../../store/auth';
import { getLoginPath } from 'web/libs/paths';
import { useSelector, useDispatch } from '../../store';
import styles from './EmailPreference.scss';

interface Props extends RouteComponentProps {}

const EmailPreference: React.SFC<Props> = ({ location }) => {
  const { emailPreference } = useSelector(state => {
    return {
      emailPreference: state.auth.emailPreference
    };
  });
  const dispatch = useDispatch();

  const emailPreferenceData = emailPreference.data;
  const { isFetched, isFetching, errorMessage } = emailPreference;
  const { token } = parseSearchString(location.search);

  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');

  useEffect(() => {
    if (isFetched || isFetching) {
      return;
    }

    dispatch(getEmailPreference(token));
  }, [dispatch, isFetched, isFetching, token]);

  return (
    <div className={styles.wrapper}>
      <Helmet>
        <title>Email Preference</title>
      </Helmet>

      <Link to="/">
        <Logo fill="#252833" width={60} height={60} />
      </Link>
      <h1 className={styles.heading}>Dnote email settings</h1>

      <div className="container">
        <div className={styles.body}>
          <Flash
            when={errorMessage !== ''}
            kind="danger"
            wrapperClassName={styles.flash}
          >
            Error fetching email preference: {errorMessage}.{' '}
            <span>
              Please <Link to={getLoginPath()}>login</Link> and try again.
            </span>
          </Flash>

          <Flash
            when={successMsg !== ''}
            kind="success"
            wrapperClassName={classnames(styles.flash, 'T-success')}
          >
            {successMsg}
          </Flash>
          <Flash
            when={failureMsg !== ''}
            kind="danger"
            wrapperClassName={styles.flash}
          >
            {failureMsg}
          </Flash>

          {isFetched && (
            <EmailPreferenceForm
              token={token}
              emailPreference={emailPreferenceData}
              setSuccessMsg={setSuccessMsg}
              setFailureMsg={setFailureMsg}
            />
          )}
        </div>
        <div className={styles.footer}>
          <Link to="/">Back to Dnote home</Link>
        </div>
      </div>
    </div>
  );
};

export default withRouter(EmailPreference);
