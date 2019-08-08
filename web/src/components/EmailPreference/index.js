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
import classnames from 'classnames';
import { Link, withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import Helmet from 'react-helmet';

import EmailPreferenceForm from '../Common/EmailPreferenceForm';
import Logo from '../Icons/Logo';
import Flash from '../Common/Flash';

import { parseSearchString } from '../../libs/url';
import { getEmailPreference } from '../../store/auth';
import { getLoginPath } from '../../libs/paths';

import styles from './EmailPreference.module.scss';

function EmailPreference({
  location,
  emailPreferenceData,
  doGetEmailPreference
}) {
  const emailPreference = emailPreferenceData.data;
  const { isFetched, isFetching, errorMessage } = emailPreferenceData;
  const { token } = parseSearchString(location.search);

  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');

  useState(() => {
    if (isFetched || isFetching) {
      return;
    }

    doGetEmailPreference(token);
  }, []);

  return (
    <div className={styles.wrapper}>
      <Helmet>
        <title>Email Preference</title>
      </Helmet>

      <Link to="/">
        <Logo fill="#252833" width="60" height="60" />
      </Link>
      <h1 className={styles.heading}>Dnote email settings</h1>

      <div className="container">
        <div className={styles.body}>
          {errorMessage && (
            <Flash type="danger" wrapperClassName={styles.flash}>
              Error fetching email preference: {errorMessage}.{' '}
              <span>
                Please <Link to={getLoginPath()}>login</Link> and try again.
              </span>
            </Flash>
          )}
          {successMsg && (
            <Flash
              type="success"
              wrapperClassName={classnames(styles.flash, 'T-success')}
            >
              {successMsg}
            </Flash>
          )}
          {failureMsg && (
            <Flash type="danger" wrapperClassName={styles.flash}>
              {failureMsg}
            </Flash>
          )}
          {isFetched && (
            <EmailPreferenceForm
              token={token}
              emailPreference={emailPreference}
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
}

function mapStateToProps(state) {
  return {
    emailPreferenceData: state.auth.emailPreference
  };
}

const mapDispatchToProps = {
  doGetEmailPreference: getEmailPreference
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(EmailPreference)
);
