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

import { getNotePath } from '../../libs/paths';
import { excerpt } from '../../libs/string';
import { NoteData } from '../../operations/types';
import { nanosecToSec } from '../../helpers/time';
import styles from './NoteItem.scss';

// renderContent renders the first line of the note
function renderContent(content) {
  let linebreakIdx = content.indexOf('\n');

  if (linebreakIdx === -1) {
    linebreakIdx = content.indexOf('\r\n');
  }

  let firstline;
  if (linebreakIdx === -1) {
    firstline = content;
  } else {
    firstline = content.substr(0, linebreakIdx);
  }

  return excerpt(firstline, 70);
}

interface Props {
  note: NoteData;
}

const NoteItem: React.SFC<Props> = ({ note }) => {
  return (
    <li className={classnames('T-note-item', styles.wrapper, {})}>
      <Link
        className={styles.link}
        to={getNotePath(note.uuid)}
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
