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
import { connect } from 'react-redux';

import { updateEmailPreference } from '../../services/users';
import { receiveEmailPreference } from '../../store/auth';
import Button from './Button';

import styles from './EmailPreferenceForm.module.scss';

const digestWeekly = 'weekly';
const digestNever = 'never';

function getDigestFrequency(emailPreference) {
  if (emailPreference.digest_weekly) {
    return digestWeekly;
  }

  return digestNever;
}

function EmailPreferenceForm({
  emailPreference,
  doReceiveEmailPreference,
  token,
  setSuccessMsg,
  setFailureMsg,
  actionsClassName
}) {
  const freq = getDigestFrequency(emailPreference);
  const [digestFrequency, setDigestFrequency] = useState(freq);
  const [inProgress, setInProgress] = useState(false);

  function handleSubmit(e) {
    e.preventDefault();

    setSuccessMsg('');
    setFailureMsg('');
    setInProgress(true);

    updateEmailPreference({ digestFrequency, token })
      .then(updatedPreference => {
        doReceiveEmailPreference(updatedPreference);

        setSuccessMsg('Updated email preference');
        setInProgress(false);
      })
      .catch(err => {
        setFailureMsg(`Failed to update. Error: ${err.message}`);
        setInProgress(false);
      });
  }

  return (
    <div>
      <div>
        <form id="T-email-pref-form" onSubmit={handleSubmit}>
          <div className={styles.heading}>Email digest frequency</div>

          <div>
            <div className={styles.radio}>
              <label htmlFor="digest-never">
                <input
                  id="digest-never"
                  type="radio"
                  name="digest"
                  value={digestNever}
                  checked={digestFrequency === digestNever}
                  onChange={e => {
                    const val = e.target.value;
                    setDigestFrequency(val);
                  }}
                />
                Never
              </label>
            </div>

            <div className={styles.radio}>
              <label htmlFor="digest-weekly">
                <input
                  id="digest-weekly"
                  type="radio"
                  name="digest"
                  value={digestWeekly}
                  checked={digestFrequency === digestWeekly}
                  onChange={e => {
                    const val = e.target.value;
                    setDigestFrequency(val);
                  }}
                />
                Weekly (Friday)
              </label>
            </div>
          </div>

          <div className={actionsClassName}>
            <Button type="submit" kind="first" isBusy={inProgress}>
              Update
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}

const mapDispatchToProps = {
  doReceiveEmailPreference: receiveEmailPreference
};

export default connect(
  null,
  mapDispatchToProps
)(EmailPreferenceForm);
