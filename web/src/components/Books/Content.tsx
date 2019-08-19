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

import React, { Fragment, useState, useEffect, useRef } from 'react';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';
import { History } from 'history';

import CreateBookModal from './CreateBookModal';
import BookList from './BookList';
import EmptyList from './EmptyList';
import SearchInput from '../Common/SearchInput';
import Button from '../Common/Button';
import { escapesRegExp } from '../../libs/string';
import { getHomePath } from '../../libs/paths';
import DeleteBookModal from './DeleteBookModal';
import {
  KeydownSelectFn,
  useSearchMenuKeydown,
  useScrollToFocused
} from '../../libs/hooks/dom';
import { useSelector } from '../../store';
import { BookData } from '../../operations/books';
import styles from './Content.scss';

function filterBooks(books: BookData[], searchInput: string): BookData[] {
  if (searchInput === '') {
    return books;
  }

  const input = escapesRegExp(searchInput);
  const re = new RegExp(input, 'i');

  return books.filter(book => {
    return re.test(book.label);
  });
}

function handleMenuKeydownSelect(history: History): KeydownSelectFn<BookData> {
  return option => {
    const destination = getHomePath({
      book: option.label
    });

    history.push(destination);
  };
}

function useFocusInput(
  isFetching: boolean,
  inputRef: React.MutableRefObject<any>
) {
  useEffect(() => {
    if (!isFetching) {
      if (inputRef.current) {
        inputRef.current.focus();
      }
    }
  }, [isFetching, inputRef]);
}

function useSetFocusedOptionOnInputFocus({
  searchValue,
  searchFocus,
  setFocusedIdx,
  containerEl
}) {
  useEffect(() => {
    if (searchFocus) {
      setFocusedIdx(0);
    } else {
      setFocusedIdx(-1);
    }
  }, [searchValue, searchFocus, containerEl, setFocusedIdx]);
}

function useFocusInputOnReset(
  searchValue: string,
  inputRef: React.MutableRefObject<any>
) {
  useEffect(() => {
    if (searchValue === '') {
      if (inputRef.current !== null) {
        inputRef.current.focus();
      }
    }
  }, [searchValue, inputRef]);
}

interface Props extends RouteComponentProps {
  setSuccessMessage: (string) => void;
}

const Content: React.SFC<Props> = ({ history, setSuccessMessage }) => {
  const { books } = useSelector(state => {
    return {
      books: state.books
    };
  });

  const [searchValue, setSearchValue] = useState('');
  const [searchFocus, setSearchFocus] = useState(false);
  const [focusedIdx, setFocusedIdx] = useState(-1);
  const [focusedOptEl, setFocusedOptEl] = useState(null);
  const [bookUUIDToDelete, setBookUUIDToDelete] = useState('');
  const [isCreateBookModalOpen, setIsCreateBookModalOpen] = useState(false);
  const inputRef = useRef(null);

  const filteredBooks = filterBooks(books.data, searchValue);
  const containerEl = document.body;

  useFocusInput(books.isFetching, inputRef);
  useSetFocusedOptionOnInputFocus({
    searchValue,
    searchFocus,
    setFocusedIdx,
    containerEl
  });
  useSearchMenuKeydown<BookData>({
    options: filteredBooks,
    containerEl,
    focusedIdx,
    setFocusedIdx,
    onKeydownSelect: handleMenuKeydownSelect(history)
  });
  useScrollToFocused({
    shouldScroll: true,
    focusedOptEl,
    containerEl
  });
  useFocusInputOnReset(searchValue, inputRef);

  const hasNoBooks = books.isFetched && filteredBooks.length === 0;

  return (
    <Fragment>
      <div className={styles.actions}>
        <SearchInput
          placeholder="Find a book"
          value={searchValue}
          onChange={e => {
            const val = e.target.value;
            setSearchValue(val);
          }}
          inputClassName={classnames(
            'text-input-small',
            styles['search-input']
          )}
          disabled={books.isFetching}
          inputRef={inputRef}
          onFocus={() => {
            setSearchFocus(true);
          }}
          onBlur={() => {
            setSearchFocus(false);
          }}
          onReset={() => {
            setSearchValue('');
          }}
        />

        <Button
          id="T-create-book-btn"
          type="button"
          kind="third"
          size="normal"
          className={styles['create-book-button']}
          disabled={books.isFetching}
          onClick={() => {
            setIsCreateBookModalOpen(true);
          }}
        >
          Create book
        </Button>
      </div>

      {hasNoBooks ? (
        <EmptyList />
      ) : (
        <BookList
          isFetching={books.isFetching}
          isFetched={books.isFetched}
          books={filteredBooks}
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
      />

      <CreateBookModal
        isOpen={isCreateBookModalOpen}
        onDismiss={() => {
          setIsCreateBookModalOpen(false);
        }}
        onSuccess={() => {
          setSearchValue('');
        }}
        setSuccessMessage={setSuccessMessage}
      />
    </Fragment>
  );
};

export default React.memo(withRouter(Content));
