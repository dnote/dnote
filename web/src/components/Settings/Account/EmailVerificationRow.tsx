import React, { useState } from 'react';
import classnames from 'classnames';

import SettingRow from '../SettingRow';
import * as usersService from '../../../services/users';

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
};

export default EmailVerificationRow;
