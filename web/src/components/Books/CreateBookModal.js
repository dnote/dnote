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
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';

import Modal, { Header } from '../Common/Modal';
import * as booksOperation from '../../operations/books';
import { addBook } from '../../actions/books';
import { getHomePath } from '../../libs/paths';
import Button from '../Common/Button';
import Flash from '../Common/Flash';
import { checkDuplicate, validateBookName } from '../../libs/books';

import styles from './CreateBookModal.module.scss';
import bodyStyles from '../Common/Modal/ModalBody.module.scss';

function CreateBookModal({
  isOpen,
  onDismiss,
  doAddBook,
  history,
  booksData,
  demo
}) {
  const [inProgress, setInProgress] = useState(false);
  const [bookName, setBookName] = useState('');
  const [errMessage, setErrMessage] = useState('');

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
          type="danger"
          onDismiss={() => {
            setErrMessage('');
          }}
          hasBorder={false}
          id={msgId}
        >
          {errMessage}
        </Flash>
      )}

      <form
        className={bodyStyles.wrapper}
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
          if (checkDuplicate(booksData.items, bookName)) {
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

          booksOperation
            .create({ name: bookName })
            .then(book => {
              doAddBook(book);
              setInProgress(false);

              const dest = getHomePath({ book: book.uuid });
              history.push(dest);
            })
            .catch(err => {
              console.log('Error creating a book', err);
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
            kind="second"
            onClick={onDismiss}
            disabled={inProgress}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            kind="third"
            disabled={demo || inProgress}
            isBusy={inProgress}
          >
            Create
          </Button>
        </div>
      </form>
    </Modal>
  );
}

function mapStateToProps(state) {
  return {
    booksData: state.books
  };
}

const mapDispatchToProps = {
  doAddBook: addBook
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(CreateBookModal)
);
