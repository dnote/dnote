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

import Flash from '../../Common/Flash';
import { getEmailPreference } from '../../../store/auth';
import { useSelector, useDispatch } from '../../../store';
import FrequencyModal from './FrequencyModal';
import SettingRow from '../SettingRow';
import { SettingSections, getSettingsPath } from 'web/libs/paths';
import styles from '../Settings.scss';

function getFrequencyLabel(emailPreference) {
  if (emailPreference.digest_weekly) {
    return 'Weekly';
  }

  return 'Never';
}

interface Props {}

const Email: React.SFC<Props> = () => {
  const dispatch = useDispatch();

  const { user, emailPreference } = useSelector(state => {
    return {
      user: state.auth.user.data,
      emailPreference: state.auth.emailPreference
    };
  });

  useEffect(() => {
    if (!emailPreference.isFetched) {
      dispatch(getEmailPreference());
    }
  }, [dispatch, emailPreference.isFetched]);

  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');
  const [isFrequencyModalOpen, setIsFrequencyModalOpen] = useState(false);

  return (
    <div>
      <Helmet>
        <title>Notification</title>
      </Helmet>

      <Flash
        when={successMsg !== ''}
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

      <Flash
        when={emailPreference.errorMessage !== ''}
        kind="danger"
        wrapperClassName={styles.flash}
      >
        <div>Error fetching notification preference:</div>
        {emailPreference.errorMessage}
      </Flash>

      <Flash
        when={!user.emailVerified}
        kind="info"
        wrapperClassName={styles.flash}
        contentClassName={styles['verification-banner']}
      >
        <div>
          You need to verify your email before Dnote can send you digests.
        </div>
        <Link
          to={getSettingsPath(SettingSections.account)}
          className={classnames(
            'button button-normal button-second',
            styles['verification-banner-cta']
          )}
        >
          Go to account settings
        </Link>
      </Flash>

      <div className={styles.wrapper}>
        <section className={styles.section}>
          <h2 className={styles['section-heading']}>Email digest</h2>

          <SettingRow
            name="Frequency"
            value={getFrequencyLabel(emailPreference.data)}
            actionContent={
              <button
                id="T-edit-frequency-button"
                className={classnames('button-no-ui', styles.edit)}
                type="button"
                onClick={() => {
                  setIsFrequencyModalOpen(true);
                }}
                disabled={!emailPreference.isFetched}
              >
                Edit
              </button>
            }
          />
        </section>
      </div>

      <FrequencyModal
        emailPreference={emailPreference}
        isOpen={isFrequencyModalOpen}
        onDismiss={() => {
          setIsFrequencyModalOpen(false);
        }}
      />
    </div>
  );
};

export default Email;
