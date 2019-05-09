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

import React, { useState } from 'react';
import classnames from 'classnames';

import Popover from '../Popover';

import styles from './Tooltip.module.scss';

function Tooltip({
  id,
  alignment,
  direction,
  contentClassName,
  wrapperClassName,
  triggerClassName,
  overlay,
  children
}) {
  const [isOpen, setIsOpen] = useState(false);

  function show() {
    setIsOpen(true);
  }

  function hide() {
    setIsOpen(false);
  }

  return (
    <Popover
      renderTrigger={triggerProps => {
        return (
          <span
            className={classnames(
              triggerClassName,
              triggerProps.triggerClassName
            )}
            aria-describedby={id}
            tabIndex="-1"
            onFocus={show}
            onMouseEnter={show}
            onMouseLeave={hide}
            onBlur={hide}
          >
            {children}
          </span>
        );
      }}
      contentClassName={classnames(styles.backdrop, contentClassName)}
      wrapperClassName={classnames(styles.wrapper, wrapperClassName)}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      alignment={alignment}
      direction={direction}
      contentId={id}
      closeOnEscapeKeydown={false}
      closeOnOutsideClick={false}
      contentHasBorder={false}
      hasArrow
      renderContent={() => {
        return (
          <div
            className={classnames(styles.overlay, {
              [styles.left]: alignment === 'left',
              [styles.right]: alignment === 'right'
            })}
          >
            {overlay}
          </div>
        );
      }}
    />
  );
}

export default Tooltip;
