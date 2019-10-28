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

import { RepetitionRuleData } from 'jslib/operations/types';
import services from 'web/libs/services';
import Button from '../../Common/Button';
import styles from './EmailPreferenceRepetition.scss';

interface Props {
  data: RepetitionRuleData;
  setSuccessMsg: (string) => void;
  setFailureMsg: (string) => void;
  token?: string;
}

const Content: React.SFC<Props> = ({
  data,
  token,
  setSuccessMsg,
  setFailureMsg
}) => {
  const [inProgress, setInProgress] = useState(false);
  const [isEnabled, setIsEnabled] = useState(data.enabled);

  function handleSubmit(e) {
    e.preventDefault();

    setSuccessMsg('');
    setFailureMsg('');
    setInProgress(true);

    services.repetitionRules
      .update(data.uuid, { enabled: isEnabled }, { token })
      .then(() => {
        setSuccessMsg('Updated the repetition.');
        setInProgress(false);
      })
      .catch(err => {
        setFailureMsg(`Failed to update. Error: ${err.message}`);
        setInProgress(false);
      });
  }

  return (
    <div>
      <p>Toggle the repetition for "{data.title}"</p>

      <form id="T-pref-repetition-form" onSubmit={handleSubmit}>
        <div>
          <div className={styles.radio}>
            <label htmlFor="repetition-off">
              <input
                id="repetition-off"
                type="radio"
                name="repetition"
                value="off"
                checked={!isEnabled}
                onChange={e => {
                  const val = e.target.value;
                  setIsEnabled(false);
                }}
              />
              Disable
            </label>
          </div>

          <div className={styles.radio}>
            <label htmlFor="repetition-on">
              <input
                id="repetition-on"
                type="radio"
                name="repetition"
                value="on"
                checked={isEnabled}
                onChange={e => {
                  const val = e.target.value;
                  setIsEnabled(true);
                }}
              />
              Enable
            </label>
          </div>
        </div>

        <div>
          <Button type="submit" kind="first" size="normal" isBusy={inProgress}>
            Update
          </Button>
        </div>
      </form>
    </div>
  );
};

export default Content;
