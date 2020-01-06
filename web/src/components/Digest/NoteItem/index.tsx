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

import React, { Fragment, useState } from 'react';
import { Link } from 'react-router-dom';

import { DigestNoteData } from 'jslib/operations/types';
import { getNotePath } from 'web/libs/paths';
import Note from '../../Common/Note';
import Flash from '../../Common/Flash';
import NoteItemHeader from './Header';
import styles from '../Digest.scss';

interface Props {
  note: DigestNoteData;
  onSetReviewed: (string, boolean) => Promise<any>;
}

const NoteItem: React.FunctionComponent<Props> = ({ note, onSetReviewed }) => {
  const [collapsed, setCollapsed] = useState(note.isReviewed);
  const [errorMessage, setErrMessage] = useState('');

  return (
    <li className={styles.item}>
      <Note
        collapsed={collapsed}
        note={note}
        header={
          <Fragment>
            <NoteItemHeader
              note={note}
              collapsed={collapsed}
              setCollapsed={setCollapsed}
              onSetReviewed={onSetReviewed}
              setErrMessage={setErrMessage}
            />

            <Flash kind="danger" when={errorMessage !== ''}>
              {errorMessage}
            </Flash>
          </Fragment>
        }
        footerActions={
          <Link
            className={styles.action}
            to={getNotePath(note.uuid)}
            rel="noopener noreferrer"
            target="_blank"
          >
            Go to note ›
          </Link>
        }
        footerUseTimeAgo
      />
    </li>
  );
};

export default NoteItem;
