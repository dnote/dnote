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

import Menu, { MenuOption } from '../../Common/Menu';
import { Alignment, Direction } from '../../Common/Menu/types';
import styles from './SelectMenu.scss';

interface Props {
  defaultCurrentOptionIdx: number;
  options: MenuOption[];
  optRefs: any[];
  triggerText: string;
  isOpen: boolean;
  setIsOpen: (boolean) => void;
  triggerId: string;
  headerText: string;
  menuId: string;
  alignment: Alignment;
  alignmentMd?: Alignment;
  direction: Direction;
  disabled?: boolean;
  wrapperClassName?: string;
  triggerClassName?: string;
}

const SelectMenu: React.FunctionComponent<Props> = ({
  defaultCurrentOptionIdx,
  options,
  optRefs,
  triggerText,
  disabled,
  isOpen,
  setIsOpen,
  headerText,
  triggerId,
  menuId,
  alignment,
  alignmentMd,
  direction,
  wrapperClassName,
  triggerClassName
}) => {
  return (
    <Menu
      defaultCurrentOptionIdx={defaultCurrentOptionIdx}
      options={options}
      disabled={disabled}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      optRefs={optRefs}
      triggerId={triggerId}
      menuId={menuId}
      triggerContent={
        <div
          className={classnames(styles['trigger-content'], triggerClassName)}
        >
          {triggerText}
          <span className="dropdown-caret" />
        </div>
      }
      headerContent={<div className={styles.header}>{headerText}</div>}
      triggerClassName={classnames(styles.trigger, triggerClassName, {
        [styles['trigger-active']]: isOpen
      })}
      contentClassName={styles.content}
      wrapperClassName={wrapperClassName}
      alignment={alignment}
      alignmentMd={alignmentMd}
      direction={direction}
    />
  );
};

export default SelectMenu;
