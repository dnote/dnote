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
import { withRouter, RouteComponentProps } from 'react-router-dom';

import { homePathDef, getHomePath } from 'web/libs/paths';
import operations from 'web/libs/operations';
import Modal, { Header, Body } from '../Common/Modal';
import Flash from '../Common/Flash';
import { setMessage } from '../../store/ui';
import { useDispatch } from '../../store';
import Button from '../Common/Button';
import styles from './DeleteNoteModal.scss';

interface Props extends RouteComponentProps {
  isOpen: boolean;
  onDismiss: () => void;
  noteUUID: string;
}

const DeleteNoteModal: React.SFC<Props> = ({
  isOpen,
  onDismiss,
  noteUUID,
  history
}) => {
  const [inProgress, setInProgress] = useState(false);
  const [errMessage, setErrMessage] = useState('');
  const dispatch = useDispatch();

  const labelId = 'delete-note-modal-label';
  const descId = 'delete-note-modal-desc';

  return (
    <Modal
      modalId="T-delete-note-modal"
      isOpen={isOpen}
      onDismiss={onDismiss}
      ariaLabelledBy={labelId}
      ariaDescribedBy={descId}
      size="small"
    >
      <Header
        labelId={labelId}
        heading="Delete this note"
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
        This action will permanently remove the current note.
      </Flash>

      <Body>
        <form
          onSubmit={e => {
            e.preventDefault();

            setInProgress(true);

            operations.notes
              .remove(noteUUID)
              .then(() => {
                // dispatch(removeBook(bookUUID));
                setInProgress(false);
                onDismiss();

                // Scroll to top so that the message is visible.
                dispatch(
                  setMessage({
                    message: `Successfully removed the note`,
                    kind: 'success',
                    path: homePathDef
                  })
                );

                history.push(getHomePath());
              })
              .catch(err => {
                console.log('Error deleting note', err);
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

export default withRouter(DeleteNoteModal);
