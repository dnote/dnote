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

import Modal, { Header, Body } from '../../Common/Modal';
import Flash from '../../Common/Flash';
import EmailPreferenceForm from '../../Common/EmailPreferenceForm';

import settingsStyles from '../Settings.scss';

interface Props {
  emailPreference: any;
  isOpen: boolean;
  onDismiss: () => void;
}

const FrequencyModal: React.SFC<Props> = ({
  emailPreference,
  isOpen,
  onDismiss
}) => {
  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');

  const labelId = 'frequency-modal';

  return (
    <Modal
      modalId="T-frequency-modal"
      isOpen={isOpen}
      onDismiss={onDismiss}
      ariaLabelledBy={labelId}
    >
      <Header
        labelId={labelId}
        heading="Change diegst frequency"
        onDismiss={onDismiss}
      />

      <Flash
        when={successMsg !== ''}
        id="T-frequency-modal-success"
        kind="success"
        wrapperClassName={settingsStyles.flash}
        onDismiss={() => {
          setSuccessMsg('');
        }}
      >
        {successMsg}
      </Flash>
      <Flash
        when={failureMsg !== ''}
        kind="danger"
        wrapperClassName={settingsStyles.flash}
        onDismiss={() => {
          setFailureMsg('');
        }}
      >
        {failureMsg}
      </Flash>

      <Body>
        {emailPreference.isFetched ? (
          <EmailPreferenceForm
            emailPreference={emailPreference.data}
            actionsClassName={settingsStyles.actions}
            setSuccessMsg={setSuccessMsg}
            setFailureMsg={setFailureMsg}
          />
        ) : (
          <div>Fetching email preference...</div>
        )}
      </Body>
    </Modal>
  );
};

export default FrequencyModal;
