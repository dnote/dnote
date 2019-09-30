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
import { withRouter, RouteComponentProps } from 'react-router-dom';

import Modal, { Header, Body } from '../Common/Modal';
import { createBook } from '../../store/books';
import { useSelector, useDispatch } from '../../store';
import Button from '../Common/Button';
import Flash from '../Common/Flash';
import { checkDuplicate, validateBookName } from 'jslib/helpers/books';

import styles from './CreateBookModal.scss';

interface Props extends RouteComponentProps {
  isOpen: boolean;
  onDismiss: () => void;
  onSuccess: () => void;
  setSuccessMessage: (string) => void;
}

const CreateBookModal: React.SFC<Props> = ({
  isOpen,
  onDismiss,
  onSuccess,
  setSuccessMessage
}) => {
  const [inProgress, setInProgress] = useState(false);
  const [bookName, setBookName] = useState('');
  const [errMessage, setErrMessage] = useState('');

  const dispatch = useDispatch();
  const { books } = useSelector(state => {
    return {
      books: state.books
    };
  });

  const labelId = 'new-book-modal-label';
  const nameInputId = 'new-book-modal-name-input';

  let msgId;
  if (errMessage) {
    msgId = 'new-book-modal-err';
  }

  return (
    <Modal
      modalId="T-create-book-modal"
      isOpen={isOpen}
      onDismiss={onDismiss}
      ariaLabelledBy={labelId}
      ariaDescribedBy={msgId}
      size="small"
    >
      <Header labelId={labelId} heading="Create a book" onDismiss={onDismiss} />

      {msgId && (
        <Flash
          kind="danger"
          onDismiss={() => {
            setErrMessage('');
          }}
          hasBorder={false}
          id={msgId}
        >
          {errMessage}
        </Flash>
      )}

      <Body>
        <form
          onSubmit={e => {
            e.preventDefault();
            setInProgress(true);

            if (!bookName) {
              setInProgress(false);
              setErrMessage('Book label is empty');
              return;
            }

            // Check if the book label already exists. If the client somehow posts a duplicate label,
            // Duplicate book labels will be resolved when they are locally synced, anyway.
            // TODO: resolve any duplicate book labels on the web as well.
            if (checkDuplicate(books.data, bookName)) {
              setInProgress(false);
              setErrMessage('Duplicate book exists');
              return;
            }

            try {
              validateBookName(bookName);
            } catch (err) {
              setInProgress(false);
              setErrMessage(err.message);
              return;
            }

            dispatch(createBook(bookName))
              .then(() => {
                setInProgress(false);

                setSuccessMessage(`Created a book: ${bookName}`);
                setInProgress(false);
                setBookName('');

                onSuccess();
                onDismiss();
              })
              .catch(err => {
                setInProgress(false);
                setErrMessage(err.message);
              });
          }}
        >
          <label htmlFor={nameInputId} className={styles.label}>
            <div className={styles['label-text']}>
              Please enter the name of the book
            </div>
            <input
              id={nameInputId}
              autoFocus
              type="text"
              placeholder="Wisdom"
              className={classnames('text-input', styles.input)}
              value={bookName}
              onChange={e => {
                const val = e.target.value;
                setBookName(val);
              }}
            />
          </label>

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
              kind="third"
              disabled={inProgress}
              isBusy={inProgress}
            >
              Create
            </Button>
          </div>
        </form>
      </Body>
    </Modal>
  );
};

export default withRouter(CreateBookModal);
