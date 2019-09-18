import React, { useState } from 'react';
import classnames from 'classnames';

import { useDispatch } from '../../../store';
import { getSubscription } from '../../../store/auth';
import SettingRow from '../SettingRow';
import * as paymentService from '../../../services/payment';
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

            paymentService
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
