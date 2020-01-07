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
import Helmet from 'react-helmet';
import { withRouter, RouteComponentProps } from 'react-router';

import { parseSearchString } from 'jslib/helpers/url';
import Flash from '../../Common/Flash';
import { useSelector } from '../../../store';
import Form from './Form';
import settingsStyles from '../Settings.scss';
import { useDispatch } from '../../../store';
import { getEmailPreference } from '../../../store/auth';
import styles from './Notifications.scss';

interface Props extends RouteComponentProps {}

const Notifications: React.FunctionComponent<Props> = ({ location }) => {
  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');
  const dispatch = useDispatch();

  const { token } = parseSearchString(location.search);

  const { emailPreference } = useSelector(state => {
    return {
      emailPreference: state.auth.emailPreference
    };
  });

  useEffect(() => {
    if (!emailPreference.isFetched) {
      dispatch(getEmailPreference());
    }
  }, [dispatch, emailPreference.isFetched]);

  return (
    <div>
      <Helmet>
        <title>Notifications</title>
      </Helmet>

      <h1 className="sr-only">Notifications</h1>

      <Flash
        id="T-notifications-success"
        when={successMsg !== ''}
        kind="success"
        onDismiss={() => {
          setSuccessMsg('');
        }}
      >
        {successMsg}
      </Flash>

      <Flash
        when={failureMsg !== ''}
        kind="danger"
        onDismiss={() => {
          setFailureMsg('');
        }}
      >
        {failureMsg}
      </Flash>

      <Flash when={emailPreference.errorMessage !== ''} kind="danger">
        {emailPreference.errorMessage}
      </Flash>

      <div className={settingsStyles.wrapper}>
        <section className={settingsStyles.section}>
          <h2 className={settingsStyles['section-heading']}>
            Email Preferences
          </h2>

          <div className={styles.body}>
            {emailPreference.isFetched ? (
              <Form
                emailPreference={emailPreference.data}
                setSuccessMsg={setSuccessMsg}
                setFailureMsg={setFailureMsg}
                token={token}
              />
            ) : (
              <span>Loading email preferences...</span>
            )}
          </div>
        </section>
      </div>
    </div>
  );
};

export default withRouter(Notifications);
