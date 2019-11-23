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

import React, { Fragment } from 'react';

import PlanRow from './PlanRow';
import CancelRow from './CancelRow';
import ReactivateRow from './ReactivateRow';
import { SubscriptionData } from '../../../../store/auth/type';
import Placeholder from './Placeholder';

interface Props {
  subscription: SubscriptionData;
  setIsPlanModalOpen: (boolean) => void;
  setSuccessMsg: (string) => void;
  setFailureMsg: (string) => void;
  isFetched: boolean;
}

const PlanSection: React.FunctionComponent<Props> = ({
  subscription,
  isFetched,
  setIsPlanModalOpen,
  setSuccessMsg,
  setFailureMsg
}) => {
  if (!isFetched) {
    return <Placeholder />;
  }

  return (
    <Fragment>
      <PlanRow subscription={subscription} />

      {subscription.id && !subscription.cancel_at_period_end && (
        <CancelRow setIsPlanModalOpen={setIsPlanModalOpen} />
      )}
      {subscription.id && subscription.cancel_at_period_end && (
        <ReactivateRow
          subscriptionId={subscription.id}
          setSuccessMsg={setSuccessMsg}
          setFailureMsg={setFailureMsg}
        />
      )}
    </Fragment>
  );
};

export default PlanSection;
