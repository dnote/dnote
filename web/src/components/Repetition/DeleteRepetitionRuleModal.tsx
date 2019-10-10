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
import classnames from 'classnames';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import services from 'web/libs/services';
import { RepetitionRuleData } from 'jslib/operations/types';
import Modal, { Header, Body } from '../Common/Modal';
import Flash from '../Common/Flash';
import { removeRepetitionRule } from '../../store/repetitionRules';
import { useSelector, useDispatch } from '../../store';
import Button from '../Common/Button';
import styles from './DeleteRepetitionRuleModal.scss';

function getRepetitionRuleByUUID(
  repetitionRules,
  uuid
): RepetitionRuleData | null {
  for (let i = 0; i < repetitionRules.length; ++i) {
    const r = repetitionRules[i];

    if (r.uuid === uuid) {
      return r;
    }
  }

  return null;
}

interface Props extends RouteComponentProps {
  isOpen: boolean;
  onDismiss: () => void;
  setSuccessMessage: (string) => void;
  repetitionRuleUUID: string;
}

const DeleteRepetitionModal: React.FunctionComponent<Props> = ({
  isOpen,
  onDismiss,
  setSuccessMessage,
  repetitionRuleUUID
}) => {
  const [inProgress, setInProgress] = useState(false);
  const [errMessage, setErrMessage] = useState('');
  const dispatch = useDispatch();

  const { repetitionRules } = useSelector(state => {
    return {
      repetitionRules: state.repetitionRules
    };
  });

  const rule = getRepetitionRuleByUUID(
    repetitionRules.data,
    repetitionRuleUUID
  );

  const labelId = 'delete-rule-modal-label';
  const nameInputId = 'delete-rule-modal-name-input';
  const descId = 'delete-rule-modal-desc';

  useEffect(() => {
    if (!isOpen) {
      setErrMessage('');
    }
  }, [isOpen]);

  if (rule === null) {
    return null;
  }

  return (
    <Modal
      modalId="T-delete-rule-modal"
      isOpen={isOpen}
      onDismiss={onDismiss}
      ariaLabelledBy={labelId}
      ariaDescribedBy={descId}
      size="small"
    >
      <Header
        labelId={labelId}
        heading="Delete the repetition rule"
        onDismiss={onDismiss}
      />

      <Flash
        kind="danger"
        onDismiss={() => {
          setErrMessage('');
        }}
        hasBorder={false}
        when={Boolean(errMessage)}
        noMargin
      >
        {errMessage}
      </Flash>

      <Flash kind="warning" id={descId} noMargin>
        <span>
          This action will permanently remove the following repetition rule:{' '}
        </span>
        <span className={styles['rule-label']}>{rule.title}</span>
      </Flash>

      <Body>
        <form
          onSubmit={e => {
            e.preventDefault();

            setSuccessMessage('');
            setInProgress(true);

            services.repetitionRules
              .remove(repetitionRuleUUID)
              .then(() => {
                dispatch(removeRepetitionRule(repetitionRuleUUID));
                setInProgress(false);
                onDismiss();

                // Scroll to top so that the message is visible.
                setSuccessMessage(
                  `Successfully removed the rule "${rule.title}"`
                );
                window.scrollTo(0, 0);
              })
              .catch(err => {
                console.log('Error deleting rule', err);
                setInProgress(false);
                setErrMessage(err.message);
              });
          }}
        >
          <div className={styles.actions}>
            <Button
              type="button"
              size="normal"
              kind="second"
              onClick={onDismiss}
              disabled={inProgress}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              size="normal"
              kind="danger"
              disabled={inProgress}
              isBusy={inProgress}
            >
              Delete
            </Button>
          </div>
        </form>
      </Body>
    </Modal>
  );
};

export default withRouter(DeleteRepetitionModal);
