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

import React, { useEffect, useRef } from 'react';
import classnames from 'classnames';

import styles from './PopoverContent.scss';
import { KEYCODE_ESC } from '../../../helpers/keyboard';
import { Alignment, Direction } from './types';

interface Props {
  contentId: string;
  contentClassName: string;
  onDismiss: () => void;
  triggerEl: HTMLElement;
  children: React.ReactNode;
  alignment: Alignment;
  direction: Direction;
  hasBorder: boolean;
  closeOnOutsideClick: boolean;
  closeOnEscapeKeydown: boolean;
}

const PopoverContent: React.SFC<Props> = ({
  contentId,
  contentClassName,
  onDismiss,
  triggerEl,
  children,
  alignment,
  direction,
  hasBorder,
  closeOnOutsideClick,
  closeOnEscapeKeydown
}) => {
  const contentRef = useRef(null);

  function handleOutsideClick(e) {
    const contentEl = contentRef.current;

    if (!contentEl || !triggerEl) {
      return;
    }

    if (triggerEl.contains(e.target)) {
      return;
    }

    if (!contentEl.contains(e.target)) {
      onDismiss();
    }
  }

  function handleKeydown(e) {
    if (e.keyCode === KEYCODE_ESC) {
      onDismiss();
    }
  }

  useEffect(() => {
    if (closeOnOutsideClick) {
      document.addEventListener('click', handleOutsideClick);
      document.addEventListener('touchstart', handleOutsideClick);
    }

    return () => {
      document.removeEventListener('click', handleOutsideClick);
      document.removeEventListener('touchstart', handleOutsideClick);
    };
  });

  useEffect(() => {
    if (closeOnEscapeKeydown) {
      document.addEventListener('keydown', handleKeydown);
    }

    return () => {
      document.addEventListener('keydown', handleKeydown);
    };
  });

  return (
    <div
      className={classnames(styles.content, contentClassName, {
        // alignment
        [styles['left-align']]: alignment === 'left',
        [styles['right-align']]: alignment === 'right',
        [styles['top-align']]: alignment === 'top',
        [styles['bottom-align']]: alignment === 'bottom',
        // direction
        [styles['top-direction']]: direction === 'top',
        [styles['bottom-direction']]: direction === 'bottom',
        [styles['left-direction']]: direction === 'left',
        [styles['right-direction']]: direction === 'right',
        [styles.border]: hasBorder
      })}
      ref={contentRef}
      id={contentId}
    >
      {children}
    </div>
  );
};

export default PopoverContent;
