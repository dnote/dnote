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

import TrashIcon from '../../Icons/Trash';
import Tooltip from '../../Common/Tooltip';
import styles from './Actions.module.scss';

function Actions({ bookUUID, onDeleteBook, shown }) {
  return (
    <div
      className={classnames(styles['actions-wrapper'], {
        [styles.shown]: shown
      })}
    >
      <Tooltip
        id="tooltip-delete-book"
        alignment="right"
        direction="bottom"
        overlay={<span>Delete this book</span>}
        wrapperClassName={styles['action-tooltip-wrapper']}
        triggerClassName={styles['action-tooltip-trigger']}
      >
        <button
          type="button"
          className={classnames(
            'T-delete-book-btn button-no-ui',
            styles.action
          )}
          onClick={() => {
            onDeleteBook(bookUUID);
          }}
        >
          <TrashIcon width="16" height="16" />
        </button>
      </Tooltip>
    </div>
  );
}

export default Actions;
