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

import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import classnames from 'classnames';
import Helmet from 'react-helmet';
import { connect } from 'react-redux';

import Header from '../../Common/Page/Header';
import Body from '../../Common/Page/Body';
import Flash from '../../Common/Flash';
import { getEmailPreference } from '../../../actions/auth';
import FrequencyModal from './FrequencyModal';
import SettingRow from '../SettingRow';
import { getSettingsPath } from '../../../libs/paths';

import settingsStyles from '../Settings.module.scss';

function getFrequencyLabel(emailPreference) {
  if (emailPreference.digest_weekly) {
    return 'Weekly';
  }

  return 'Never';
}

function Email({ emailPreferenceData, doGetEmailPreference }) {
  useEffect(() => {
    if (!emailPreferenceData.isFetched) {
      doGetEmailPreference();
    }
  }, [doGetEmailPreference, emailPreferenceData.isFetched]);

  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');
  const [isFrequencyModalOpen, setIsFrequencyModalOpen] = useState(false);

  return (
    <div>
      <Helmet>
        <title>Notification</title>
      </Helmet>

      <Header heading="Notification" />

      <Body>
        <div className="container-wide">
          {successMsg && (
            <div className="row">
              <div className="col-12 col-lg-10">
                <Flash
                  type="success"
                  wrapperClassName={settingsStyles.flash}
                  onDismiss={() => {
                    setSuccessMsg('');
                  }}
                >
                  {successMsg}
                </Flash>
              </div>
            </div>
          )}

          {failureMsg && (
            <div className="row">
              <div className="col-12 col-lg-10">
                <Flash
                  type="danger"
                  wrapperClassName={settingsStyles.flash}
                  onDismiss={() => {
                    setFailureMsg('');
                  }}
                >
                  {failureMsg}
                </Flash>
              </div>
            </div>
          )}

          {emailPreferenceData.errorMessage && (
            <div className="row">
              <div className="col-12 col-lg-10">
                <Flash type="danger" wrapperClassName={settingsStyles.flash}>
                  <div>Error fetching notification preference:</div>
                  {emailPreferenceData.errorMessage}
                </Flash>
              </div>
            </div>
          )}

          <div className="row">
            <div className="col-12 col-lg-10">
              <Flash
                type="info"
                wrapperClassName={settingsStyles.flash}
                contentClassName={settingsStyles['verification-banner']}
              >
                <div>
                  You need to verify your email before Dnote can send you
                  digests.
                </div>
                <Link
                  to={getSettingsPath('account')}
                  className={classnames(
                    'button button-normal button-second',
                    settingsStyles['verification-banner-cta']
                  )}
                >
                  Go to account settings
                </Link>
              </Flash>
            </div>
          </div>

          <div className="row">
            <div className="col-12 col-lg-10">
              <section className={settingsStyles.section}>
                <h2 className={settingsStyles['section-heading']}>Digest</h2>

                <SettingRow
                  name="Frequency"
                  value={getFrequencyLabel(emailPreferenceData.data)}
                  actionContent={
                    <button
                      id="T-edit-frequency-button"
                      className={classnames(
                        'button-no-ui',
                        settingsStyles.edit
                      )}
                      type="button"
                      onClick={() => {
                        setIsFrequencyModalOpen(true);
                      }}
                      disabled={!emailPreferenceData.isFetched}
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

      <FrequencyModal
        emailPreferenceData={emailPreferenceData}
        isOpen={isFrequencyModalOpen}
        onDismiss={() => {
          setIsFrequencyModalOpen(false);
        }}
      />
    </div>
  );
}

function mapStateToProps(state) {
  return {
    userData: state.auth.user,
    emailPreferenceData: state.auth.emailPreference
  };
}

const mapDispatchToProps = {
  doGetEmailPreference: getEmailPreference
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Email);
