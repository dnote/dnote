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

import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';

import services from 'web/libs/services';
import { getCurrentUser } from '../../../store/auth';
import { useDispatch } from '../../../store';
import Button from '../../Common/Button';
import Modal, { Header, Body } from '../../Common/Modal';
import Flash from '../../Common/Flash';
import modalStyles from '../../Common/Modal/Modal.scss';
import settingsStyles from '../Settings.scss';

interface Props {
  currentEmail: string;
  isOpen: boolean;
  onDismiss: () => void;
}

const EmailModal: React.SFC<Props> = ({ currentEmail, isOpen, onDismiss }) => {
  const [passwordVal, setPasswordVal] = useState('');
  const [emailVal, setEmailVal] = useState('');
  const [inProgress, setInProgress] = useState(false);
  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');
  const dispatch = useDispatch();

  useEffect(() => {
    if (!isOpen) {
      setPasswordVal('');
      setEmailVal('');
    }
  }, [isOpen]);

  async function handleSubmit(e) {
    e.preventDefault();

    setSuccessMsg('');
    setFailureMsg('');
    setInProgress(true);

    try {
      if (currentEmail === emailVal) {
        throw new Error('The new email is the same as the old email');
      }

      await services.users.updateProfile({
        email: emailVal
      });

      await dispatch(getCurrentUser({ refresh: true }));

      setSuccessMsg('Updated the email');
      setEmailVal('');
      setPasswordVal('');
      setInProgress(false);
    } catch (err) {
      setFailureMsg(`Failed to update the email. Error: ${err.message}`);
      setInProgress(false);
    }
  }

  const labelId = 'email-modal';

  return (
    <Modal
      modalId="T-change-email-modal"
      isOpen={isOpen}
      onDismiss={onDismiss}
      ariaLabelledBy={labelId}
    >
      <Header labelId={labelId} heading="Change email" onDismiss={onDismiss} />

      <Flash
        when={successMsg !== ''}
        id="T-change-email-modal-success"
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

      <Body>
        <form onSubmit={handleSubmit} autoComplete="off">
          {/* prevent browsers from automatically filling the input fields */}
          <input type="password" style={{ display: 'none' }} readOnly />

          <div className={modalStyles['input-row']}>
            <label className="input-label" htmlFor="email-form-email-input">
              New email
            </label>

            <input
              id="email-form-email-input"
              type="email"
              placeholder="your@email.com"
              value={emailVal}
              onChange={e => {
                const val = e.target.value;
                setEmailVal(val);
              }}
              className="form-control"
              autoComplete="new-password"
            />
          </div>

          <div className={modalStyles['input-row']}>
            <label className="input-label" htmlFor="email-form-password-input">
              Current password
            </label>

            <input
              id="email-form-password-input"
              type="password"
              placeholder="********"
              value={passwordVal}
              onChange={e => {
                const val = e.target.value;
                setPasswordVal(val);
              }}
              className="form-control"
              autoComplete="off"
            />
          </div>

          <div className={modalStyles.actions}>
            <Button
              type="submit"
              kind="first"
              size="normal"
              isBusy={inProgress}
            >
              Update
            </Button>

            <Button
              type="button"
              kind="second"
              size="normal"
              isBusy={inProgress}
              onClick={onDismiss}
            >
              Cancel
            </Button>
          </div>
        </form>
      </Body>
    </Modal>
  );
};

const mapDispatchToProps = {
  doGetCurrentUser: getCurrentUser
};

export default connect(
  null,
  mapDispatchToProps
)(EmailModal);
