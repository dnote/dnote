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

import React, { useState, useRef } from 'react';
import { Link } from 'react-router-dom';

import { getDigestsPath } from 'web/libs/paths';
import SelectMenu from '../../Common/PageToolbar/SelectMenu';
import selectMenuStyles from '../../Common/PageToolbar/SelectMenu.scss';
import { Status } from '../types';
import styles from './Toolbar.scss';

interface Props {
  status: Status;
  disabled?: boolean;
}

const StatusMenu: React.FunctionComponent<Props> = ({ status, disabled }) => {
  const [isOpen, setIsOpen] = useState(false);
  const optRefs = [useRef(null), useRef(null), useRef(null)];

  const options = [
    {
      name: 'all',
      value: (
        <Link
          role="menuitem"
          className={selectMenuStyles.link}
          to={getDigestsPath({ status: Status.All })}
          onClick={() => {
            setIsOpen(false);
          }}
          ref={optRefs[0]}
          tabIndex={-1}
        >
          All
        </Link>
      )
    },
    {
      name: 'unread',
      value: (
        <Link
          role="menuitem"
          className={selectMenuStyles.link}
          to={getDigestsPath({ status: Status.Unread })}
          onClick={() => {
            setIsOpen(false);
          }}
          ref={optRefs[1]}
          tabIndex={-1}
        >
          Unread
        </Link>
      )
    },
    {
      name: 'read',
      value: (
        <Link
          role="menuitem"
          className={selectMenuStyles.link}
          to={getDigestsPath({ status: Status.Read })}
          onClick={() => {
            setIsOpen(false);
          }}
          ref={optRefs[2]}
          tabIndex={-1}
        >
          Read
        </Link>
      )
    }
  ];

  let defaultCurrentOptionIdx: number;
  let triggerText: string;
  if (status === Status.Read) {
    defaultCurrentOptionIdx = 2;
    triggerText = 'Read';
  } else if (status === Status.Unread) {
    defaultCurrentOptionIdx = 1;
    triggerText = 'Unread';
  } else {
    defaultCurrentOptionIdx = 0;
    triggerText = 'All';
  }

  return (
    <SelectMenu
      defaultCurrentOptionIdx={defaultCurrentOptionIdx}
      options={options}
      disabled={disabled}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      optRefs={optRefs}
      triggerId="sort-menu-trigger"
      menuId="sort-menu"
      headerText="Status"
      triggerText={`Status: ${triggerText}`}
      wrapperClassName={styles['select-menu-wrapper']}
      alignment="left"
      direction="bottom"
    />
  );
};

export default StatusMenu;
