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
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import { DigestNoteData } from 'jslib/operations/types';
import { getHomePath } from 'web/libs/paths';
import ReviewButton from './ReviewButton';
import Button from '../../Common/Button';
import CaretIcon from '../../Icons/CaretSolid';

import noteStyles from '../../Common/Note/Note.scss';
import styles from './Header.scss';

interface Props {
  note: DigestNoteData;
  setCollapsed: (boolean) => void;
  onSetReviewed: (string, boolean) => Promise<any>;
  setErrMessage: (string) => void;
  collapsed: boolean;
}

const Header: React.FunctionComponent<Props> = ({
  note,
  collapsed,
  setCollapsed,
  onSetReviewed,
  setErrMessage
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
        <Button
          className={styles['header-action']}
          type="button"
          kind="no-ui"
          onClick={() => {
            setCollapsed(!collapsed);
          }}
        >
          <CaretIcon
            fill={fill}
            width={12}
            height={12}
            className={classnames({ [styles['caret-collapsed']]: collapsed })}
          />
        </Button>

        <h1
          className={classnames(noteStyles['book-label'], styles['book-label'])}
        >
          <Link to={getHomePath({ book: note.book.label })}>
            {note.book.label}
          </Link>
        </h1>
      </div>

      <div className={noteStyles['header-right']}>
        <ReviewButton
          isReviewed={note.isReviewed}
          noteUUID={note.uuid}
          setCollapsed={setCollapsed}
          onSetReviewed={onSetReviewed}
          setErrMessage={setErrMessage}
        />
      </div>
    </header>
  );
};

export default Header;
