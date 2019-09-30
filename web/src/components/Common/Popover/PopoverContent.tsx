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

import { KEYCODE_ESC } from 'jslib/helpers/keyboard';
import styles from './PopoverContent.scss';
import { Alignment, Direction } from './types';

interface Props {
  contentId: string;
  onDismiss: () => void;
  children: React.ReactNode;
  direction: Direction;
  closeOnOutsideClick?: boolean;
  closeOnEscapeKeydown?: boolean;
  contentClassName?: string;
  alignment?: Alignment;
  triggerEl?: HTMLElement;
  wrapperEl?: any;
  hasBorder?: boolean;
}

const PopoverContent: React.SFC<Props> = ({
  contentId,
  contentClassName,
  onDismiss,
  triggerEl,
  wrapperEl,
  children,
  alignment,
  direction,
  hasBorder,
  closeOnOutsideClick,
  closeOnEscapeKeydown
}) => {
  const contentRef = useRef(null);

  useEffect(() => {
    function handleOutsideClick(e) {
      const contentEl = contentRef.current;

      if (!contentEl) {
        return;
      }

      if (triggerEl && triggerEl.contains(e.target)) {
        return;
      }

      if (!contentEl.contains(e.target)) {
        onDismiss();
      }
    }

    // addLinkListeners adds event listeners to all links to handle outside click
    // because click events on react-router's Link elements do not bubble up to document
    function addLinkListeners() {
      const links = document.links;
      for (let i = 0, linksLength = links.length; i < linksLength; i++) {
        const link = links[i];

        link.addEventListener('click', handleOutsideClick);
      }
    }

    // removeLinkListeners cleans up any event listeners for handling outside click attached to links
    function removeLinkListeners() {
      const links = document.links;
      for (let i = 0, linksLength = links.length; i < linksLength; i++) {
        const link = links[i];

        link.removeEventListener('click', handleOutsideClick);
      }
    }

    if (closeOnOutsideClick) {
      document.addEventListener('click', handleOutsideClick);
      document.addEventListener('touchstart', handleOutsideClick);
      addLinkListeners();
    }

    return () => {
      document.removeEventListener('click', handleOutsideClick);
      document.removeEventListener('touchstart', handleOutsideClick);
      removeLinkListeners();
    };
  }, [closeOnOutsideClick, onDismiss, triggerEl, wrapperEl]);

  useEffect(() => {
    function handleKeydown(e: KeyboardEvent) {
      if (e.keyCode === KEYCODE_ESC) {
        e.stopPropagation();
        onDismiss();
      }
    }

    let targetEl;
    if (wrapperEl) {
      targetEl = wrapperEl;
    } else {
      targetEl = document;
    }

    if (closeOnEscapeKeydown) {
      targetEl.addEventListener('keyup', handleKeydown);
    }

    return () => {
      targetEl.removeEventListener('keyup', handleKeydown);
    };
  }, [closeOnEscapeKeydown, onDismiss, wrapperEl]);

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
