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

interface Props {
  id: string;
  value: string;
  itemClassName: string;
  children: React.ReactNode;
  disabled: boolean;
  setIsOpen: (boolean) => void;
  selectedOptRef: React.MutableRefObject<any>;
  setFocusedOptEl: (HTMLElement) => void;
  isFocused: boolean;
  isSelected?: boolean;
}

const Item: React.SFC<Props> = ({
  children,
  id,
  value,
  disabled,
  itemClassName,
  setIsOpen,
  selectedOptRef,
  setFocusedOptEl,
  isSelected,
  isFocused
}) => {
  return (
    <li
      id={id}
      key={value}
      className={classnames(itemClassName)}
      role="none"
      onClick={() => {
        if (disabled) {
          return;
        }

        setIsOpen(false);
      }}
      ref={el => {
        if (isSelected) {
          // eslint-disable-next-line no-param-reassign
          selectedOptRef.current = el;
        }

        if (isFocused) {
          setFocusedOptEl(el);
        }
      }}
    >
      {children}
    </li>
  );
};

export default Item;
