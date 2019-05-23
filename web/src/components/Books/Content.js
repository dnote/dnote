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

import React, { useState, useEffect, useRef } from 'react';
import { withRouter } from 'react-router';
import classnames from 'classnames';

import BookList from './BookList';
import EmptyList from './EmptyList';
import SearchInput from '../Common/SearchInput';
import Flash from '../Common/Flash';
import Button from '../Common/Button';

import { escapesRegExp } from '../../libs/string';
import { getHomePath } from '../../libs/paths';
import DeleteBookModal from './DeleteBookModal';
import { useSearchMenuKeydown, useScrollToFocused } from '../../libs/hooks/dom';
import { getOptIdxByValue } from '../../helpers/accessibility';
import styles from './Content.module.scss';

function filterBooks(books, searchInput) {
  if (!searchInput) {
    return books;
  }

  const input = escapesRegExp(searchInput);
  const re = new RegExp(input, 'i');

  return books.filter(book => {
    return re.test(book.label);
  });
}

function handleMenuKeydownSelect(demo, history) {
  return option => {
    const destination = getHomePath(
      {
        book: option.uuid
      },
      { demo }
    );

    history.push(destination);
  };
}

function Content({
  books,
  isFetching,
  isFetched,
  demo,
  location,
  match,
  history,
  containerEl,
  onStartCreateBook
}) {
  const [searchValue, setSearchValue] = useState('');
  const [searchFocus, setSearchFocus] = useState('');
  const [focusedIdx, setFocusedIdx] = useState(-1);
  const [focusedOptEl, setFocusedOptEl] = useState(null);
  const [successMessage, setSuccessMessage] = useState('');
  const [bookUUIDToDelete, setBookUUIDToDelete] = useState('');
  const inputRef = useRef(null);

  const filterdBooks = filterBooks(books, searchValue);

  const currentIdx = getOptIdxByValue(filterdBooks, searchValue);
  useEffect(() => {
    if (searchFocus) {
      setFocusedIdx(currentIdx);
    } else {
      setFocusedIdx(-1);
    }
  }, [searchValue, searchFocus, currentIdx, containerEl]);

  useSearchMenuKeydown({
    options: filterdBooks,
    containerEl,
    focusedIdx,
    setFocusedIdx,
    onKeydownSelect: handleMenuKeydownSelect(demo, history),
    location,
    match,
    history
  });
  useScrollToFocused({
    shouldScroll: true,
    offset: -80,
    focusedOptEl,
    containerEl
  });

  useEffect(() => {
    if (!isFetching) {
      if (inputRef.current) {
        inputRef.current.focus();
      }
    }
  }, [isFetching]);

  const hasNoBooks = isFetched && books.length === 0;

  return (
    <div>
      {successMessage && (
        <Flash
          type="success"
          onDismiss={() => {
            setSuccessMessage('');
          }}
        >
          {successMessage}
        </Flash>
      )}

      <div className={styles.actions}>
        <SearchInput
          size="medium"
          placeholder="Find a book"
          value={searchValue}
          setValue={setSearchValue}
          inputClassName={classnames(
            'text-input-large',
            styles['search-input']
          )}
          wrapperClassName={styles['search-input-wrapper']}
          disabled={isFetching || hasNoBooks}
          inputRef={inputRef}
          onFocus={() => {
            setSearchFocus(true);
          }}
          onBlur={() => {
            setSearchFocus(false);
          }}
        />

        <Button
          id="T-create-book-btn"
          type="button"
          kind="third"
          size="normal"
          className={styles['create-book-button']}
          disabled={isFetching}
          onClick={() => {
            onStartCreateBook(true);
          }}
        >
          Create a book
        </Button>
      </div>

      {hasNoBooks ? (
        <EmptyList />
      ) : (
        <BookList
          isFetching={isFetching}
          isFetched={isFetched}
          books={filterdBooks}
          demo={demo}
          focusedIdx={focusedIdx}
          setFocusedOptEl={setFocusedOptEl}
          onDeleteBook={bookUUID => {
            setBookUUIDToDelete(bookUUID);
          }}
        />
      )}

      <DeleteBookModal
        isOpen={Boolean(bookUUIDToDelete)}
        onDismiss={() => {
          setBookUUIDToDelete(null);
        }}
        bookUUID={bookUUIDToDelete}
        setSuccessMessage={setSuccessMessage}
        containerEl={containerEl}
        demo={demo}
      />
    </div>
  );
}

export default withRouter(Content);
