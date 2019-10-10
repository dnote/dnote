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

import React, { Fragment, useState, useRef } from 'react';
import classnames from 'classnames';

import Overlay from './Overlay';
import { Alignment, Direction } from '../Popover/types';
import { isMobileWidth } from 'web/libs/dom';
import {
  KEYCODE_ESC,
  KEYCODE_ENTER,
  KEYCODE_SPACE
} from 'jslib/helpers/keyboard';

interface Props {
  id: string;
  alignment: Alignment;
  direction: Direction;
  overlay: React.ReactElement;
  children: React.ReactChild;
  contentClassName?: string;
  wrapperClassName?: string;
  triggerClassName?: string;
}

const Tooltip: React.FunctionComponent<Props> = ({
  id,
  alignment,
  direction,
  wrapperClassName,
  overlay,
  children
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const triggerRef = useRef(null);
  const touchingRef = useRef(false);

  function show() {
    setIsOpen(true);
  }

  function hide() {
    setIsOpen(false);
  }

  return (
    <span onMouseEnter={show} onMouseLeave={hide}>
      <span
        className={wrapperClassName}
        aria-describedby={isOpen ? id : undefined}
        ref={triggerRef}
      >
        {children}
      </span>

      <Overlay
        id={id}
        isOpen={isOpen}
        triggerEl={triggerRef.current}
        alignment={alignment}
        direction={direction}
      >
        {overlay}
      </Overlay>
    </span>
  );
};

export default Tooltip;
