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

import Menu, { MenuOption } from '../../Common/Menu';
import DotsIcon from '../../Icons/Dots';
import styles from './ItemActions.scss';

interface Props {
  id: string;
  triggerId: string;
  isActive: boolean;
  options: MenuOption[];
  optRefs: React.MutableRefObject<any>[];
  isOpen: boolean;
  setIsOpen: React.Dispatch<any>;
  wrapperClassName?: string;
}

const ItemActions: React.FunctionComponent<Props> = ({
  id,
  triggerId,
  isActive,
  options,
  optRefs,
  isOpen,
  setIsOpen,
  wrapperClassName
}) => {
  return (
    <div
      className={classnames(styles.wrapper, wrapperClassName, {
        [styles['is-open']]: isOpen,
        [styles.active]: isActive
      })}
    >
      <Menu
        options={options}
        isOpen={isOpen}
        setIsOpen={setIsOpen}
        optRefs={optRefs}
        menuId={id}
        triggerId={triggerId}
        triggerContent={<DotsIcon width={12} height={12} />}
        triggerClassName={styles['trigger']}
        contentClassName={styles['content']}
        alignment="right"
        direction="bottom"
      />
    </div>
  );
};

export default ItemActions;
