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
import { connect } from 'react-redux';
import classnames from 'classnames';

import Header from '../../Common/Page/Header';
import Body from '../../Common/Page/Body';
import EmailModal from './EmailModal';
import PasswordModal from './PasswordModal';
import SettingRow from '../SettingRow';
import Flash from '../../Common/Flash';
import * as usersService from '../../../services/users';

import settingsStyles from '../Settings.module.scss';

function EmailVerificationRow({ verified, setSuccessMsg, setFailureMsg }) {
  const [isSent, setIsSent] = useState(false);
  const [inProgress, setInProgress] = useState(false);

  let value;
  if (verified) {
    value = 'Yes';
  } else {
    value = 'No';
  }

  let actionContent = null;
  if (!verified) {
    actionContent = (
      <button
        className={classnames('button-no-ui', settingsStyles.edit)}
        type="button"
        id="T-send-verification-button"
        onClick={() => {
          setInProgress(true);

          usersService
            .sendEmailVerificationEmail()
            .then(() => {
              setIsSent(true);
              setInProgress(false);
              setSuccessMsg(
                'A verification email is on the way. It might take a minute. Pleae check your inbox.'
              );
            })
            .catch(err => {
              setInProgress(false);
              setFailureMsg(
                `Failed to send a verification email. Error: ${err.message}`
              );
            });
        }}
        disabled={inProgress || isSent}
      >
        {isSent ? 'Verification sent' : 'Send verification'}
      </button>
    );
  }

  let desc = '';
  if (!verified) {
    desc = 'You need to verify the email to receive email digests';
  }

  return (
    <SettingRow
      name="Email verified"
      desc={desc}
      value={value}
      actionContent={actionContent}
    />
  );
}

function Account({ userState }) {
  const [emailModalOpen, setEmailModalOpen] = useState(false);
  const [passwordModalOpen, setPasswordModalOpen] = useState(false);
  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');

  const user = userState.data;

  return (
    <div id="T-account-page">
      <Helmet>
        <title>Account</title>
      </Helmet>

      <Header heading="Account" />
      <Body>
        <div className="container-wide">
          <div className="row">
            <div className="col-12 col-md-12 col-lg-10">
              {successMsg && (
                <Flash
                  id="T-success-flash"
                  type="success"
                  wrapperClassName={settingsStyles.flash}
                  onDismiss={() => {
                    setSuccessMsg('');
                  }}
                >
                  {successMsg}
                </Flash>
              )}
              {failureMsg && (
                <Flash
                  type="danger"
                  wrapperClassName={settingsStyles.flash}
                  onDismiss={() => {
                    setFailureMsg('');
                  }}
                >
                  {failureMsg}
                </Flash>
              )}
            </div>

            <div className="col-12 col-md-12 col-lg-10">
              <section className={settingsStyles.section}>
                <h2 className={settingsStyles['section-heading']}>Profile</h2>

                <SettingRow
                  name="Email"
                  value={user.email}
                  actionContent={
                    <button
                      id="T-change-email-button"
                      className={classnames(
                        'button-no-ui',
                        settingsStyles.edit
                      )}
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
                  verified={user.email_verified}
                  email={user.email}
                  setSuccessMsg={setSuccessMsg}
                  setFailureMsg={setFailureMsg}
                />

                <SettingRow
                  name="Password"
                  desc=" Set a unique password to protect your data."
                  actionContent={
                    <button
                      id="T-change-password-button"
                      className={classnames(
                        'button-no-ui',
                        settingsStyles.edit
                      )}
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
          </div>
        </div>
      </Body>

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
        email={user.email}
      />
    </div>
  );
}

function mapStateToProps(state) {
  return {
    userState: state.auth.user
  };
}

export default connect(mapStateToProps)(Account);
