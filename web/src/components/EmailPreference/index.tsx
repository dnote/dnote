/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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
import { RouteComponentProps } from 'react-router';
import { Link, withRouter } from 'react-router-dom';
import Helmet from 'react-helmet';

import { parseSearchString } from 'jslib/helpers/url';
import { getHomePath, getLoginPath } from 'web/libs/paths';
import Form from '../Settings/Notifications/Form';
import { getEmailPreference } from '../../store/auth';
import Logo from '../Icons/Logo';
import { useSelector, useDispatch } from '../../store';
import Flash from '../Common/Flash';
import styles from './EmailPreference.scss';

interface Props extends RouteComponentProps {}

const EmailPreference: React.FunctionComponent<Props> = ({ location }) => {
  const { token } = parseSearchString(location.search);

  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');

  const dispatch = useDispatch();

  const { emailPreference } = useSelector(state => {
    return {
      emailPreference: state.auth.emailPreference
    };
  });

  useEffect(() => {
    if (!emailPreference.isFetched) {
      dispatch(getEmailPreference(token));
    }
  }, [dispatch, emailPreference.isFetched, token]);

  return (
    <div className={styles.wrapper}>
      <Helmet>
        <title>Email Preferences</title>
      </Helmet>

      <Link to="/">
        <Logo fill="#252833" width={60} height={60} />
      </Link>
      <h1 className={styles.heading}>Dnote email preferences</h1>

      <div className="container">
        <div className={styles.body}>
          <Flash
            kind="danger"
            wrapperClassName={styles.flash}
            when={emailPreference.errorMessage !== ''}
          >
            Error fetching email preference: {emailPreference.errorMessage}.
            Please try again after <Link to={getLoginPath()}>logging in</Link>.
          </Flash>

          <Flash
            id="T-notifications-success"
            kind="success"
            wrapperClassName={classnames(styles.flash, 'T-success')}
            when={successMsg !== ''}
          >
            {successMsg}
          </Flash>

          <Flash
            kind="danger"
            wrapperClassName={styles.flash}
            when={failureMsg !== ''}
          >
            {failureMsg}
          </Flash>

          {emailPreference.isFetched && (
            <Form
              token={token}
              emailPreference={emailPreference.data}
              setSuccessMsg={setSuccessMsg}
              setFailureMsg={setFailureMsg}
            />
          )}
        </div>
        <div className={styles.footer}>
          <Link to={getHomePath()}>Back to Dnote home</Link>
        </div>
      </div>
    </div>
  );
};

export default withRouter(EmailPreference);
