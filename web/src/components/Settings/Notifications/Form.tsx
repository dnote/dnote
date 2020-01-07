/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

import services from 'web/libs/services';
import { EmailPrefData } from 'jslib/operations/types';
import Button from '../../Common/Button';
import { receiveEmailPreference } from '../../../store/auth';
import { useDispatch } from '../../../store';
import styles from './Form.scss';

enum Action {
  setInactiveReminder,
  setProductUpdate
}

function formReducer(state, action): EmailPrefData {
  switch (action.type) {
    case Action.setInactiveReminder:
      return {
        ...state,
        inactiveReminder: action.data
      };
    case Action.setProductUpdate:
      return {
        ...state,
        productUpdate: action.data
      };
    default:
      return state;
  }
}

interface Props {
  emailPreference: EmailPrefData;
  setSuccessMsg: (string) => void;
  setFailureMsg: (string) => void;
  token?: string;
}

const Form: React.FunctionComponent<Props> = ({
  emailPreference,
  setSuccessMsg,
  setFailureMsg,
  token
}) => {
  const [inProgress, setInProgress] = useState(false);
  const dispatch = useDispatch();

  const [formState, formDispatch] = useReducer(formReducer, emailPreference);

  function handleSubmit(e) {
    e.preventDefault();

    setSuccessMsg('');
    setFailureMsg('');
    setInProgress(true);

    services.users
      .updateEmailPreference({
        inactiveReminder: formState.inactiveReminder,
        productUpdate: formState.productUpdate,
        token
      })
      .then(updatedPreference => {
        dispatch(receiveEmailPreference(updatedPreference));

        setSuccessMsg('Updated email preference');
        setInProgress(false);
      })
      .catch(err => {
        setFailureMsg(`Failed to update. Error: ${err.message}`);
        setInProgress(false);
      });
  }

  return (
    <form id="T-notifications-form" onSubmit={handleSubmit}>
      <div className={styles.section}>
        <h3 className={styles.heading}>Alerts</h3>
        <p className={styles.subtext}>Email me when:</p>
        <ul className="list-unstyled">
          <li>
            <input
              type="checkbox"
              id="inactive-reminder"
              checked={formState.inactiveReminder}
              onChange={e => {
                const { checked } = e.target;

                formDispatch({
                  type: Action.setInactiveReminder,
                  data: checked
                });
              }}
            />
            <label className={styles.label} htmlFor="inactive-reminder">
              I stop learning new things
            </label>
          </li>
        </ul>
      </div>

      <div className={styles.section}>
        <h3 className={styles.heading}>News</h3>

        <p className={styles.subtext}>Email me about:</p>
        <ul className="list-unstyled">
          <li>
            <input
              type="checkbox"
              id="product-update"
              checked={formState.productUpdate}
              onChange={e => {
                const { checked } = e.target;

                formDispatch({
                  type: Action.setProductUpdate,
                  data: checked
                });
              }}
            />
            <label className={styles.label} htmlFor="product-update">
              New features and updates
            </label>
          </li>
        </ul>
      </div>

      <div className={styles.actions}>
        <Button type="submit" kind="first" size="normal" isBusy={inProgress}>
          Update
        </Button>
      </div>
    </form>
  );
};

export default Form;
