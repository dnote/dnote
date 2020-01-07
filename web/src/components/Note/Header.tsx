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
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import { NoteData } from 'jslib/operations/types';
import { getHomePath } from 'web/libs/paths';
import BookIcon from '../Icons/Book';
import HeaderRight from './HeaderRight';

import noteStyles from '../Common/Note/Note.scss';
import styles from './Header.scss';

interface Props {
  note: NoteData;
  isOwner: boolean;
  collapsed?: boolean;
}

const Header: React.FunctionComponent<Props> = ({
  note,
  isOwner,
  collapsed
}) => {
  let fill;
  if (collapsed) {
    fill = '#8c8c8c';
  } else {
    fill = '#000000';
  }

  return (
    <header className={noteStyles.header}>
      <div className={noteStyles['header-left']}>
        <BookIcon
          fill={fill}
          width={20}
          height={20}
          className={styles['book-icon']}
        />

        <h1
          className={classnames(noteStyles['book-label'], styles['book-label'])}
        >
          <Link to={getHomePath({ book: note.book.label })}>
            {note.book.label}
          </Link>
        </h1>
      </div>

      <div className={noteStyles['header-right']}>
        <HeaderRight isOwner={isOwner} isPublic={note.public} />
      </div>
    </header>
  );
};

export default Header;
