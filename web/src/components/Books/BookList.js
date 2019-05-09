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

import React from 'react';
import classnames from 'classnames';

import BookItem from './BookItem';
import BookHolder from './BookHolder';
import styles from './BookList.module.scss';

const placeholders = new Array(12);
for (let i = 0; i < placeholders.length; ++i) {
  placeholders[i] = <BookHolder key={i} />;
}

function Placeholders() {
  return placeholders;
}

function Books({ books, demo, focusedIdx, setFocusedOptEl, onDeleteBook }) {
  return books.map((book, idx) => {
    const isFocused = idx === focusedIdx;

    return (
      <BookItem
        key={book.uuid}
        book={book}
        demo={demo}
        isFocused={isFocused}
        setFocusedOptEl={setFocusedOptEl}
        onDeleteBook={onDeleteBook}
      />
    );
  });
}

export default ({
  isFetching,
  isFetched,
  books,
  demo,
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
        <Placeholders />
      ) : (
        <Books
          books={books}
          demo={demo}
          focusedIdx={focusedIdx}
          setFocusedOptEl={setFocusedOptEl}
          onDeleteBook={onDeleteBook}
        />
      )}
    </ul>
  );
};
