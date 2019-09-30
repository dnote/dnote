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

import operations from 'web/libs/operations';
import Modal, { Header, Body } from '../Common/Modal';
import Flash from '../Common/Flash';
import { removeBook } from '../../store/books';
import { useSelector, useDispatch } from '../../store';
import Button from '../Common/Button';

import styles from './DeleteBookModal.scss';

function getBookByUUID(books, uuid) {
  for (let i = 0; i < books.length; ++i) {
    const book = books[i];

    if (book.uuid === uuid) {
      return book;
    }
  }

  return {};
}

interface Props extends RouteComponentProps {
  isOpen: boolean;
  onDismiss: () => void;
  setSuccessMessage: (string) => void;
  bookUUID: string;
}

const DeleteBookModal: React.SFC<Props> = ({
  isOpen,
  onDismiss,
  setSuccessMessage,
  bookUUID
}) => {
  const [inProgress, setInProgress] = useState(false);
  const [bookLabel, setBookLabel] = useState('');
  const [errMessage, setErrMessage] = useState('');
  const dispatch = useDispatch();

  const { books } = useSelector(state => {
    return {
      books: state.books
    };
  });

  const book = getBookByUUID(books.data, bookUUID);

  const labelId = 'delete-book-modal-label';
  const nameInputId = 'delete-book-modal-name-input';
  const descId = 'delete-book-modal-desc';

  useEffect(() => {
    if (!isOpen) {
      setBookLabel('');
      setErrMessage('');
    }
  }, [isOpen]);

  return (
    <Modal
      modalId="T-delete-book-modal"
      isOpen={isOpen}
      onDismiss={onDismiss}
      ariaLabelledBy={labelId}
      ariaDescribedBy={descId}
      size="small"
    >
      <Header
        labelId={labelId}
        heading="Delete the book"
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
          This action will permanently remove the following book including its
          notes:
        </span>
        <span className={styles['book-label']}>{book.label}</span>
      </Flash>

      <Body>
        <form
          onSubmit={e => {
            e.preventDefault();

            setSuccessMessage('');
            setInProgress(true);

            if (bookLabel !== book.label) {
              setErrMessage('The book label did not match');
              setInProgress(false);
              return;
            }

            operations.books
              .remove(bookUUID)
              .then(() => {
                dispatch(removeBook(bookUUID));
                setInProgress(false);
                onDismiss();

                // Scroll to top so that the message is visible.
                setSuccessMessage(
                  `Successfully removed the book: ${book.label}`
                );
                window.scrollTo(0, 0);
              })
              .catch(err => {
                console.log('Error deleting book', err);
                setInProgress(false);
                setErrMessage(err.message);
              });
          }}
        >
          <label htmlFor={nameInputId} className={styles.label}>
            <div className={styles['label-text']}>
              To confirm, please enter the label of the book.
            </div>
            <input
              id={nameInputId}
              autoFocus
              type="text"
              placeholder="Wisdom"
              className={classnames('text-input', styles.input)}
              value={bookLabel}
              onChange={e => {
                const val = e.target.value;
                setBookLabel(val);
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

export default withRouter(DeleteBookModal);
