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
import { useDispatch } from '../../../store';
import { getSubscription } from '../../../store/auth';
import SettingRow from '../SettingRow';
import styles from '../Settings.scss';

interface Props {
  subscriptionId: string;
  setSuccessMsg: (string) => void;
  setFailureMsg: (string) => void;
}

const ReactivateRow: React.SFC<Props> = ({
  subscriptionId,
  setSuccessMsg,
  setFailureMsg
}) => {
  const dispatch = useDispatch();

  const [inProgress, setInProgress] = useState(false);

  return (
    <SettingRow
      name="Reactivate your plan"
      desc="You can reactivate your plan if you have changed your mind."
      actionContent={
        <button
          className={classnames('button-no-ui', styles.edit)}
          type="button"
          disabled={inProgress}
          onClick={() => {
            setInProgress(true);

            services.payment
              .reactivateSubscription({ subscriptionId })
              .then(() => {
                return dispatch(getSubscription()).then(() => {
                  setSuccessMsg(
                    'Your plan was reactivated. The billing cycle will be the same.'
                  );
                });
              })
              .catch(err => {
                setFailureMsg(
                  `Failed to reactivate the plan. Error: ${err.message}. Please contact sung@getdnote.com.`
                );
                setInProgress(false);
              });
          }}
        >
          {inProgress ? 'Reactivating...' : 'Reactivate plan'}
        </button>
      }
    />
  );
};

export default ReactivateRow;
