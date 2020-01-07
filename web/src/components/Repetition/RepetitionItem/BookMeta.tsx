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

import React from 'react';

import { BookDomain } from 'jslib/operations/types';
import { pluralize } from 'web/libs/string';
import styles from './RepetitionItem.scss';

interface ContentProps {
  bookDomain: BookDomain;
  bookCount: number;
}

const Content: React.FunctionComponent<ContentProps> = ({
  bookDomain,
  bookCount
}) => {
  if (bookDomain === BookDomain.All) {
    return <span>From all books</span>;
  }

  let verb;
  if (bookDomain === BookDomain.Excluding) {
    verb = 'Excluding';
  } else if (bookDomain === BookDomain.Including) {
    verb = 'From';
  }

  return (
    <span>
      {verb} {bookCount} {pluralize('book', bookCount)}
    </span>
  );
};

interface Props {
  bookDomain: BookDomain;
  bookCount: number;
}

const BookMeta: React.FunctionComponent<Props> = ({
  bookDomain,
  bookCount
}) => {
  return (
    <span className={styles['book-meta']}>
      <Content bookDomain={bookDomain} bookCount={bookCount} />
    </span>
  );
};

export default BookMeta;
