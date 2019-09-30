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

import services from 'web/libs/services';
import SettingRow from '../SettingRow';
import settingsStyles from '../Settings.scss';

interface Props {
  verified: boolean;
  setSuccessMsg: (string) => void;
  setFailureMsg: (string) => void;
}

const EmailVerificationRow: React.SFC<Props> = ({
  verified,
  setSuccessMsg,
  setFailureMsg
}) => {
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

          services.users
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
};

export default EmailVerificationRow;
