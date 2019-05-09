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

import * as paymentService from '../../../services/payment';
import Button from '../../Common/Button';
import Modal, { Header, Body } from '../../Common/Modal';

import settingsStyles from '../Settings.module.scss';

function CancelPlanModal({
  isOpen,
  onDismiss,
  subscriptionId,
  setSuccessMsg,
  setFailureMsg,
  doGetSubscription
}) {
  const [inProgress, setInProgress] = useState(false);

  async function handleSubmit(e) {
    e.preventDefault();

    setSuccessMsg('');
    setFailureMsg('');
    setInProgress(true);

    try {
      await paymentService.cancelSubscription({ subscriptionId });
      await doGetSubscription();

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

          <div className={settingsStyles.actions}>
            <Button
              type="button"
              kind="first"
              isBusy={inProgress}
              onClick={onDismiss}
            >
              No, I changed my mind. Go back.
            </Button>

            <Button type="submit" kind="second" isBusy={inProgress}>
              Cancel my plan.
            </Button>
          </div>
        </form>
      </Body>
    </Modal>
  );
}

const mapDispatchToProps = {};

export default connect(
  null,
  mapDispatchToProps
)(CancelPlanModal);
