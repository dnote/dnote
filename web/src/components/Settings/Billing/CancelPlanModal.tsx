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

import services from 'web/libs/services';
import { useDispatch } from '../../../store';
import Button from '../../Common/Button';
import Modal, { Header, Body } from '../../Common/Modal';
import modalStyles from '../../Common/Modal/Modal.scss';
import { getSubscription } from '../../../store/auth';

interface Props {
  isOpen: boolean;
  onDismiss: () => void;
  subscriptionId: string;
  setSuccessMsg: (string) => void;
  setFailureMsg: (string) => void;
}

const CancelPlanModal: React.SFC<Props> = ({
  isOpen,
  onDismiss,
  subscriptionId,
  setSuccessMsg,
  setFailureMsg
}) => {
  const dispatch = useDispatch();
  const [inProgress, setInProgress] = useState(false);

  async function handleSubmit(e) {
    e.preventDefault();

    setSuccessMsg('');
    setFailureMsg('');
    setInProgress(true);

    try {
      await services.payment.cancelSubscription({ subscriptionId });
      await dispatch(getSubscription());

      setSuccessMsg(
        'Your subscription is cancelled. You can still continue using Dnote until the end of billing cycle.'
      );
      setInProgress(false);
      onDismiss();
    } catch (err) {
      setFailureMsg(`Failed to cancel the subscription. Error: ${err.message}`);
      setInProgress(false);
    }
  }

  const labelId = 'plan-cancel-modal';

  return (
    <Modal isOpen={isOpen} onDismiss={onDismiss} ariaLabelledBy={labelId}>
      <Header
        labelId={labelId}
        heading="Cancel the plan"
        onDismiss={onDismiss}
      />

      <Body>
        <form onSubmit={handleSubmit} autoComplete="off">
          <div>Sorry to see you go. Hope Dnote was helpful to you.</div>

          <div className={modalStyles.actions}>
            <Button
              type="button"
              kind="first"
              size="small"
              isBusy={inProgress}
              onClick={onDismiss}
            >
              No, I changed my mind. Go back.
            </Button>

            <Button
              type="submit"
              kind="second"
              size="small"
              isBusy={inProgress}
            >
              Cancel my plan.
            </Button>
          </div>
        </form>
      </Body>
    </Modal>
  );
};

export default CancelPlanModal;
