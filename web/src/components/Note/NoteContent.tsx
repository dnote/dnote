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

/* eslint-disable react/no-danger */

import React from 'react';
import classnames from 'classnames';
import { Link } from 'react-router-dom';

import BookIcon from '../Icons/Book';
import { parseMarkdown } from '../../helpers/markdown';
import { nanosecToMillisec, getShortMonthName } from '../../helpers/time';
import { useSelector } from '../../store';
import { getNoteEditPath } from '../../libs/paths';
import styles from './NoteContent.scss';

function formatAddedOn(ts: number): string {
  const ms = nanosecToMillisec(ts);
  const d = new Date(ms);

  const month = getShortMonthName(d);
  const date = d.getDate();
  const year = d.getFullYear();

  return `${month} ${date}, ${year}`;
}

function getDatetimeISOString(ts: number): string {
  const ms = nanosecToMillisec(ts);

  return new Date(ms).toISOString();
}

interface Props {}

const Content: React.SFC<Props> = () => {
  const { note, user } = useSelector(state => {
    return {
      note: state.note.data,
      user: state.auth.user.data
    };
  });

  return (
    <article className={styles.frame}>
      <header className={styles.header}>
        <BookIcon
          fill="#000000"
          width={20}
          height={20}
          className={styles['book-icon']}
        />

        <h1 className={styles['book-label']}>{note.book.label}</h1>
      </header>

      <section
        className={classnames('markdown-body', styles.content)}
        dangerouslySetInnerHTML={{
          __html: parseMarkdown(note.content)
        }}
      />

      <footer className={styles.footer}>
        <time
          className={styles.ts}
          dateTime={getDatetimeISOString(note.added_on)}
        >
          {formatAddedOn(note.added_on)}
        </time>

        {note.user.uuid === user.uuid && (
          <div className={styles.actions}>
            <Link to={getNoteEditPath(note.uuid)} className={styles.action}>
              Edit
            </Link>
          </div>
        )}
      </footer>
    </article>
  );
};

export default React.memo(Content);
