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
import { Link, withRouter, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';

import { getDigestPath } from 'web/libs/paths';
import { parseSearchString } from 'jslib/helpers/url';
import { blacklist } from 'jslib/helpers/obj';
import SelectMenu from '../../Common/PageToolbar/SelectMenu';
import selectMenuStyles from '../../Common/PageToolbar/SelectMenu.scss';
import { Status } from '../types';
import styles from './Toolbar.scss';

interface Props extends RouteComponentProps {
  digestUUID: string;
  status: Status;
  disabled?: boolean;
}

const StatusMenu: React.FunctionComponent<Props> = ({
  digestUUID,
  status,
  disabled,
  location
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const optRefs = [useRef(null), useRef(null), useRef(null)];
  const searchObj = parseSearchString(location.search);

  const options = [
    {
      name: 'all',
      value: (
        <Link
          role="menuitem"
          className={selectMenuStyles.link}
          to={getDigestPath(digestUUID, blacklist(searchObj, ['status']))}
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
      name: 'unreviewed',
      value: (
        <Link
          role="menuitem"
          className={selectMenuStyles.link}
          to={getDigestPath(digestUUID, {
            ...searchObj,
            status: Status.Unreviewed
          })}
          onClick={() => {
            setIsOpen(false);
          }}
          ref={optRefs[1]}
          tabIndex={-1}
        >
          Unreviewed
        </Link>
      )
    },
    {
      name: 'reviewed',
      value: (
        <Link
          role="menuitem"
          className={selectMenuStyles.link}
          to={getDigestPath(digestUUID, {
            ...searchObj,
            status: Status.Reviewed
          })}
          onClick={() => {
            setIsOpen(false);
          }}
          ref={optRefs[2]}
          tabIndex={-1}
        >
          Reviewed
        </Link>
      )
    }
  ];

  const isActive = status === Status.Reviewed || status === Status.Unreviewed;

  let defaultCurrentOptionIdx: number;
  let statusText: string;
  if (status === Status.Reviewed) {
    defaultCurrentOptionIdx = 2;
    statusText = 'Reviewed';
  } else if (status === Status.Unreviewed) {
    defaultCurrentOptionIdx = 1;
    statusText = 'Unreviewed';
  } else {
    defaultCurrentOptionIdx = 0;
    statusText = 'All';
  }

  return (
    <SelectMenu
      wrapperClassName={styles['menu-trigger']}
      defaultCurrentOptionIdx={defaultCurrentOptionIdx}
      options={options}
      disabled={disabled}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      optRefs={optRefs}
      triggerId="status-menu-trigger"
      menuId="status-menu"
      headerText="Filter by status"
      triggerClassName={classnames('button-no-padding', {
        [styles['active-menu-trigger']]: isActive
      })}
      triggerText={` Status: ${statusText} `}
      alignment="right"
      direction="bottom"
    />
  );
};

export default withRouter(StatusMenu);
