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

import { updatePassword, presignin } from '../../../services/users';
import { loginHelper, aes256GcmEncrypt } from '../../../crypto';
import { b64ToBuf, bufToB64 } from '../../../libs/encoding';
import Button from '../../Common/Button';
import Modal, { Header, Body } from '../../Common/Modal';
import Flash from '../../Common/Flash';

import settingsStyles from '../Settings.module.scss';

function PasswordModal({ email, isOpen, onDismiss }) {
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [newPasswordConfirmation, setNewPasswordConfirmation] = useState('');
  const [inProgress, setInProgress] = useState(false);
  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');

  useEffect(() => {
    if (!isOpen) {
      setOldPassword('');
      setNewPassword('');
      setNewPasswordConfirmation('');
    }
  }, [isOpen]);

  async function handleSubmit(e) {
    e.preventDefault();

    setSuccessMsg('');
    setFailureMsg('');
    setInProgress(true);

    try {
      if (newPassword !== newPasswordConfirmation) {
        throw new Error('Password and its confirmation do not match');
      }

      const { iteration } = await presignin({ email, oldPassword });
      const { authKey: oldAuthKey } = await loginHelper({
        email,
        password: oldPassword,
        iteration
      });
      const {
        masterKey: newMasterKey,
        authKey: newAuthKey
      } = await loginHelper({ email, password: newPassword, iteration });

      const cipherKey = localStorage.getItem('cipherKey');
      const cipherKeyBuf = b64ToBuf(cipherKey);

      const newCipherKeyEncBuf = await aes256GcmEncrypt(
        b64ToBuf(newMasterKey),
        cipherKeyBuf
      );
      const newCipherKeyEnc = bufToB64(newCipherKeyEncBuf);

      await updatePassword({
        oldAuthKey,
        newAuthKey,
        newCipherKeyEnc,
        newKdfIteration: iteration
      });

      setSuccessMsg('Updated the password');
      setOldPassword('');
      setNewPassword('');
      setNewPasswordConfirmation('');
      setInProgress(true);
    } catch (err) {
      setFailureMsg(`Failed to update. ${err.message}`);
      setInProgress(false);
    }
  }

  const labelId = 'password-modal';

  return (
    <Modal modalId="T-password-modal" isOpen={isOpen} onDismiss={onDismiss}>
      <Header
        labelId={labelId}
        heading="Change password"
        onDismiss={onDismiss}
        ariaLabelledBy={labelId}
      />

      {successMsg && (
        <Flash
          id="T-password-modal-success"
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

      <Body>
        <form onSubmit={handleSubmit}>
          {/* prevent browsers from automatically filling the input fields */}
          <input type="password" style={{ display: 'none' }} readOnly />

          <div className={settingsStyles['input-row']}>
            <label
              className={settingsStyles['input-label']}
              htmlFor="old-password-input"
            >
              Old password
              <input
                id="old-password-input"
                type="password"
                placeholder="********"
                value={oldPassword}
                onChange={e => {
                  const val = e.target.value;

                  setOldPassword(val);
                }}
                className="form-control"
                autoComplete={false.toString()}
              />
            </label>
          </div>
          <div className={settingsStyles['input-row']}>
            <label
              className={settingsStyles['input-label']}
              htmlFor="new-password-input"
            >
              New Password
              <input
                id="new-password-input"
                type="password"
                placeholder="********"
                value={newPassword}
                onChange={e => {
                  const val = e.target.value;

                  setNewPassword(val);
                }}
                className="form-control"
              />
            </label>
          </div>
          <div className={settingsStyles['input-row']}>
            <label
              className={settingsStyles['input-label']}
              htmlFor="new-password-confirmation-input"
            >
              New Password
              <input
                id="new-password-confirmation-input"
                type="password"
                placeholder="********"
                value={newPasswordConfirmation}
                onChange={e => {
                  const val = e.target.value;

                  setNewPasswordConfirmation(val);
                }}
                className="form-control"
              />
            </label>
          </div>

          <div className={settingsStyles.actions}>
            <Button type="submit" kind="first" isBusy={inProgress}>
              Update
            </Button>

            <Button
              type="button"
              kind="second"
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
}

export default PasswordModal;
