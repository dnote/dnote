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
import classnames from 'classnames';

import EmailModal from './EmailModal';
import PasswordModal from './PasswordModal';
import SettingRow from '../SettingRow';
import EmailVerificationRow from './EmailVerificationRow';
import Flash from '../../Common/Flash';
import { useSelector } from '../../../store';

import styles from '../Settings.scss';

interface Props {}

const Account: React.SFC<Props> = () => {
  const [emailModalOpen, setEmailModalOpen] = useState(false);
  const [passwordModalOpen, setPasswordModalOpen] = useState(false);
  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');

  const { user } = useSelector(state => {
    return {
      user: state.auth.user.data
    };
  });

  return (
    <div id="T-account-page">
      <Helmet>
        <title>Account</title>
      </Helmet>

      <h1 className="sr-only">Account settings</h1>

      <Flash
        when={successMsg !== ''}
        id="T-success-flash"
        kind="success"
        wrapperClassName={styles.flash}
        onDismiss={() => {
          setSuccessMsg('');
        }}
      >
        {successMsg}
      </Flash>
      <Flash
        when={failureMsg !== ''}
        kind="danger"
        wrapperClassName={styles.flash}
        onDismiss={() => {
          setFailureMsg('');
        }}
      >
        {failureMsg}
      </Flash>

      <div className={styles.wrapper}>
        <section className={styles.section}>
          <h2 className={styles['section-heading']}>Profile</h2>

          <SettingRow
            name="Email"
            value={user.email}
            actionContent={
              <button
                id="T-change-email-button"
                className={classnames('button-no-ui', styles.edit)}
                type="button"
                onClick={() => {
                  setEmailModalOpen(true);
                }}
              >
                Edit
              </button>
            }
          />

          <EmailVerificationRow
            verified={user.emailVerified}
            setSuccessMsg={setSuccessMsg}
            setFailureMsg={setFailureMsg}
          />

          <SettingRow
            name="Password"
            desc=" Set a unique password to protect your data."
            actionContent={
              <button
                id="T-change-password-button"
                className={classnames('button-no-ui', styles.edit)}
                type="button"
                onClick={() => {
                  setPasswordModalOpen(true);
                }}
              >
                Edit
              </button>
            }
          />
        </section>
      </div>

      <EmailModal
        isOpen={emailModalOpen}
        onDismiss={() => {
          setEmailModalOpen(false);
        }}
        currentEmail={user.email}
      />

      <PasswordModal
        isOpen={passwordModalOpen}
        onDismiss={() => {
          setPasswordModalOpen(false);
        }}
      />
    </div>
  );
};

export default Account;
