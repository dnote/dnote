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

import React, { Fragment } from 'react';
import classnames from 'classnames';

import BookItem from './BookItem';
import BookHolder from './BookHolder';
import { BookData } from 'jslib/operations/types';
import styles from './BookList.scss';

function Placeholder() {
  const ret = [];

  for (let i = 0; i < 12; i++) {
    ret.push(<BookHolder key={i} />);
  }

  return <Fragment>{ret}</Fragment>;
}

interface Props {
  isFetching: boolean;
  isFetched: boolean;
  books: BookData[];
  focusedIdx: number;
  setFocusedOptEl: (HTMLElement) => void;
  onDeleteBook: (string) => void;
}

const BookList: React.SFC<Props> = ({
  isFetching,
  isFetched,
  books,
  focusedIdx,
  setFocusedOptEl,
  onDeleteBook
}) => {
  return (
    <ul
      id="book-list"
      className={classnames(styles.list, { loaded: isFetched })}
    >
      {isFetching ? (
        <Placeholder />
      ) : (
        books.map((book, idx) => {
          const isFocused = idx === focusedIdx;

          return (
            <BookItem
              key={book.uuid}
              book={book}
              isFocused={isFocused}
              setFocusedOptEl={setFocusedOptEl}
              onDeleteBook={onDeleteBook}
            />
          );
        })
      )}
    </ul>
  );
};

export default React.memo(BookList);
