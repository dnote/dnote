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

import BookIcon from '../../Icons/Book';
import { parseMarkdown } from '../../../helpers/markdown';
import { nanosecToMillisec, getShortMonthName } from '../../../helpers/time';

import styles from './Note.module.scss';

function formatAddedOn(ts) {
  const ms = nanosecToMillisec(ts);
  const d = new Date(ms);

  const month = getShortMonthName(d);
  const date = d.getDate();
  const year = d.getFullYear();

  return `${month} ${date}, ${year}`;
}

function Note({ note }) {
  return (
    <div className={styles.wrapper}>
      <div className={styles.header}>
        <BookIcon
          fill="#000000"
          width={20}
          height={20}
          className={styles['book-icon']}
        />

        <div className={styles['book-label']}>{note.book.label}</div>
      </div>
      <div
        className={classnames('markdown-body', styles.content)}
        // eslint-disable-next-line react/no-danger
        dangerouslySetInnerHTML={{
          __html: parseMarkdown(note.content)
        }}
      />
      <div className={styles.footer}>
        <div className={styles.ts}>{formatAddedOn(note.added_on)}</div>
      </div>
    </div>
  );
}

export default Note;
