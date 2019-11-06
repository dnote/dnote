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

import { getNoteEditPath, getHomePath } from 'web/libs/paths';
import { tokenize, TokenKind } from 'web/libs/fts/lexer';
import BookIcon from '../Icons/Book';
import GlobeIcon from '../Icons/Globe';
import { parseMarkdown } from '../../helpers/markdown';
import { nanosecToMillisec, getMonthName } from '../../helpers/time';
import formatTime from '../../helpers/time/format';
import { useSelector } from '../../store';
import Time from '../Common/Time';
import Tooltip from '../Common/Tooltip';
import styles from './Content.scss';

function formatAddedOn(ts: number): string {
  const ms = nanosecToMillisec(ts);
  const d = new Date(ms);

  return formatTime(d, '%MMMM %DD, %YYYY');
}

function formatFTSSelection(content: string): string {
  if (content.indexOf('<dnotehl>') === -1) {
    return content;
  }

  const tokens = tokenize(content);

  let output = '';
  let buf = [];

  for (let i = 0; i < tokens.length; i++) {
    const t = tokens[i];

    if (t.kind === TokenKind.hlBegin || t.kind === TokenKind.eol) {
      output += buf.join('');

      buf = [];
    } else if (t.kind === TokenKind.hlEnd) {
      output += `<span class="${styles.match}">
        ${buf.join('')}
      </span>`;

      buf = [];
    } else {
      buf.push(t.value);
    }
  }

  return output;
}

function formatContent(content: string): string {
  const formatted = formatFTSSelection(content);
  return parseMarkdown(formatted);
}

interface Props {
  onDeleteModalOpen: () => void;
  onShareModalOpen: () => void;
}

const Content: React.FunctionComponent<Props> = ({
  onDeleteModalOpen,
  onShareModalOpen
}) => {
  const { note, user } = useSelector(state => {
    return {
      note: state.note.data,
      user: state.auth.user.data
    };
  });

  const publicTooltip = 'Anyone on the Internet can see this note.';

  return (
    <article className={styles.frame}>
      <header className={styles.header}>
        <div className={styles['header-left']}>
          <BookIcon
            fill="#000000"
            width={20}
            height={20}
            className={styles['book-icon']}
          />

          <h1 className={styles['book-label']}>
            <Link
              to={getHomePath({ book: note.book.label })}
              className={styles['book-label-link']}
            >
              {note.book.label}
            </Link>
          </h1>
        </div>

        <div className={styles['header-right']}>
          {note.public && (
            <Tooltip
              id="note-public-indicator"
              alignment="right"
              direction="bottom"
              overlay={publicTooltip}
            >
              <GlobeIcon
                fill="#8c8c8c"
                width={16}
                height={16}
                ariaLabel={publicTooltip}
              />
            </Tooltip>
          )}
        </div>
      </header>

      <section
        className={classnames('markdown-body', styles.content)}
        dangerouslySetInnerHTML={{
          __html: formatContent(note.content)
        }}
      />

      <footer className={styles.footer}>
        <Time
          id="note-ts"
          text={formatAddedOn(note.added_on)}
          ms={nanosecToMillisec(note.added_on)}
          wrapperClassName={styles.ts}
          tooltipAlignment="left"
          tooltipDirection="bottom"
        />

        {note.user.uuid === user.uuid && (
          <div className={styles.actions}>
            <button
              id="T-share-note-button"
              type="button"
              className={classnames('button-no-ui', styles.action)}
              onClick={e => {
                e.preventDefault();

                onShareModalOpen();
              }}
            >
              Share
            </button>

            <button
              id="T-delete-note-button"
              type="button"
              className={classnames('button-no-ui', styles.action)}
              onClick={e => {
                e.preventDefault();

                onDeleteModalOpen();
              }}
            >
              Delete
            </button>

            <Link
              id="T-edit-note-button"
              to={getNoteEditPath(note.uuid)}
              className={styles.action}
            >
              Edit
            </Link>
          </div>
        )}
      </footer>
    </article>
  );
};

export default React.memo(Content);
