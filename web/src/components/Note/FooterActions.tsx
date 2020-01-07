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
import classnames from 'classnames';
import { Link } from 'react-router-dom';

import { getNoteEditPath } from 'web/libs/paths';
import styles from './FooterActions.scss';

interface Props {
  isOwner: boolean;
  noteUUID: string;
  onShareModalOpen: () => void;
  onDeleteModalOpen: () => void;
}

const FooterActions: React.FunctionComponent<Props> = ({
  isOwner,
  noteUUID,
  onShareModalOpen,
  onDeleteModalOpen
}) => {
  if (!isOwner) {
    return null;
  }

  return (
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
        to={getNoteEditPath(noteUUID)}
        className={styles.action}
      >
        Edit
      </Link>
    </div>
  );
};

export default FooterActions;
