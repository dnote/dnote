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

import React, { useState, useReducer } from 'react';
import classnames from 'classnames';

import Modal, { Header, Body } from '../Common/Modal';
import Flash from '../Common/Flash';
import { daysToSec } from '../../helpers/time';
import Button from '../Common/Button';
import GlobeIcon from '../Icons/Globe';
import styles from './CreateModal.scss';
import Form, { FormState } from './Form';
import modalStyles from '../Common/Modal/Modal.scss';

interface Props {
  isOpen: boolean;
  onDismiss: () => void;
}

const CreateRuleModal: React.FunctionComponent<Props> = ({
  isOpen,
  onDismiss
}) => {
  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');
  const [inProgress, setInProgress] = useState(false);

  const labelId = 'create-repetition-modal';

  function handleSubmit(state: FormState) {}

  return (
    <Modal
      modalId="T-create-rule-modal"
      modalClassName={styles.content}
      isOpen={isOpen}
      onDismiss={onDismiss}
      ariaLabelledBy={labelId}
    >
      <Header
        labelId={labelId}
        heading="Create a spaced repetition"
        onDismiss={onDismiss}
      />

      <Flash
        when={successMsg !== ''}
        id="T-create-rule-modal-success"
        kind="success"
        onDismiss={() => {
          setSuccessMsg('');
        }}
        noMargin
      >
        {successMsg}
      </Flash>
      <Flash
        when={failureMsg !== ''}
        kind="danger"
        onDismiss={() => {
          setFailureMsg('');
        }}
        noMargin
      >
        {failureMsg}
      </Flash>

      <Body></Body>
    </Modal>
  );
};

export default CreateRuleModal;
