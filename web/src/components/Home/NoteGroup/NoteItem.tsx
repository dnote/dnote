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
import moment from 'moment';
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import { getNotePath } from 'web/libs/paths';
import { tokenize, TokenKind } from 'web/libs/fts/lexer';
import { NoteData } from 'jslib/operations/types';
import { excerpt } from 'web/libs/string';
import { Filters } from 'jslib/helpers/filters';
import { nanosecToSec } from '../../../helpers/time';
import styles from './NoteItem.scss';

function formatFTSSnippet(content: string): React.ReactNode[] {
  const tokens = tokenize(content);

  const output: React.ReactNode[] = [];
  let buf = [];

  for (let i = 0; i < tokens.length; i++) {
    const t = tokens[i];

    if (t.kind === TokenKind.hlBegin || t.kind === TokenKind.eol) {
      output.push(buf.join(''));

      buf = [];
    } else if (t.kind === TokenKind.hlEnd) {
      const comp = (
        <span className={styles.match} key={i}>
          {buf.join('')}
        </span>
      );
      output.push(comp);

      buf = [];
    } else {
      buf.push(t.value);
    }
  }

  return output;
}

function renderContent(content: string): React.ReactNode[] {
  if (content.indexOf('<dnotehl>') > -1) {
    return formatFTSSnippet(content);
  }

  return excerpt(content, 160);
}

interface Props {
  note: NoteData;
  filters: Filters;
}

const NoteItem: React.SFC<Props> = ({ note, filters }) => {
  return (
    <li className={classnames('T-note-item', styles.wrapper, {})}>
      <Link
        className={styles.link}
        to={getNotePath(note.uuid, filters)}
        draggable={false}
      >
        <div className={styles.body}>
          <div className={styles.header}>
            <h3 className={styles['book-label']}>{note.book.label}</h3>
            <div className={styles.ts}>
              {moment.unix(nanosecToSec(note.added_on)).fromNow()}
            </div>
          </div>
          <div className={styles.content}>{renderContent(note.content)}</div>
        </div>
      </Link>
    </li>
  );
};

export default NoteItem;
