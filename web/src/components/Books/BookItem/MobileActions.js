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

import React, { useState, useRef } from 'react';
import classnames from 'classnames';

import Menu from '../../Common/Menu';
import DotsIcon from '../../Icons/Dots';
import styles from './MobileActions.module.scss';

function MobileActions({ bookUUID, onDeleteBook }) {
  const [isOpen, setIsOpen] = useState(false);

  const optRefs = [useRef(null)];
  const options = [
    {
      name: 'home',
      value: (
        <button
          ref={optRefs[0]}
          type="button"
          className={classnames(
            'button-no-ui button-stretch',
            styles.action,
            styles.danger
          )}
          onClick={() => {
            setIsOpen(false);
            onDeleteBook(bookUUID);
          }}
        >
          Remove
        </button>
      )
    }
  ];

  return (
    <Menu
      options={options}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      optRefs={optRefs}
      menuId="mobile-book-actions"
      triggerContent={<DotsIcon width="12" height="12" />}
      wrapperClassName={styles.wrapper}
      triggerClassName={styles.trigger}
      contentClassName={styles.content}
      alignment="top"
      direction="left"
    />
  );
}

export default MobileActions;
