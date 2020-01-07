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

import React, { useRef, useState } from 'react';
import { Alignment, Direction } from '../Popover/types';
import Overlay from './Overlay';

interface Props {
  id: string;
  alignment: Alignment;
  direction: Direction;
  overlay: React.ReactNode;
  children: React.ReactChild;
  contentClassName?: string;
  wrapperClassName?: string;
  triggerClassName?: string;
  noArrow?: boolean;
}

const Tooltip: React.FunctionComponent<Props> = ({
  id,
  alignment,
  direction,
  wrapperClassName,
  overlay,
  children,
  noArrow = false
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const triggerRef = useRef(null);

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
        noArrow={noArrow}
      >
        {overlay}
      </Overlay>
    </span>
  );
};

export default Tooltip;
